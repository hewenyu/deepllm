package planner

import (
	"context"
	"deepllm/components/agent"
	"deepllm/components/agent/accommodation"
	"deepllm/components/agent/dining"
	"deepllm/components/agent/weather"
	"deepllm/components/mock"
	"deepllm/internal/data"
	"fmt"
	"time"
)

// PlannerAgent handles overall trip planning
type PlannerAgent struct {
	*agent.BaseAgent
	weatherAgent       *weather.WeatherAgent
	accommodationAgent *accommodation.AccommodationAgent
	diningAgent        *dining.DiningAgent
}

// NewPlannerAgent creates a new planner agent
func NewPlannerAgent(model mock.ChatModel, tools []mock.Tool, dataQuery *data.DataQuery) *PlannerAgent {
	return &PlannerAgent{
		BaseAgent:          agent.NewBaseAgent("planner", model, tools, dataQuery),
		weatherAgent:       weather.NewWeatherAgent(model, tools, dataQuery),
		accommodationAgent: accommodation.NewAccommodationAgent(model, tools, dataQuery),
		diningAgent:        dining.NewDiningAgent(model, tools, dataQuery),
	}
}

// Process processes the trip planning request
func (p *PlannerAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
	request, ok := input.(*data.TripPlanRequest)
	if !ok {
		return nil, fmt.Errorf("invalid input type for planner agent")
	}

	// Get weather recommendations
	weatherResult, err := p.weatherAgent.Process(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather recommendations: %v", err)
	}

	// Get accommodation recommendations
	accommodationResult, err := p.accommodationAgent.Process(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get accommodation recommendations: %v", err)
	}

	// Get dining recommendations
	diningResult, err := p.diningAgent.Process(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get dining recommendations: %v", err)
	}

	// Load attractions data
	attractions, err := p.DataQuery.Loader.LoadAttractions()
	if err != nil {
		return nil, fmt.Errorf("failed to load attractions: %v", err)
	}

	// Filter attractions by preferences and budget
	nearbyAttractions := data.FindNearbyAttractions(attractions, request.Location, 10.0) // Within 10km
	filteredAttractions := p.DataQuery.FilterAttractionsByPreferences(
		nearbyAttractions,
		request.Preferences.Activities,
	)
	budgetedAttractions := p.DataQuery.FilterByBudget(
		filteredAttractions,
		request.Budget.Activity,
	)

	// Use LLM to create the final trip plan
	systemPrompt := p.BuildPrompt(
		"Trip Planning Specialist",
		"create a comprehensive trip plan that combines weather conditions, accommodations, dining options, and attractions",
	)

	agent, err := p.CreateReactAgent(ctx, systemPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %v", err)
	}

	// Calculate trip duration
	duration := int(request.EndDate.Sub(request.StartDate).Hours() / 24)

	// Prepare input for LLM
	messages := []*mock.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role: "user",
			Content: fmt.Sprintf("Please create a comprehensive trip plan for %d days based on the following information:\n"+
				"Trip dates: %s to %s\n"+
				"Location: %s\n"+
				"Party size: %d\n"+
				"Budget:\n"+
				"  - Total: %.2f\n"+
				"  - Hotel per night: %.2f\n"+
				"  - Food per day: %.2f\n"+
				"  - Activities per day: %.2f\n"+
				"Preferences:\n"+
				"  - Activities: %v\n"+
				"  - Cuisine: %v\n"+
				"  - Hotel: %v\n"+
				"Special requirements: %v\n\n"+
				"Weather recommendations: %s\n"+
				"Accommodation recommendations: %s\n"+
				"Dining recommendations: %s\n"+
				"Available attractions: %+v",
				duration,
				request.StartDate.Format("2006-01-02"),
				request.EndDate.Format("2006-01-02"),
				request.Location.Name,
				request.PartySize,
				request.Budget.Total,
				request.Budget.Hotel,
				request.Budget.Food,
				request.Budget.Activity,
				request.Preferences.Activities,
				request.Preferences.Cuisine,
				request.Preferences.Hotel,
				request.Requirements,
				weatherResult.(*mock.Message).Content,
				accommodationResult.(*mock.Message).Content,
				diningResult.(*mock.Message).Content,
				budgetedAttractions,
			),
		},
	}

	// Get the final trip plan from LLM
	result, err := agent.Invoke(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to create trip plan: %v", err)
	}

	return result, nil
}

// createDailySchedule creates a schedule for a single day
func (p *PlannerAgent) createDailySchedule(
	ctx context.Context,
	date time.Time,
	weather data.Weather,
	attractions []data.Attraction,
	restaurants []data.Restaurant,
) (*data.DailyPlan, error) {
	// Check weather suitability for outdoor activities
	isOutdoorSuitable, weatherNote := p.weatherAgent.GetWeatherSuitability(weather)

	// Filter attractions based on weather
	var suitableAttractions []data.Attraction
	if isOutdoorSuitable {
		suitableAttractions = attractions
	} else {
		// Only include indoor attractions when weather is not suitable for outdoor activities
		for _, attraction := range attractions {
			for _, category := range attraction.Category {
				if category == "indoor" {
					suitableAttractions = append(suitableAttractions, attraction)
					break
				}
			}
		}
	}

	// Create daily plan
	plan := &data.DailyPlan{
		Date:    date,
		Weather: weather,
	}

	// Add activities and meals (this is a simplified version)
	if len(suitableAttractions) > 0 {
		plan.Activities = []data.Activity{
			{
				Attraction: suitableAttractions[0],
				StartTime:  time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, date.Location()),
				EndTime:    time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, date.Location()),
				Notes:      weatherNote,
			},
		}
	}

	if len(restaurants) > 0 {
		plan.Meals = []data.Meal{
			{
				Restaurant: restaurants[0],
				Time:       time.Date(date.Year(), date.Month(), date.Day(), 12, 30, 0, 0, date.Location()),
				Type:       "lunch",
			},
		}
	}

	// Calculate total cost
	for _, activity := range plan.Activities {
		plan.TotalCost += activity.Cost
	}
	for _, meal := range plan.Meals {
		plan.TotalCost += meal.Cost
	}

	return plan, nil
}
