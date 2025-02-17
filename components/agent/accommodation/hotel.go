package accommodation

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/hewenyu/deepllm/internal/data"
)

const earthRadiusKm = 6371.0

// HotelAgent specializes in hotel recommendations
type HotelAgent struct {
	store *data.Store
}

// NewHotelAgent creates a new hotel recommendation agent
func NewHotelAgent(store *data.Store) *HotelAgent {
	return &HotelAgent{
		store: store,
	}
}

// AccommodationRequest represents a request for hotel recommendations
type AccommodationRequest struct {
	Location     data.Location `json:"location"`
	CheckIn      time.Time     `json:"check_in"`
	CheckOut     time.Time     `json:"check_out"`
	Budget       float64       `json:"budget_per_night"` // Per night budget
	GuestCount   int           `json:"guest_count"`      // Number of guests
	Distance     float64       `json:"max_distance_km"`  // Maximum distance in km
	Preferences  []string      `json:"preferences"`      // e.g., ["商务", "亲子", "度假"]
	Requirements []string      `json:"requirements"`     // e.g., ["无烟房", "双床"]
}

// HotelRecommendation contains hotel recommendation details
type HotelRecommendation struct {
	Hotel        *data.Hotel  `json:"hotel"`
	DistanceKm   float64      `json:"distance_km"`
	ReasonToBook []string     `json:"reasons"`    // Why this hotel is recommended
	RoomTypes    []RoomChoice `json:"room_types"` // Suitable room types
	SpecialNotes []string     `json:"notes"`      // Special considerations
}

// RoomChoice represents a recommended room type
type RoomChoice struct {
	Type     string   `json:"type"`
	Price    float64  `json:"price"`
	Size     float64  `json:"size_sqm"`
	Features []string `json:"features"`
	Notes    string   `json:"notes"`
}

// Recommend generates hotel recommendations based on the request
func (a *HotelAgent) Recommend(ctx context.Context, req AccommodationRequest) ([]HotelRecommendation, error) {
	// Get hotels within budget
	hotels := a.store.FindHotelsByPriceRange(0, req.Budget)

	// Score and filter hotels
	var recommendations []HotelRecommendation
	for _, h := range hotels {
		if dist := calculateDistance(req.Location, h.Coordinates); dist <= req.Distance {
			if score := a.scoreHotel(h, req); score > 0 {
				rec := HotelRecommendation{
					Hotel:        &h,
					DistanceKm:   dist,
					ReasonToBook: a.generateReasons(h, req),
					RoomTypes:    a.findSuitableRooms(h, req),
					SpecialNotes: a.generateNotes(h, req),
				}
				recommendations = append(recommendations, rec)
			}
		}
	}

	// Sort by score and distance
	sort.Slice(recommendations, func(i, j int) bool {
		scoreI := a.scoreHotel(*recommendations[i].Hotel, req)
		scoreJ := a.scoreHotel(*recommendations[j].Hotel, req)
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

// scoreHotel scores a hotel based on how well it matches the request
func (a *HotelAgent) scoreHotel(h data.Hotel, req AccommodationRequest) float64 {
	score := 0.0

	// Check if any room is within budget
	hasAffordableRoom := false
	for _, room := range h.Rooms {
		if room.Price <= req.Budget {
			hasAffordableRoom = true
			break
		}
	}
	if !hasAffordableRoom {
		return 0
	}

	// Location score (inverse of distance)
	distanceScore := 1.0 - (calculateDistance(req.Location, h.Coordinates) / req.Distance)
	score += distanceScore * 3

	// Amenity match score
	for _, pref := range req.Preferences {
		if containsAny(h.Amenities, []string{pref}) {
			score += 2
		}
	}

	// Category score
	switch h.Category {
	case "五星级":
		score += 5
	case "四星级":
		score += 4
	case "精品酒店":
		score += 3
	}

	return score
}

// findSuitableRooms finds room types that match the request criteria
func (a *HotelAgent) findSuitableRooms(h data.Hotel, req AccommodationRequest) []RoomChoice {
	var choices []RoomChoice
	for _, room := range h.Rooms {
		if room.Price <= req.Budget {
			choice := RoomChoice{
				Type:     room.Type,
				Price:    room.Price,
				Size:     room.SizeSqm,
				Features: room.Features,
			}

			// Add notes based on guest count and requirements
			if req.GuestCount > 2 && containsAny(room.Features, []string{"大床", "单床"}) {
				choice.Notes = "可能不适合当前人数"
			} else if containsAny(room.Features, req.Requirements) {
				choice.Notes = "符合需求"
			}

			choices = append(choices, choice)
		}
	}
	return choices
}

// generateReasons generates reasons why this hotel is recommended
func (a *HotelAgent) generateReasons(h data.Hotel, req AccommodationRequest) []string {
	var reasons []string

	// Location advantages
	if len(h.Transport.NearbyStations) > 0 {
		reasons = append(reasons, "交通便利: 临近"+join(h.Transport.NearbyStations, ", "))
	}

	// Matching amenities
	var matchedAmenities []string
	for _, amenity := range h.Amenities {
		if containsAny(req.Preferences, []string{amenity}) {
			matchedAmenities = append(matchedAmenities, amenity)
		}
	}
	if len(matchedAmenities) > 0 {
		reasons = append(reasons, "设施齐全: "+join(matchedAmenities, ", "))
	}

	// Price consideration
	if h.PriceRange.Level == "经济" || h.PriceRange.Level == "中等" {
		reasons = append(reasons, "性价比高")
	}

	return reasons
}

// generateNotes generates special notes about the hotel
func (a *HotelAgent) generateNotes(h data.Hotel, req AccommodationRequest) []string {
	var notes []string

	// Transportation notes
	notes = append(notes, "距离机场: "+h.Transport.FromAirport.TaxiTime)

	// Price fluctuation warning
	if h.PriceRange.Notes != "" {
		notes = append(notes, "价格提示: "+h.PriceRange.Notes)
	}

	return notes
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
