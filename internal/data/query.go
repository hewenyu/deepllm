package data

import (
	"strings"
	"time"
)

// DataQuery provides methods to query and filter data
type DataQuery struct {
	Loader *DataLoader
}

// NewDataQuery creates a new DataQuery instance
func NewDataQuery(loader *DataLoader) *DataQuery {
	return &DataQuery{
		Loader: loader,
	}
}

// FilterAttractionsByPreferences filters attractions based on user preferences
func (q *DataQuery) FilterAttractionsByPreferences(attractions []Attraction, preferences []string) []Attraction {
	if len(preferences) == 0 {
		return attractions
	}

	var filtered []Attraction
	for _, attraction := range attractions {
		for _, pref := range preferences {
			prefLower := strings.ToLower(pref)
			// Check if preference matches category or tags
			for _, category := range attraction.Category {
				if strings.Contains(strings.ToLower(category), prefLower) {
					filtered = append(filtered, attraction)
					goto nextAttraction
				}
			}
			for _, tag := range attraction.Tags {
				if strings.Contains(strings.ToLower(tag), prefLower) {
					filtered = append(filtered, attraction)
					goto nextAttraction
				}
			}
		}
	nextAttraction:
	}
	return filtered
}

// FilterRestaurantsByPreferences filters restaurants based on cuisine preferences
func (q *DataQuery) FilterRestaurantsByPreferences(restaurants []Restaurant, cuisinePrefs []string) []Restaurant {
	if len(cuisinePrefs) == 0 {
		return restaurants
	}

	var filtered []Restaurant
	for _, restaurant := range restaurants {
		for _, pref := range cuisinePrefs {
			prefLower := strings.ToLower(pref)
			for _, cuisine := range restaurant.Cuisine {
				if strings.Contains(strings.ToLower(cuisine), prefLower) {
					filtered = append(filtered, restaurant)
					goto nextRestaurant
				}
			}
		}
	nextRestaurant:
	}
	return filtered
}

// FilterHotelsByPreferences filters hotels based on preferences and budget
func (q *DataQuery) FilterHotelsByPreferences(hotels []Hotel, preferences []string, maxPrice float64) []Hotel {
	var filtered []Hotel
	for _, hotel := range hotels {
		if hotel.PricePerNight > maxPrice {
			continue
		}

		if len(preferences) == 0 {
			filtered = append(filtered, hotel)
			continue
		}

		for _, pref := range preferences {
			prefLower := strings.ToLower(pref)
			for _, amenity := range hotel.Amenities {
				if strings.Contains(strings.ToLower(amenity), prefLower) {
					filtered = append(filtered, hotel)
					goto nextHotel
				}
			}
		}
	nextHotel:
	}
	return filtered
}

// GetWeatherForDate gets weather information for a specific date and location
func (q *DataQuery) GetWeatherForDate(weather []Weather, date time.Time, location Location) (Weather, bool) {
	for _, w := range weather {
		if w.Date.Year() == date.Year() &&
			w.Date.Month() == date.Month() &&
			w.Date.Day() == date.Day() &&
			CalculateDistance(w.Location, location) < 10 { // Within 10km radius
			return w, true
		}
	}
	return Weather{}, false
}

// FilterByBudget filters places by budget constraints
func (q *DataQuery) FilterByBudget(attractions []Attraction, maxBudget float64) []Attraction {
	var filtered []Attraction
	for _, attraction := range attractions {
		if attraction.Price <= maxBudget {
			filtered = append(filtered, attraction)
		}
	}
	return filtered
}

// SortByRating sorts places by their rating in descending order
func (q *DataQuery) SortByRating(attractions []Attraction) []Attraction {
	sorted := make([]Attraction, len(attractions))
	copy(sorted, attractions)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Rating > sorted[i].Rating {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}
