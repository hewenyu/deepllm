package data

import (
	"context"
	"sync"
)

// Store manages all tourism related data
type Store struct {
	loader *DataLoader
	cache  struct {
		districts   []District
		attractions []Attraction
		restaurants []Restaurant
		hotels      []Hotel
		weather     *WeatherForecast
		sync.RWMutex
	}
}

// NewStore creates a new data store instance
func NewStore(dataPath string) *Store {
	return &Store{
		loader: NewDataLoader(dataPath),
	}
}

// LoadAll loads all data into memory
func (s *Store) LoadAll(ctx context.Context) error {
	s.cache.Lock()
	defer s.cache.Unlock()

	var districts struct {
		Districts []District `json:"districts"`
	}
	if err := s.loader.loadJSONFile(TypeDistrict, &districts); err != nil {
		return err
	}
	s.cache.districts = districts.Districts

	var attractions struct {
		Attractions []Attraction `json:"attractions"`
	}
	if err := s.loader.loadJSONFile(TypeAttraction, &attractions); err != nil {
		return err
	}
	s.cache.attractions = attractions.Attractions

	var restaurants struct {
		Restaurants []Restaurant `json:"restaurants"`
	}
	if err := s.loader.loadJSONFile(TypeRestaurant, &restaurants); err != nil {
		return err
	}
	s.cache.restaurants = restaurants.Restaurants

	var hotels struct {
		Hotels []Hotel `json:"hotels"`
	}
	if err := s.loader.loadJSONFile(TypeHotel, &hotels); err != nil {
		return err
	}
	s.cache.hotels = hotels.Hotels

	var weather WeatherForecast
	if err := s.loader.loadJSONFile(TypeWeather, &weather); err != nil {
		return err
	}
	s.cache.weather = &weather

	return nil
}

// Query Methods

// GetDistrict returns district by ID
func (s *Store) GetDistrict(id string) *District {
	s.cache.RLock()
	defer s.cache.RUnlock()

	for _, d := range s.cache.districts {
		if d.ID == id {
			return &d
		}
	}
	return nil
}

// GetAttractionsByDistrict returns attractions in a district
func (s *Store) GetAttractionsByDistrict(districtID string) []Attraction {
	s.cache.RLock()
	defer s.cache.RUnlock()

	var results []Attraction
	for _, a := range s.cache.attractions {
		if a.DistrictID == districtID {
			results = append(results, a)
		}
	}
	return results
}

// GetRestaurantsByDistrict returns restaurants in a district
func (s *Store) GetRestaurantsByDistrict(districtID string) []Restaurant {
	s.cache.RLock()
	defer s.cache.RUnlock()

	var results []Restaurant
	for _, r := range s.cache.restaurants {
		if r.DistrictID == districtID {
			results = append(results, r)
		}
	}
	return results
}

// GetHotelsByDistrict returns hotels in a district
func (s *Store) GetHotelsByDistrict(districtID string) []Hotel {
	s.cache.RLock()
	defer s.cache.RUnlock()

	var results []Hotel
	for _, h := range s.cache.hotels {
		if h.DistrictID == districtID {
			results = append(results, h)
		}
	}
	return results
}

// GetWeatherForecast returns the current weather forecast
func (s *Store) GetWeatherForecast() *WeatherForecast {
	s.cache.RLock()
	defer s.cache.RUnlock()
	return s.cache.weather
}

// Advanced Query Methods

// FindNearbyAttractions returns attractions within given distance (km) from location
func (s *Store) FindNearbyAttractions(loc Location, distanceKm float64) []Attraction {
	s.cache.RLock()
	defer s.cache.RUnlock()

	var results []Attraction
	for _, a := range s.cache.attractions {
		if haversineDistance(loc, a.Coordinates) <= distanceKm {
			results = append(results, a)
		}
	}
	return results
}

// FindRestaurantsByCuisine returns restaurants of given cuisine type
func (s *Store) FindRestaurantsByCuisine(cuisineType string) []Restaurant {
	s.cache.RLock()
	defer s.cache.RUnlock()

	var results []Restaurant
	for _, r := range s.cache.restaurants {
		if r.CuisineType == cuisineType {
			results = append(results, r)
		}
	}
	return results
}

// FindHotelsByPriceRange returns hotels within given price range
func (s *Store) FindHotelsByPriceRange(minPrice, maxPrice float64) []Hotel {
	s.cache.RLock()
	defer s.cache.RUnlock()

	var results []Hotel
	for _, h := range s.cache.hotels {
		if h.PriceRange.Min >= minPrice && h.PriceRange.Max <= maxPrice {
			results = append(results, h)
		}
	}
	return results
}

// FindNearbyRestaurants returns restaurants within given distance (km) from location
func (s *Store) FindNearbyRestaurants(loc Location, distanceKm float64) []Restaurant {
	s.cache.RLock()
	defer s.cache.RUnlock()

	var results []Restaurant
	for _, r := range s.cache.restaurants {
		if haversineDistance(loc, r.Coordinates) <= distanceKm {
			results = append(results, r)
		}
	}
	return results
}
