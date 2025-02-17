package coordinator

import (
	"context"
	"fmt"
	"time"

	"github.com/hewenyu/deepllm/components/agent/accommodation"
	"github.com/hewenyu/deepllm/components/agent/dining"
	"github.com/hewenyu/deepllm/components/agent/weather"
	"github.com/hewenyu/deepllm/internal/data"
)

// TripPlanner coordinates multiple agents for trip planning
type TripPlanner struct {
	store           *data.Store
	weatherAgent    *weather.WeatherAgent
	restaurantAgent *dining.RestaurantAgent
	hotelAgent      *accommodation.HotelAgent
}

// NewTripPlanner creates a new trip planner
func NewTripPlanner(store *data.Store) *TripPlanner {
	return &TripPlanner{
		store:           store,
		weatherAgent:    weather.NewWeatherAgent(store),
		restaurantAgent: dining.NewRestaurantAgent(store),
		hotelAgent:      accommodation.NewHotelAgent(store),
	}
}

// TripPlanRequest represents a trip planning request
type TripPlanRequest struct {
	StartDate time.Time     `json:"start_date"`
	EndDate   time.Time     `json:"end_date"`
	Location  data.Location `json:"location"` // 主要活动区域
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

// DailyPlan represents a single day's itinerary
type DailyPlan struct {
	Date       string                        `json:"date"`
	Weather    *weather.WeatherAdvice        `json:"weather"`
	Activities []Activity                    `json:"activities"`
	Dining     []dining.DiningRecommendation `json:"dining"`
	Notes      []string                      `json:"notes"`
}

// Activity represents a planned activity
type Activity struct {
	Time       string           `json:"time"`
	Type       string           `json:"type"` // "景点", "活动", "休息"
	Location   data.Location    `json:"location"`
	Attraction *data.Attraction `json:"attraction,omitempty"`
	Duration   int              `json:"duration_minutes"`
	Notes      []string         `json:"notes"`
}

// TripPlan represents a complete trip plan
type TripPlan struct {
	Overview struct {
		Duration   int      `json:"duration_days"`
		TotalCost  float64  `json:"total_cost"`
		Highlights []string `json:"highlights"`
	} `json:"overview"`
	Accommodation *accommodation.HotelRecommendation `json:"accommodation"`
	DailyPlans    []DailyPlan                        `json:"daily_plans"`
	Tips          []string                           `json:"tips"`
}

// Plan generates a complete trip plan
func (p *TripPlanner) Plan(ctx context.Context, req TripPlanRequest) (*TripPlan, error) {
	days := int(req.EndDate.Sub(req.StartDate).Hours() / 24)
	if days < 1 {
		return nil, fmt.Errorf("invalid date range")
	}

	// Create plan structure
	plan := &TripPlan{}
	plan.Overview.Duration = days

	// Find suitable hotel
	hotelReq := accommodation.AccommodationRequest{
		Location:     req.Location,
		CheckIn:      req.StartDate,
		CheckOut:     req.EndDate,
		Budget:       req.Budget.Hotel,
		GuestCount:   req.PartySize,
		Distance:     2.0, // 2km radius
		Preferences:  req.Preferences.Hotel,
		Requirements: req.Requirements,
	}

	hotels, err := p.hotelAgent.Recommend(ctx, hotelReq)
	if err != nil {
		return nil, fmt.Errorf("hotel recommendation failed: %v", err)
	}
	if len(hotels) > 0 {
		plan.Accommodation = &hotels[0]
	}

	// Generate daily plans
	plan.DailyPlans = make([]DailyPlan, days)
	for i := 0; i < days; i++ {
		date := req.StartDate.AddDate(0, 0, i)
		dailyPlan, err := p.planDay(ctx, date, req)
		if err != nil {
			return nil, fmt.Errorf("planning day %d failed: %v", i+1, err)
		}
		plan.DailyPlans[i] = *dailyPlan
	}

	// Calculate total cost and generate highlights
	p.finalizeTrip(plan, req)

	return plan, nil
}

// planDay generates a single day's itinerary
func (p *TripPlanner) planDay(ctx context.Context, date time.Time, req TripPlanRequest) (*DailyPlan, error) {
	plan := &DailyPlan{
		Date: date.Format("2006-01-02"),
	}

	// Get weather advice
	weatherAdvice, err := p.weatherAgent.GetAdvice(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("weather advice failed: %v", err)
	}
	plan.Weather = weatherAdvice

	// Plan meals
	lunchReq := dining.DiningRequest{
		Location:    req.Location,
		Time:        date.Add(12 * time.Hour),
		Budget:      req.Budget.Food / 2, // Split budget between lunch and dinner
		Cuisine:     req.Preferences.Cuisine,
		PartySize:   req.PartySize,
		Distance:    2.0,
		Preferences: req.Preferences.Activities,
	}

	lunch, err := p.restaurantAgent.Recommend(ctx, lunchReq)
	if err == nil && len(lunch) > 0 {
		plan.Dining = append(plan.Dining, lunch[0])
	}

	dinnerReq := lunchReq
	dinnerReq.Time = date.Add(18 * time.Hour)
	dinner, err := p.restaurantAgent.Recommend(ctx, dinnerReq)
	if err == nil && len(dinner) > 0 {
		plan.Dining = append(plan.Dining, dinner[0])
	}

	// Plan activities based on weather
	p.planActivities(plan, weatherAdvice, req)

	return plan, nil
}

// planActivities plans activities based on weather and preferences
func (p *TripPlanner) planActivities(plan *DailyPlan, weather *weather.WeatherAdvice, req TripPlanRequest) {
	// Morning activity (9:00-12:00)
	if len(weather.Suitable) > 0 || len(weather.OutdoorOptions) > 0 {
		// Good weather - plan outdoor activity
		plan.Activities = append(plan.Activities, Activity{
			Time:     "09:00",
			Type:     "景点",
			Duration: 180, // 3 hours
			Notes:    weather.Precautions,
		})
	} else {
		// Bad weather - plan indoor activity
		plan.Activities = append(plan.Activities, Activity{
			Time:     "10:00",
			Type:     "室内活动",
			Duration: 120, // 2 hours
			Notes:    weather.Precautions,
		})
	}

	// Afternoon activity (14:00-17:00)
	if contains(weather.Unsuitable, "长时间户外活动") {
		// Plan indoor activities
		plan.Activities = append(plan.Activities, Activity{
			Time:     "14:00",
			Type:     "室内活动",
			Duration: 180,
			Notes:    append(weather.Precautions, "选择室内景点"),
		})
	} else {
		plan.Activities = append(plan.Activities, Activity{
			Time:     "14:00",
			Type:     "景点",
			Duration: 180,
			Notes:    weather.Precautions,
		})
	}

	// Evening activity (20:00-21:30)
	plan.Activities = append(plan.Activities, Activity{
		Time:     "20:00",
		Type:     "休闲活动",
		Duration: 90,
		Notes:    []string{"夜景观赏", "文化体验"},
	})
}

// finalizeTrip adds finishing touches to the trip plan
func (p *TripPlanner) finalizeTrip(plan *TripPlan, req TripPlanRequest) {
	// Add general travel tips
	plan.Tips = []string{
		"建议提前预订热门景点门票",
		"准备雨具以防不时之需",
		"关注天气变化适时调整行程",
		"重要物品请随身携带",
	}

	// Add location-specific tips
	if req.Location.Latitude >= 30.2 && req.Location.Latitude <= 30.3 &&
		req.Location.Longitude >= 120.1 && req.Location.Longitude <= 120.2 {
		plan.Tips = append(plan.Tips,
			"西湖景区周末人流量较大",
			"建议选择地铁等公共交通工具",
			"可以考虑购买景区联票",
		)
	}

	// Extract highlights
	plan.Overview.Highlights = make([]string, 0)
	for _, day := range plan.DailyPlans {
		if len(day.Weather.Suitable) > 0 {
			plan.Overview.Highlights = append(plan.Overview.Highlights,
				fmt.Sprintf("%s适合：%s", day.Date, join(day.Weather.Suitable, "、")))
		}
	}
}

// Helper functions

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func join(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}
	if len(slice) == 2 {
		return slice[0] + sep + slice[1]
	}
	return slice[0] + sep + slice[1] + "等"
}
