package data

import "time"

// Geographic Models
type District struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Coordinates    Location `json:"coordinates"`
	AreaKm2        float64  `json:"area_km2"`
	Transportation []string `json:"transportation"`
	Landmarks      []string `json:"landmarks"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Tourism Models
type Attraction struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	DistrictID      string   `json:"district_id"`
	Description     string   `json:"description"`
	Coordinates     Location `json:"coordinates"`
	Price           Price    `json:"price"`
	OpeningHours    Hours    `json:"opening_hours"`
	RecommendedTime struct {
		Hours     int      `json:"hours"`
		BestTimes []string `json:"best_times"`
	} `json:"recommended_time"`
	Highlights []string          `json:"highlights"`
	Tags       []string          `json:"tags"`
	CrowdLevel map[string]string `json:"crowd_level"`
}

type Restaurant struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DistrictID  string   `json:"district_id"`
	Description string   `json:"description"`
	Coordinates Location `json:"coordinates"`
	CuisineType string   `json:"cuisine_type"`
	PriceRange  struct {
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
		Currency string  `json:"currency"`
		Level    string  `json:"level"`
	} `json:"price_range"`
	OpeningHours      Hours    `json:"opening_hours"`
	SignatureDishes   []string `json:"signature_dishes"`
	Features          []string `json:"features"`
	ReservationNeeded bool     `json:"reservations_required"`
	Contact           Contact  `json:"contact"`
}

type Hotel struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DistrictID  string   `json:"district_id"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Coordinates Location `json:"coordinates"`
	PriceRange  struct {
		Min      float64 `json:"min"`
		Max      float64 `json:"max"`
		Currency string  `json:"currency"`
		Level    string  `json:"level"` // 经济, 中等, 高端
		Notes    string  `json:"notes"`
	} `json:"price_range"`
	Rooms     []Room    `json:"rooms"`
	Amenities []string  `json:"amenities"`
	Transport Transport `json:"transportation"`
	Contact   Contact   `json:"contact"`
}

// Weather Models
type WeatherForecast struct {
	City           string          `json:"city"`
	UpdateTime     time.Time       `json:"update_time"`
	Source         string          `json:"source"`
	DailyForecasts []DailyForecast `json:"daily_forecasts"`
	SpecialNotices []Notice        `json:"special_notices"`
}

// Common Structs
type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Notes    string  `json:"notes"`
}

type Hours struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	Notes     string `json:"notes,omitempty"`
	BreakTime *struct {
		Start string `json:"start,omitempty"`
		End   string `json:"end,omitempty"`
	} `json:"break_time,omitempty"`
}

type Room struct {
	Type     string   `json:"type"`
	SizeSqm  float64  `json:"size_sqm"`
	Price    float64  `json:"price"`
	Features []string `json:"features"`
}

type Transport struct {
	FromAirport struct {
		TaxiTime   string  `json:"taxi_time"`
		DistanceKm float64 `json:"distance_km"`
	} `json:"from_airport"`
	NearbyStations []string `json:"nearby_stations"`
}

type Contact struct {
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Email   string `json:"email,omitempty"`
}

type DailyForecast struct {
	Date    string `json:"date"`
	Weather struct {
		Day   string `json:"day"`
		Night string `json:"night"`
	} `json:"weather"`
	Temperature   Temperature `json:"temperature"`
	Humidity      Range       `json:"humidity"`
	Wind          Wind        `json:"wind"`
	Precipitation struct {
		Probability float64 `json:"probability"`
		Amount      float64 `json:"amount"`
		Unit        string  `json:"unit"`
	} `json:"precipitation"`
	AirQuality AirQuality `json:"air_quality"`
	Suggestion Suggestion `json:"suggestion"`
}

type Temperature struct {
	Max  float64 `json:"max"`
	Min  float64 `json:"min"`
	Unit string  `json:"unit"`
}

type Range struct {
	Max  float64 `json:"max"`
	Min  float64 `json:"min"`
	Unit string  `json:"unit"`
}

type Wind struct {
	Direction string `json:"direction"`
	Speed     Range  `json:"speed"`
}

type AirQuality struct {
	AQI              int    `json:"aqi"`
	Level            string `json:"level"`
	PrimaryPollutant string `json:"primary_pollutant"`
}

type Suggestion struct {
	Tourism string `json:"tourism"`
	Comfort string `json:"comfort"`
	Notes   string `json:"notes"`
}

type Notice struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
