package weather

import (
	"context"
	"deepllm/components/agent"
	"deepllm/components/mock"
	"deepllm/internal/data"
	"fmt"
	"time"
)

// WeatherAgent handles weather-based recommendations
type WeatherAgent struct {
	*agent.BaseAgent
}

// NewWeatherAgent creates a new weather agent
func NewWeatherAgent(model mock.ChatModel, tools []mock.Tool, dataQuery *data.DataQuery) *WeatherAgent {
	return &WeatherAgent{
		BaseAgent: agent.NewBaseAgent("weather", model, tools, dataQuery),
	}
}

// Process processes the weather request
func (w *WeatherAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
	request, ok := input.(*data.TripPlanRequest)
	if !ok {
		return nil, fmt.Errorf("invalid input type for weather agent")
	}

	// Load weather data
	weatherData, err := w.DataQuery.Loader.LoadWeather()
	if err != nil {
		return nil, fmt.Errorf("failed to load weather data: %v", err)
	}

	// Get weather forecasts for each day of the trip
	var forecasts []data.Weather
	for date := request.StartDate; !date.After(request.EndDate); date = date.Add(24 * time.Hour) {
		if weather, found := w.DataQuery.GetWeatherForDate(weatherData, date, request.Location); found {
			forecasts = append(forecasts, weather)
		}
	}

	// Use LLM to analyze weather and provide recommendations
	systemPrompt := w.BuildPrompt(
		"Weather Advisory Specialist",
		"analyze weather conditions and provide recommendations for activities and necessary preparations",
	)

	agent, err := w.CreateReactAgent(ctx, systemPrompt)
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
			Content: fmt.Sprintf("Please analyze the weather conditions and provide recommendations for the trip:\n"+
				"Trip dates: %s to %s\n"+
				"Location: %s\n"+
				"Planned activities: %v\n"+
				"Weather forecasts:\n%+v",
				request.StartDate.Format("2006-01-02"),
				request.EndDate.Format("2006-01-02"),
				request.Location.Name,
				request.Preferences.Activities,
				forecasts,
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

// GetWeatherSuitability determines if the weather is suitable for outdoor activities
func (w *WeatherAgent) GetWeatherSuitability(weather data.Weather) (bool, string) {
	// Check for severe weather conditions
	if weather.Precipitation > 25.0 { // Heavy rain (mm)
		return false, "Heavy rain expected"
	}
	if weather.WindSpeed > 30.0 { // Strong wind (km/h)
		return false, "Strong winds expected"
	}

	// Check temperature comfort range (15-30°C is generally comfortable)
	if weather.Temperature.Max > 35.0 {
		return false, "Temperature too high for outdoor activities"
	}
	if weather.Temperature.Min < 5.0 {
		return false, "Temperature too low for outdoor activities"
	}

	return true, "Weather conditions are suitable for outdoor activities"
}
