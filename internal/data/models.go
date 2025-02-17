package data

import "time"

// Location represents a geographic location with coordinates
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

// TripPlanRequest represents a trip planning request
type TripPlanRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Location  Location  `json:"location"` // 主要活动区域
	Budget    struct {
		Total    float64 `json:"total"`    // 总预算
		Hotel    float64 `json:"hotel"`    // 住宿预算/晚
		Food     float64 `json:"food"`     // 餐饮预算/天
		Activity float64 `json:"activity"` // 活动预算/天
	} `json:"budget"`
	Preferences struct {
		Activities []string `json:"activities"` // 活动偏好
		Cuisine    []string `json:"cuisine"`    // 餐饮偏好
		Hotel      []string `json:"hotel"`      // 住宿偏好
	} `json:"preferences"`
	PartySize    int      `json:"party_size"`   // 出行人数
	Requirements []string `json:"requirements"` // 特殊需求
}

// Attraction represents a tourist attraction
type Attraction struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
	Category    []string `json:"category"`
	Price       float64  `json:"price"`
	OpenHours   []string `json:"open_hours"`
	Tags        []string `json:"tags"`
	Rating      float64  `json:"rating"`
	Images      []string `json:"images,omitempty"`
}

// Restaurant represents a dining establishment
type Restaurant struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Location    Location `json:"location"`
	Cuisine     []string `json:"cuisine"`
	PriceRange  string   `json:"price_range"` // $ $$ $$$ $$$$
	Rating      float64  `json:"rating"`
	OpenHours   []string `json:"open_hours"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// Hotel represents an accommodation option
type Hotel struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Location      Location `json:"location"`
	Stars         int      `json:"stars"`
	PricePerNight float64  `json:"price_per_night"`
	Amenities     []string `json:"amenities"`
	Description   string   `json:"description"`
	Rating        float64  `json:"rating"`
	Images        []string `json:"images,omitempty"`
}

// Weather represents weather information for a specific date and location
type Weather struct {
	Date        time.Time `json:"date"`
	Location    Location  `json:"location"`
	Temperature struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"temperature"`
	Condition     string  `json:"condition"`
	Humidity      float64 `json:"humidity"`
	WindSpeed     float64 `json:"wind_speed"`
	Precipitation float64 `json:"precipitation"`
}

// DailyPlan represents a single day's itinerary
type DailyPlan struct {
	Date       time.Time  `json:"date"`
	Weather    Weather    `json:"weather"`
	Activities []Activity `json:"activities"`
	Meals      []Meal     `json:"meals"`
	TotalCost  float64    `json:"total_cost"`
}

// Activity represents a planned activity
type Activity struct {
	Attraction Attraction `json:"attraction"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time"`
	Cost       float64    `json:"cost"`
	Notes      string     `json:"notes,omitempty"`
}

// Meal represents a planned meal
type Meal struct {
	Restaurant Restaurant `json:"restaurant"`
	Time       time.Time  `json:"time"`
	Type       string     `json:"type"` // breakfast, lunch, dinner
	Cost       float64    `json:"cost"`
}

// TripPlan represents the complete trip itinerary
type TripPlan struct {
	Request    TripPlanRequest `json:"request"`
	Hotel      Hotel           `json:"hotel"`
	DailyPlans []DailyPlan     `json:"daily_plans"`
	TotalCost  float64         `json:"total_cost"`
	Summary    string          `json:"summary"`
	Tips       []string        `json:"tips"`
}
