package coordinator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hewenyu/deepllm/components/agent"
	"github.com/hewenyu/deepllm/components/agent/accommodation"
	"github.com/hewenyu/deepllm/components/agent/dining"
	"github.com/hewenyu/deepllm/components/agent/weather"
	"github.com/hewenyu/deepllm/internal/data"
)

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

// TripPlanner coordinates multiple agents for trip planning
type TripPlanner struct {
	store           *data.Store
	config          *agent.BaseAgentOptions
	weatherAgent    *weather.WeatherAgent
	restaurantAgent *dining.RestaurantAgent
	hotelAgent      *accommodation.HotelAgent
}

// NewTripPlanner creates a new trip planner
func NewTripPlanner(cfg *agent.BaseAgentOptions, store *data.Store) *TripPlanner {
	planner := &TripPlanner{
		store:  store,
		config: cfg,
	}

	// Initialize agents
	planner.weatherAgent = weather.NewWeatherAgent(agent.BaseAgentOptions{
		Config:      cfg.Config,
		Name:        "WeatherAdvisor",
		Description: "天气与活动建议智能体",
	}, store)

	// TODO: Initialize other agents with proper configuration
	// planner.restaurantAgent = ...
	// planner.hotelAgent = ...

	return planner
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

	if p.hotelAgent != nil {
		hotels, err := p.hotelAgent.Recommend(ctx, hotelReq)
		if err != nil {
			return nil, fmt.Errorf("hotel recommendation failed: %v", err)
		}
		if len(hotels) > 0 {
			plan.Accommodation = &hotels[0]
		}
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
	weatherReq := weather.WeatherRequest{Date: date}
	weatherResp, err := p.weatherAgent.Process(ctx, weatherReq)
	if err != nil {
		return nil, fmt.Errorf("weather advice failed: %v", err)
	}
	if advice, ok := weatherResp.Data.(*weather.WeatherAdvice); ok {
		plan.Weather = advice
	}

	// TODO: Implement restaurant recommendations and activity planning
	// Use restaurantAgent to get dining recommendations
	// Plan activities based on weather and preferences

	return plan, nil
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

	// Extract highlights from daily plans
	for _, day := range plan.DailyPlans {
		if day.Weather != nil && len(day.Weather.Suitable) > 0 {
			highlight := fmt.Sprintf("%s适合：%s", day.Date,
				strings.Join(day.Weather.Suitable, "、"))
			plan.Overview.Highlights = append(plan.Overview.Highlights, highlight)
		}
	}
}
