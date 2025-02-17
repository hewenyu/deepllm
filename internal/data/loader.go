package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DataType represents different types of tourism data
type DataType string

const (
	TypeDistrict   DataType = "geographic/districts"
	TypeAttraction DataType = "tourism/attractions"
	TypeRestaurant DataType = "tourism/restaurants"
	TypeHotel      DataType = "tourism/hotels"
	TypeWeather    DataType = "weather/forecast"
	DataBasePath   string   = "data"
)

// DataLoader handles loading of tourism related data
type DataLoader struct {
	basePath string
}

// NewDataLoader creates a new data loader instance
func NewDataLoader(basePath string) *DataLoader {
	return &DataLoader{
		basePath: basePath,
	}
}

// loadJSONFile loads and unmarshals JSON data from a file
func (l *DataLoader) loadJSONFile(dataType DataType, v interface{}) error {
	filePath := filepath.Join(l.basePath, string(dataType)+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %w", filePath, err)
	}

	return nil
}
