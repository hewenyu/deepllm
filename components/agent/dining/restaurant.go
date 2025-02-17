package dining

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/hewenyu/deepllm/internal/data"
)

const earthRadiusKm = 6371.0

// RestaurantAgent specializes in restaurant recommendations
type RestaurantAgent struct {
	store *data.Store
}

// NewRestaurantAgent creates a new restaurant recommendation agent
func NewRestaurantAgent(store *data.Store) *RestaurantAgent {
	return &RestaurantAgent{
		store: store,
	}
}

// DiningRequest represents a request for restaurant recommendations
type DiningRequest struct {
	Location    data.Location `json:"location"`
	Time        time.Time     `json:"time"`    // Dining time
	Budget      float64       `json:"budget"`  // Per person budget
	Cuisine     []string      `json:"cuisine"` // Preferred cuisine types
	PartySize   int           `json:"party_size"`
	Distance    float64       `json:"max_distance_km"` // Maximum distance in km
	Preferences []string      `json:"preferences"`     // e.g., ["安静", "景观", "茶位"]
}

// DiningRecommendation contains restaurant recommendation details
type DiningRecommendation struct {
	Restaurant     *data.Restaurant `json:"restaurant"`
	DistanceKm     float64          `json:"distance_km"`
	ReasonToVisit  []string         `json:"reasons"` // Why this restaurant is recommended
	SpecialNotes   []string         `json:"notes"`   // Special considerations
	ReservationTip string           `json:"reservation_tip"`
}

// Recommend generates restaurant recommendations based on the request
func (a *RestaurantAgent) Recommend(ctx context.Context, req DiningRequest) ([]DiningRecommendation, error) {
	// Get nearby restaurants
	nearby := a.store.FindNearbyRestaurants(req.Location, req.Distance)

	// Score and filter restaurants
	var recommendations []DiningRecommendation
	for _, r := range nearby {
		if score := a.scoreRestaurant(r, req); score > 0 {
			rec := DiningRecommendation{
				Restaurant:     &r,
				DistanceKm:     calculateDistance(req.Location, r.Coordinates),
				ReasonToVisit:  a.generateReasons(r, req),
				SpecialNotes:   a.generateNotes(r, req),
				ReservationTip: a.getReservationTip(r, req.Time, req.PartySize),
			}
			recommendations = append(recommendations, rec)
		}
	}

	// Sort by score and distance
	sort.Slice(recommendations, func(i, j int) bool {
		scoreI := a.scoreRestaurant(*recommendations[i].Restaurant, req)
		scoreJ := a.scoreRestaurant(*recommendations[j].Restaurant, req)
		if scoreI == scoreJ {
			return recommendations[i].DistanceKm < recommendations[j].DistanceKm
		}
		return scoreI > scoreJ
	})

	// Limit to top recommendations
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations, nil
}

// calculateDistance calculates the distance between two locations using the Haversine formula
func calculateDistance(a, b data.Location) float64 {
	lat1 := toRadians(a.Latitude)
	lon1 := toRadians(a.Longitude)
	lat2 := toRadians(b.Latitude)
	lon2 := toRadians(b.Longitude)

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	h := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Pow(math.Sin(dLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))

	return earthRadiusKm * c
}

// toRadians converts degrees to radians
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// scoreRestaurant scores a restaurant based on how well it matches the request
func (a *RestaurantAgent) scoreRestaurant(r data.Restaurant, req DiningRequest) float64 {
	score := 0.0

	// Check budget constraints
	avgPrice := (r.PriceRange.Min + r.PriceRange.Max) / 2
	if avgPrice > req.Budget {
		return 0 // Over budget
	}

	// Base score from price match (closer to budget = better)
	priceFactor := 1.0 - (req.Budget-avgPrice)/req.Budget
	score += priceFactor * 3

	// Cuisine type match
	if len(req.Cuisine) > 0 {
		for _, c := range req.Cuisine {
			if r.CuisineType == c {
				score += 5
				break
			}
		}
	}

	// Preference match
	for _, pref := range req.Preferences {
		if containsAny(r.Features, []string{pref}) {
			score += 2
		}
	}

	return score
}

// generateReasons generates reasons why this restaurant is recommended
func (a *RestaurantAgent) generateReasons(r data.Restaurant, req DiningRequest) []string {
	var reasons []string

	// Famous dishes
	if len(r.SignatureDishes) > 0 {
		reasons = append(reasons, "特色菜品: "+join(r.SignatureDishes, ", "))
	}

	// Features matching preferences
	var matchedFeatures []string
	for _, f := range r.Features {
		if containsAny(req.Preferences, []string{f}) {
			matchedFeatures = append(matchedFeatures, f)
		}
	}
	if len(matchedFeatures) > 0 {
		reasons = append(reasons, "符合偏好: "+join(matchedFeatures, ", "))
	}

	// Price level consideration
	if r.PriceRange.Level == "中等" || r.PriceRange.Level == "经济" {
		reasons = append(reasons, "性价比高")
	}

	return reasons
}

// generateNotes generates special notes about the restaurant
func (a *RestaurantAgent) generateNotes(r data.Restaurant, req DiningRequest) []string {
	var notes []string

	if r.ReservationNeeded {
		notes = append(notes, "建议提前预约")
	}

	// Opening hours consideration
	if r.OpeningHours.BreakTime != nil {
		notes = append(notes, "注意中午休息时间: "+
			r.OpeningHours.BreakTime.Start+"-"+
			r.OpeningHours.BreakTime.End)
	}

	return notes
}

// getReservationTip provides specific reservation advice
func (a *RestaurantAgent) getReservationTip(r data.Restaurant, time time.Time, partySize int) string {
	if !r.ReservationNeeded {
		return "无需预约"
	}

	if partySize > 6 {
		return "建议至少提前1天预约"
	}

	hour := time.Hour()
	if (hour >= 11 && hour <= 13) || (hour >= 17 && hour <= 19) {
		return "高峰时段，建议提前2小时预约"
	}

	return "建议提前预约"
}

// Helper functions

func containsAny(slice []string, targets []string) bool {
	for _, s := range slice {
		for _, t := range targets {
			if s == t {
				return true
			}
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
