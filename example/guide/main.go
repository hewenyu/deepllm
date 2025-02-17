package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hewenyu/deepllm/internal/data"
)

func main() {
	// Initialize data store
	store := data.NewStore("./data")

	// Load all data
	ctx := context.Background()
	if err := store.LoadAll(ctx); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	// Example 1: Get information about West Lake district
	district := store.GetDistrict("XH")
	if district != nil {
		fmt.Printf("District: %s\n", district.Name)
		fmt.Printf("Description: %s\n", district.Description)
		fmt.Printf("Landmarks: %v\n\n", district.Landmarks)
	}

	// Example 2: Find nearby attractions (within 2km of West Lake)
	westLakeLocation := data.Location{
		Latitude:  30.2587,
		Longitude: 120.1485,
	}
	nearbyAttractions := store.FindNearbyAttractions(westLakeLocation, 2.0)
	fmt.Printf("Attractions within 2km of West Lake:\n")
	for _, a := range nearbyAttractions {
		fmt.Printf("- %s\n", a.Name)
		fmt.Printf("  Description: %s\n", a.Description)
		fmt.Printf("  Price: %v %s\n", a.Price.Amount, a.Price.Currency)
		fmt.Printf("  Best Times: %v\n\n", a.RecommendedTime.BestTimes)
	}

	// Example 3: Find nearby restaurants with specific cuisine
	nearbyRestaurants := store.FindNearbyRestaurants(westLakeLocation, 1.0)
	fmt.Printf("Restaurants within 1km of West Lake:\n")
	for _, r := range nearbyRestaurants {
		fmt.Printf("- %s (%s)\n", r.Name, r.CuisineType)
		fmt.Printf("  Price Range: %v-%v %s (%s)\n",
			r.PriceRange.Min, r.PriceRange.Max,
			r.PriceRange.Currency, r.PriceRange.Level)
		fmt.Printf("  Signature Dishes: %v\n\n", r.SignatureDishes)
	}

	// Example 4: Check weather and suggest activities
	weather := store.GetWeatherForecast()
	if weather != nil && len(weather.DailyForecasts) > 0 {
		today := weather.DailyForecasts[0]
		fmt.Printf("Today's Weather (%s):\n", today.Date)
		fmt.Printf("Day: %s, Night: %s\n", today.Weather.Day, today.Weather.Night)
		fmt.Printf("Temperature: %v°C - %v°C\n", today.Temperature.Min, today.Temperature.Max)
		fmt.Printf("Tourism Suggestion: %s\n", today.Suggestion.Tourism)
		fmt.Printf("Notes: %s\n\n", today.Suggestion.Notes)

		// Suggest indoor/outdoor activities based on weather
		if today.Weather.Day == "晴" {
			attractions := store.GetAttractionsByDistrict("XH")
			fmt.Printf("Good weather for outdoor activities! Recommended spots:\n")
			for _, a := range attractions {
				if containsAny(a.Tags, []string{"自然景观", "户外活动"}) {
					fmt.Printf("- %s (%s)\n", a.Name, a.Description)
				}
			}
		} else {
			fmt.Printf("Consider indoor activities:\n")
			for _, a := range nearbyAttractions {
				if containsAny(a.Tags, []string{"博物馆", "艺术馆", "室内景点"}) {
					fmt.Printf("- %s (%s)\n", a.Name, a.Description)
				}
			}
		}
	}

	// Example 5: Find hotels by price range and distance
	hotels := store.FindHotelsByPriceRange(500, 2000)
	fmt.Printf("\nMid-range hotels near West Lake:\n")
	for _, h := range hotels {
		// Since we're using the store's methods for distance calculations,
		// we'll just print hotels that were found within our price range
		fmt.Printf("- %s (%s)\n", h.Name, h.Category)
		fmt.Printf("  Address: %s\n", h.Contact.Address)
		fmt.Printf("  Price Range: %v-%v %s\n",
			h.PriceRange.Min, h.PriceRange.Max, h.PriceRange.Currency)
		fmt.Printf("  Amenities: %v\n\n", h.Amenities)
	}
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
