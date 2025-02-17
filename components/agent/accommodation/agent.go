package accommodation

import (
	"context"
	"deepllm/components/agent"
	"deepllm/components/mock"
	"deepllm/internal/data"
	"fmt"
)

// AccommodationAgent handles hotel recommendations
type AccommodationAgent struct {
	*agent.BaseAgent
}

// NewAccommodationAgent creates a new accommodation agent
func NewAccommodationAgent(model mock.ChatModel, tools []mock.Tool, dataQuery *data.DataQuery) *AccommodationAgent {
	return &AccommodationAgent{
		BaseAgent: agent.NewBaseAgent("accommodation", model, tools, dataQuery),
	}
}

// Process processes the accommodation request
func (a *AccommodationAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
	request, ok := input.(*data.TripPlanRequest)
	if !ok {
		return nil, fmt.Errorf("invalid input type for accommodation agent")
	}

	// Load hotels data
	hotels, err := a.DataQuery.Loader.LoadHotels()
	if err != nil {
		return nil, fmt.Errorf("failed to load hotels: %v", err)
	}

	// Find nearby hotels
	nearbyHotels := data.FindNearbyHotels(hotels, request.Location, 5.0) // Within 5km

	// Filter by preferences and budget
	filteredHotels := a.DataQuery.FilterHotelsByPreferences(
		nearbyHotels,
		request.Preferences.Hotel,
		request.Budget.Hotel,
	)

	// Sort hotels by rating
	sortedHotels := make([]data.Hotel, len(filteredHotels))
	copy(sortedHotels, filteredHotels)
	for i := 0; i < len(sortedHotels)-1; i++ {
		for j := i + 1; j < len(sortedHotels); j++ {
			if sortedHotels[j].Rating > sortedHotels[i].Rating {
				sortedHotels[i], sortedHotels[j] = sortedHotels[j], sortedHotels[i]
			}
		}
	}

	// Use LLM to analyze and recommend hotels
	systemPrompt := a.BuildPrompt(
		"Hotel Recommendation Specialist",
		"analyze hotel options and provide personalized recommendations based on user preferences and requirements",
	)

	agent, err := a.CreateReactAgent(ctx, systemPrompt)
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
			Content: fmt.Sprintf("Please recommend hotels based on the following criteria:\n"+
				"Budget per night: %.2f\n"+
				"Preferred amenities: %v\n"+
				"Number of people: %d\n"+
				"Special requirements: %v\n\n"+
				"Available hotels: %+v",
				request.Budget.Hotel,
				request.Preferences.Hotel,
				request.PartySize,
				request.Requirements,
				sortedHotels,
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
