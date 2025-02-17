package data

import (
	"math"
)

const (
	// Earth's radius in kilometers
	earthRadiusKm = 6371.0
)

// distanceResult stores distance calculation result with original index
type distanceResult struct {
	distance float64
	index    int
}

// haversineDistance calculates the distance between two locations using the Haversine formula
func haversineDistance(a, b Location) float64 {
	// Convert latitude and longitude to radians
	lat1 := toRadians(a.Latitude)
	lon1 := toRadians(a.Longitude)
	lat2 := toRadians(b.Latitude)
	lon2 := toRadians(b.Longitude)

	// Differences in coordinates
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	// Haversine formula
	h := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Pow(math.Sin(dLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))

	// Calculate the distance
	return earthRadiusKm * c
}

// toRadians converts degrees to radians
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// LocationFilter represents a filter for location-based queries
type LocationFilter struct {
	Center     Location
	RadiusKm   float64
	MaxResults int
}

// SortLocationsByDistance sorts locations by distance from a reference point
func SortLocationsByDistance(reference Location, locations []Location) []Location {
	// Create a slice of distances with their original indices
	distances := make([]distanceResult, len(locations))

	// Calculate distances
	for i, loc := range locations {
		distances[i] = distanceResult{
			distance: haversineDistance(reference, loc),
			index:    i,
		}
	}

	// Sort by distance
	quickSortByDistance(distances)

	// Create sorted result
	result := make([]Location, len(locations))
	for i, di := range distances {
		result[i] = locations[di.index]
	}

	return result
}

// quickSortByDistance performs quicksort on the distance slice
func quickSortByDistance(distances []distanceResult) {
	if len(distances) < 2 {
		return
	}

	pivot := distances[0]
	left := 1
	right := len(distances) - 1

	for left <= right {
		if distances[left].distance <= pivot.distance {
			left++
			continue
		}
		if distances[right].distance > pivot.distance {
			right--
			continue
		}
		distances[left], distances[right] = distances[right], distances[left]
		left++
		right--
	}

	distances[0], distances[right] = distances[right], pivot
	quickSortByDistance(distances[:right])
	quickSortByDistance(distances[right+1:])
}
