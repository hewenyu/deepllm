package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/hewenyu/deepllm/internal/data"
)

// AttractionParams represents parameters for attraction queries
type AttractionParams struct {
	DistrictID string  `json:"district_id,omitempty" jsonschema:"description=District ID to search in"`
	Latitude   float64 `json:"latitude,omitempty" jsonschema:"description=Latitude of the location"`
	Longitude  float64 `json:"longitude,omitempty" jsonschema:"description=Longitude of the location"`
	Distance   float64 `json:"distance,omitempty" jsonschema:"description=Search radius in kilometers"`
}

// RestaurantParams represents parameters for restaurant queries
type RestaurantParams struct {
	DistrictID  string  `json:"district_id,omitempty" jsonschema:"description=District ID to search in"`
	Latitude    float64 `json:"latitude,omitempty" jsonschema:"description=Latitude of the location"`
	Longitude   float64 `json:"longitude,omitempty" jsonschema:"description=Longitude of the location"`
	Distance    float64 `json:"distance,omitempty" jsonschema:"description=Search radius in kilometers"`
	CuisineType string  `json:"cuisine_type,omitempty" jsonschema:"description=Type of cuisine"`
}

// NewAttractionTool creates a new attraction search tool
func NewAttractionTool(store *data.Store) (tool.InvokableTool, error) {
	return utils.InferTool(
		"search_attractions",
		"Search for attractions by district or location",
		func(_ context.Context, params *AttractionParams) (string, error) {
			var attractions []data.Attraction

			if params.DistrictID != "" {
				attractions = store.GetAttractionsByDistrict(params.DistrictID)
			} else if params.Latitude != 0 && params.Longitude != 0 && params.Distance != 0 {
				loc := data.Location{
					Latitude:  params.Latitude,
					Longitude: params.Longitude,
				}
				attractions = store.FindNearbyAttractions(loc, params.Distance)
			} else {
				return "", fmt.Errorf("either district_id or location with distance must be provided")
			}

			response := map[string]interface{}{
				"attractions": attractions,
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return "", fmt.Errorf("failed to marshal response: %v", err)
			}

			return string(jsonResponse), nil
		},
	)
}

// NewRestaurantTool creates a new restaurant search tool
func NewRestaurantTool(store *data.Store) (tool.InvokableTool, error) {
	return utils.InferTool(
		"search_restaurants",
		"Search for restaurants by district, location, or cuisine type",
		func(_ context.Context, params *RestaurantParams) (string, error) {
			var restaurants []data.Restaurant

			if params.DistrictID != "" {
				restaurants = store.GetRestaurantsByDistrict(params.DistrictID)
			} else if params.Latitude != 0 && params.Longitude != 0 && params.Distance != 0 {
				loc := data.Location{
					Latitude:  params.Latitude,
					Longitude: params.Longitude,
				}
				restaurants = store.FindNearbyRestaurants(loc, params.Distance)
			} else {
				return "", fmt.Errorf("either district_id or location with distance must be provided")
			}

			// Filter by cuisine type if specified
			if params.CuisineType != "" {
				filtered := make([]data.Restaurant, 0)
				for _, r := range restaurants {
					if r.CuisineType == params.CuisineType {
						filtered = append(filtered, r)
					}
				}
				restaurants = filtered
			}

			response := map[string]interface{}{
				"restaurants": restaurants,
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return "", fmt.Errorf("failed to marshal response: %v", err)
			}

			return string(jsonResponse), nil
		},
	)
}

// NewWeatherTool creates a new weather forecast tool
func NewWeatherTool(store *data.Store) (tool.InvokableTool, error) {
	return utils.InferTool(
		"get_weather",
		"Get current weather forecast",
		func(_ context.Context, _ *struct{}) (string, error) {
			weather := store.GetWeatherForecast()
			if weather == nil {
				return "", fmt.Errorf("weather forecast not available")
			}

			jsonResponse, err := json.Marshal(weather)
			if err != nil {
				return "", fmt.Errorf("failed to marshal response: %v", err)
			}

			return string(jsonResponse), nil
		},
	)
}
