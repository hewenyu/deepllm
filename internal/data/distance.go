package data

import "math"

const (
	// earthRadius is the radius of Earth in kilometers
	earthRadius = 6371.0
)

// CalculateDistance calculates the distance between two locations using the Haversine formula
func CalculateDistance(loc1, loc2 Location) float64 {
	lat1 := toRadians(loc1.Latitude)
	lon1 := toRadians(loc1.Longitude)
	lat2 := toRadians(loc2.Latitude)
	lon2 := toRadians(loc2.Longitude)

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// toRadians converts degrees to radians
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// FindNearbyAttractions finds attractions within the specified radius (in km) from a location
func FindNearbyAttractions(attractions []Attraction, center Location, radius float64) []Attraction {
	var nearby []Attraction
	for _, attraction := range attractions {
		if distance := CalculateDistance(center, attraction.Location); distance <= radius {
			nearby = append(nearby, attraction)
		}
	}
	return nearby
}

// FindNearbyRestaurants finds restaurants within the specified radius (in km) from a location
func FindNearbyRestaurants(restaurants []Restaurant, center Location, radius float64) []Restaurant {
	var nearby []Restaurant
	for _, restaurant := range restaurants {
		if distance := CalculateDistance(center, restaurant.Location); distance <= radius {
			nearby = append(nearby, restaurant)
		}
	}
	return nearby
}

// FindNearbyHotels finds hotels within the specified radius (in km) from a location
func FindNearbyHotels(hotels []Hotel, center Location, radius float64) []Hotel {
	var nearby []Hotel
	for _, hotel := range hotels {
		if distance := CalculateDistance(center, hotel.Location); distance <= radius {
			nearby = append(nearby, hotel)
		}
	}
	return nearby
}
