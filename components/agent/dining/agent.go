package dining

import (
	"context"
	"deepllm/components/agent"
	"deepllm/components/mock"
	"deepllm/internal/data"
	"fmt"
)

// DiningAgent handles restaurant recommendations
type DiningAgent struct {
	*agent.BaseAgent
}

// NewDiningAgent creates a new dining agent
func NewDiningAgent(model mock.ChatModel, tools []mock.Tool, dataQuery *data.DataQuery) *DiningAgent {
	return &DiningAgent{
		BaseAgent: agent.NewBaseAgent("dining", model, tools, dataQuery),
	}
}

// Process processes the dining request
func (d *DiningAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
	request, ok := input.(*data.TripPlanRequest)
	if !ok {
		return nil, fmt.Errorf("invalid input type for dining agent")
	}

	// Load restaurants data
	restaurants, err := d.DataQuery.Loader.LoadRestaurants()
	if err != nil {
		return nil, fmt.Errorf("failed to load restaurants: %v", err)
	}

	// Find nearby restaurants
	nearbyRestaurants := data.FindNearbyRestaurants(restaurants, request.Location, 3.0) // Within 3km

	// Filter by cuisine preferences
	filteredRestaurants := d.DataQuery.FilterRestaurantsByPreferences(
		nearbyRestaurants,
		request.Preferences.Cuisine,
	)

	// Sort restaurants by rating
	sortedRestaurants := make([]data.Restaurant, len(filteredRestaurants))
	copy(sortedRestaurants, filteredRestaurants)
	for i := 0; i < len(sortedRestaurants)-1; i++ {
		for j := i + 1; j < len(sortedRestaurants); j++ {
			if sortedRestaurants[j].Rating > sortedRestaurants[i].Rating {
				sortedRestaurants[i], sortedRestaurants[j] = sortedRestaurants[j], sortedRestaurants[i]
			}
		}
	}

	// Use LLM to analyze and recommend restaurants
	systemPrompt := d.BuildPrompt(
		"Restaurant Recommendation Specialist",
		"analyze restaurant options and provide personalized dining recommendations based on user preferences and requirements",
	)

	agent, err := d.CreateReactAgent(ctx, systemPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %v", err)
	}

	// Prepare input for LLM
	messages := []*mock.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role: "user",
			Content: fmt.Sprintf("Please recommend restaurants based on the following criteria:\n"+
				"Daily food budget: %.2f\n"+
				"Preferred cuisines: %v\n"+
				"Number of people: %d\n"+
				"Special requirements: %v\n\n"+
				"Available restaurants: %+v",
				request.Budget.Food,
				request.Preferences.Cuisine,
				request.PartySize,
				request.Requirements,
				sortedRestaurants,
			),
		},
	}

	// Get recommendation from LLM
	result, err := agent.Invoke(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendation: %v", err)
	}

	return result, nil
}
