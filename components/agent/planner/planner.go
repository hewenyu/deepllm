package planner

import (
	"context"

	"github.com/hewenyu/deepllm/internal/data"
)

// TripAgent represents a trip planning agent
type TripAgent struct {
	store *data.Store
}

// NewTripAgent creates a new trip planning agent
func NewTripAgent(store *data.Store) *TripAgent {
	return &TripAgent{
		store: store,
	}
}

// PlanRequest contains parameters for trip planning
type PlanRequest struct {
	Location    data.Location `json:"location"`
	Duration    int           `json:"duration_days"`    // Trip duration in days
	Budget      float64       `json:"budget"`           // Total budget in CNY
	Preferences []string      `json:"preferences"`      // e.g., ["历史", "美食", "购物"]
	Weather     bool          `json:"consider_weather"` // Whether to consider weather in planning
}

// DayPlan represents a single day's itinerary
type DayPlan struct {
	Date     string              `json:"date"`
	Weather  *data.DailyForecast `json:"weather,omitempty"`
	Schedule []Activity          `json:"schedule"`
	Notes    []string            `json:"notes"`
}

// Activity represents a planned activity
type Activity struct {
	Type       string           `json:"type"` // "景点", "餐饮", "休息"
	Time       string           `json:"time"`
	Location   data.Location    `json:"location"`
	Attraction *data.Attraction `json:"attraction,omitempty"`
	Restaurant *data.Restaurant `json:"restaurant,omitempty"`
	Duration   int              `json:"duration_minutes"`
	Notes      string           `json:"notes"`
}

// Plan generates a trip plan based on the request
func (a *TripAgent) Plan(ctx context.Context, req PlanRequest) ([]DayPlan, error) {
	// TODO: Implement multi-agent planning logic
	return nil, nil
}

// suggestAttractions suggests attractions based on preferences and weather
func (a *TripAgent) suggestAttractions(ctx context.Context, loc data.Location, weather *data.DailyForecast, prefs []string) []data.Attraction {
	// Get nearby attractions
	attractions := a.store.FindNearbyAttractions(loc, 5.0) // 5km radius

	// Filter by weather conditions if needed
	if weather != nil && weather.Weather.Day != "晴" {
		// Prefer indoor attractions in bad weather
		var filtered []data.Attraction
		for _, attr := range attractions {
			if containsAny(attr.Tags, []string{"博物馆", "艺术馆", "室内景点"}) {
				filtered = append(filtered, attr)
			}
		}
		attractions = filtered
	}

	// Match preferences
	if len(prefs) > 0 {
		var matched []data.Attraction
		for _, attr := range attractions {
			if containsAny(attr.Tags, prefs) {
				matched = append(matched, attr)
			}
		}
		attractions = matched
	}

	return attractions
}

// Helper function to check if a slice contains any of the target strings
func containsAny(slice []string, targets []string) bool {
	for _, s := range slice {
		for _, t := range targets {
			if s == t {
				return true
			}
		}
	}
	return false
}
