package main

import (
	"context"
	"deepllm/components/agent/coordinator"
	"deepllm/components/mock"
	"deepllm/components/ollama"
	"deepllm/internal/data"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Initialize data loader and query
	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "./data"
	}
	dataLoader := data.NewDataLoader(dataPath)
	dataQuery := data.NewDataQuery(dataLoader)

	// Initialize Ollama chat model
	ollamaBaseURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaBaseURL == "" {
		ollamaBaseURL = "http://localhost:11434"
	}
	ollamaModel := os.Getenv("OLLAMA_MODEL")
	if ollamaModel == "" {
		ollamaModel = "deepseek-r1:14b"
	}
	chatModel := ollama.NewChatModel(ollamaBaseURL, ollamaModel)

	// Create tools
	tools := []mock.Tool{
		mock.NewMockTool(
			"search_attractions",
			"Search for attractions near a location",
			func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
				// Mock implementation
				return map[string]interface{}{
					"status": "success",
					"data":   "Mock attraction data",
				}, nil
			},
		),
		mock.NewMockTool(
			"search_restaurants",
			"Search for restaurants near a location",
			func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
				// Mock implementation
				return map[string]interface{}{
					"status": "success",
					"data":   "Mock restaurant data",
				}, nil
			},
		),
		mock.NewMockTool(
			"search_hotels",
			"Search for hotels near a location",
			func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
				// Mock implementation
				return map[string]interface{}{
					"status": "success",
					"data":   "Mock hotel data",
				}, nil
			},
		),
		mock.NewMockTool(
			"get_weather",
			"Get weather forecast for a location",
			func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
				// Mock implementation
				return map[string]interface{}{
					"status": "success",
					"data":   "Mock weather data",
				}, nil
			},
		),
	}

	// Create coordinator agent
	coordinatorAgent := coordinator.NewCoordinatorAgent(chatModel, tools, dataQuery)

	// Create sample trip request
	request := &data.TripPlanRequest{
		StartDate: time.Now().Add(24 * time.Hour),
		EndDate:   time.Now().Add(72 * time.Hour),
		Location: data.Location{
			Name:      "West Lake",
			Latitude:  30.2587,
			Longitude: 120.1315,
		},
		Budget: struct {
			Total    float64 `json:"total"`
			Hotel    float64 `json:"hotel"`
			Food     float64 `json:"food"`
			Activity float64 `json:"activity"`
		}{
			Total:    5000,
			Hotel:    1000,
			Food:     300,
			Activity: 200,
		},
		Preferences: struct {
			Activities []string `json:"activities"`
			Cuisine    []string `json:"cuisine"`
			Hotel      []string `json:"hotel"`
		}{
			Activities: []string{"sightseeing", "shopping", "cultural"},
			Cuisine:    []string{"local", "seafood"},
			Hotel:      []string{"luxury", "lake view"},
		},
		PartySize:    2,
		Requirements: []string{"wheelchair accessible"},
	}

	// Process the request
	ctx := context.Background()
	result, err := coordinatorAgent.Process(ctx, request)
	if err != nil {
		log.Fatalf("Failed to process request: %v", err)
	}

	// Print the result
	fmt.Printf("Trip Plan:\n%s\n", result.(*mock.Message).Content)
}
