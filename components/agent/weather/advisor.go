package weather

import (
	"context"
	"fmt"
	"time"

	"github.com/hewenyu/deepllm/components/agent"
	"github.com/hewenyu/deepllm/internal/data"
)

// WeatherAgent specializes in weather-based activity recommendations
type WeatherAgent struct {
	*agent.BaseAgent
	store *data.Store
	name  string
	desc  string
}

// NewWeatherAgent creates a new weather advisor agent
func NewWeatherAgent(opts agent.BaseAgentOptions, store *data.Store) *WeatherAgent {
	return &WeatherAgent{
		BaseAgent: agent.NewBaseAgent(opts.Config),
		store:     store,
		name:      opts.Name,
		desc:      opts.Description,
	}
}

// Initialize implements agent.AgentInterface
func (a *WeatherAgent) Initialize(ctx context.Context) error {
	if a.store == nil {
		return fmt.Errorf("data store is not initialized")
	}
	return nil
}

// Name implements agent.AgentInterface
func (a *WeatherAgent) Name() string {
	if a.name == "" {
		return "WeatherAdvisor"
	}
	return a.name
}

// Description implements agent.AgentInterface
func (a *WeatherAgent) Description() string {
	if a.desc == "" {
		return "天气与活动建议智能体"
	}
	return a.desc
}

// Process implements agent.AgentInterface
func (a *WeatherAgent) Process(ctx context.Context, request interface{}) (agent.AgentResponse, error) {
	req, ok := request.(WeatherRequest)
	if !ok {
		return agent.AgentResponse{}, fmt.Errorf("invalid request type")
	}

	advice, err := a.GetAdvice(ctx, req.Date)
	if err != nil {
		return agent.AgentResponse{}, err
	}

	return agent.AgentResponse{
		Success: true,
		Data:    advice,
	}, nil
}

// WeatherRequest represents a request for weather advice
type WeatherRequest struct {
	Date time.Time `json:"date"`
}

// WeatherAdvice contains weather-based recommendations
type WeatherAdvice struct {
	Weather        *data.DailyForecast `json:"weather"`
	Suitable       []string            `json:"suitable_activities"`   // 适合的活动
	Unsuitable     []string            `json:"unsuitable_activities"` // 不适合的活动
	Precautions    []string            `json:"precautions"`           // 注意事项
	IndoorOptions  []string            `json:"indoor_options"`        // 室内备选
	OutdoorOptions []string            `json:"outdoor_options"`       // 室外备选
}

// GetAdvice provides weather-based activity recommendations
func (a *WeatherAgent) GetAdvice(ctx context.Context, date time.Time) (*WeatherAdvice, error) {
	forecast := a.store.GetWeatherForecast()
	if forecast == nil || len(forecast.DailyForecasts) == 0 {
		return nil, fmt.Errorf("no weather forecast available")
	}

	dateStr := date.Format("2006-01-02")

	// Find matching forecast
	var todayForecast *data.DailyForecast
	for _, f := range forecast.DailyForecasts {
		if f.Date == dateStr {
			todayForecast = &f
			break
		}
	}

	if todayForecast == nil {
		return nil, fmt.Errorf("no forecast found for date: %s", dateStr)
	}

	advice := &WeatherAdvice{
		Weather: todayForecast,
	}

	// Generate activity recommendations using LLM
	prompt := fmt.Sprintf(`Based on the following weather conditions, suggest activities:
Weather: %s to %s
Temperature: %.1f°C to %.1f°C
Rain Probability: %.0f%%
Wind Speed: %.0f-%.0f%s
Air Quality: %s (AQI: %d)

Please provide:
1. Suitable outdoor activities
2. Activities to avoid
3. Safety precautions
4. Indoor alternatives
5. Recommended outdoor activities if weather permits`,
		todayForecast.Weather.Day,
		todayForecast.Weather.Night,
		todayForecast.Temperature.Min,
		todayForecast.Temperature.Max,
		todayForecast.Precipitation.Probability,
		todayForecast.Wind.Speed.Min,
		todayForecast.Wind.Speed.Max,
		todayForecast.Wind.Speed.Unit,
		todayForecast.AirQuality.Level,
		todayForecast.AirQuality.AQI,
	)

	var result struct {
		Suitable       []string `json:"suitable"`
		Unsuitable     []string `json:"unsuitable"`
		Precautions    []string `json:"precautions"`
		IndoorOptions  []string `json:"indoor_options"`
		OutdoorOptions []string `json:"outdoor_options"`
	}

	if err := a.GenerateStructured(ctx, prompt, &result); err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %v", err)
	}

	advice.Suitable = result.Suitable
	advice.Unsuitable = result.Unsuitable
	advice.Precautions = result.Precautions
	advice.IndoorOptions = result.IndoorOptions
	advice.OutdoorOptions = result.OutdoorOptions

	return advice, nil
}
