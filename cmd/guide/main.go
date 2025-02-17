package main

import (
	"context"
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

	// Create a simple tool
	weatherTool := mock.NewMockTool(
		"get_weather",
		"Get weather forecast for a location",
		func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
			// Mock implementation
			return map[string]interface{}{
				"temperature": 25.0,
				"condition":   "sunny",
				"humidity":    60.0,
			}, nil
		},
	)

	// Create a simple prompt
	messages := []*mock.Message{
		{
			Role:    "system",
			Content: "You are a helpful travel assistant for Hangzhou.",
		},
		{
			Role: "user",
			Content: "What's the weather like at West Lake today? " +
				"Should I go for a boat ride?",
		},
	}

	// Call the model
	ctx := context.Background()
	chatModel.BindTools([]mock.Tool{weatherTool})
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	fmt.Printf("Assistant's response:\n%s\n", response.Content)

	// Demonstrate data querying
	fmt.Println("\nDemonstrating data querying:")

	// Load and filter attractions
	attractions, err := dataLoader.LoadAttractions()
	if err != nil {
		log.Fatalf("Failed to load attractions: %v", err)
	}

	location := data.Location{
		Name:      "West Lake",
		Latitude:  30.2587,
		Longitude: 120.1315,
	}

	nearbyAttractions := data.FindNearbyAttractions(attractions, location, 2.0) // Within 2km
	fmt.Printf("\nFound %d attractions within 2km of West Lake\n", len(nearbyAttractions))

	// Filter by preferences
	preferences := []string{"cultural", "historical"}
	filteredAttractions := dataQuery.FilterAttractionsByPreferences(nearbyAttractions, preferences)
	fmt.Printf("\nFound %d attractions matching preferences: %v\n", len(filteredAttractions), preferences)

	// Sort by rating
	sortedAttractions := dataQuery.SortByRating(filteredAttractions)
	fmt.Println("\nTop rated attractions:")
	for i, attraction := range sortedAttractions {
		if i >= 3 {
			break
		}
		fmt.Printf("%d. %s (Rating: %.1f)\n", i+1, attraction.Name, attraction.Rating)
	}

	// Demonstrate weather querying
	weatherData, err := dataLoader.LoadWeather()
	if err != nil {
		log.Fatalf("Failed to load weather data: %v", err)
	}

	weather, found := dataQuery.GetWeatherForDate(weatherData, time.Now(), location)
	if found {
		fmt.Printf("\nCurrent weather at %s:\n", location.Name)
		fmt.Printf("Temperature: %.1f°C - %.1f°C\n", weather.Temperature.Min, weather.Temperature.Max)
		fmt.Printf("Condition: %s\n", weather.Condition)
		fmt.Printf("Humidity: %.1f%%\n", weather.Humidity)
		fmt.Printf("Wind Speed: %.1f km/h\n", weather.WindSpeed)
	} else {
		fmt.Println("\nNo weather data available for the current date")
	}
}
