package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DataLoader handles loading data from JSON files
type DataLoader struct {
	BasePath string
}

// NewDataLoader creates a new DataLoader instance
func NewDataLoader(basePath string) *DataLoader {
	return &DataLoader{
		BasePath: basePath,
	}
}

// loadJSON loads data from a JSON file into the target interface
func (d *DataLoader) loadJSON(filename string, target interface{}) error {
	path := filepath.Join(d.BasePath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON from %s: %v", filename, err)
	}

	return nil
}

// LoadAttractions loads attractions data
func (d *DataLoader) LoadAttractions() ([]Attraction, error) {
	var attractions []Attraction
	err := d.loadJSON("attractions.json", &attractions)
	return attractions, err
}

// LoadRestaurants loads restaurants data
func (d *DataLoader) LoadRestaurants() ([]Restaurant, error) {
	var restaurants []Restaurant
	err := d.loadJSON("restaurants.json", &restaurants)
	return restaurants, err
}

// LoadHotels loads hotels data
func (d *DataLoader) LoadHotels() ([]Hotel, error) {
	var hotels []Hotel
	err := d.loadJSON("hotels.json", &hotels)
	return hotels, err
}

// LoadWeather loads weather data
func (d *DataLoader) LoadWeather() ([]Weather, error) {
	var weather []Weather
	err := d.loadJSON("weather.json", &weather)
	return weather, err
}
