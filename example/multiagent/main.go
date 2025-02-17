package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hewenyu/deepllm/components/agent/coordinator"
	"github.com/hewenyu/deepllm/internal/data"
)

func main() {
	// Initialize data store
	store := data.NewStore("./data")

	// Load all data
	ctx := context.Background()
	if err := store.LoadAll(ctx); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	// Create trip planner
	planner := coordinator.NewTripPlanner(store)

	// Create sample trip request
	req := coordinator.TripPlanRequest{
		StartDate: time.Now().AddDate(0, 0, 1), // Tomorrow
		EndDate:   time.Now().AddDate(0, 0, 3), // 3 days trip
		Location: data.Location{
			Latitude:  30.2587,
			Longitude: 120.1485,
		},
		Budget: struct {
			Total    float64 `json:"total"`
			Hotel    float64 `json:"hotel"`
			Food     float64 `json:"food"`
			Activity float64 `json:"activity"`
		}{
			Total:    3000,
			Hotel:    800,
			Food:     300,
			Activity: 200,
		},
		Preferences: struct {
			Activities []string `json:"activities"`
			Cuisine    []string `json:"cuisine"`
			Hotel      []string `json:"hotel"`
		}{
			Activities: []string{"自然景观", "文化古迹", "购物"},
			Cuisine:    []string{"杭帮菜", "创意菜"},
			Hotel:      []string{"商务", "休闲"},
		},
		PartySize:    2,
		Requirements: []string{"无烟房", "双床"},
	}

	// Generate trip plan
	plan, err := planner.Plan(ctx, req)
	if err != nil {
		log.Fatalf("Failed to generate plan: %v", err)
	}

	// Print trip overview
	fmt.Printf("=== 行程概览 ===\n")
	fmt.Printf("行程天数: %d天\n", plan.Overview.Duration)
	fmt.Printf("亮点推荐:\n")
	for _, highlight := range plan.Overview.Highlights {
		fmt.Printf("- %s\n", highlight)
	}

	// Print accommodation details
	if plan.Accommodation != nil {
		fmt.Printf("\n=== 住宿安排 ===\n")
		fmt.Printf("酒店: %s (%s)\n", plan.Accommodation.Hotel.Name, plan.Accommodation.Hotel.Category)
		fmt.Printf("地址: %s\n", plan.Accommodation.Hotel.Contact.Address)
		fmt.Printf("推荐理由:\n")
		for _, reason := range plan.Accommodation.ReasonToBook {
			fmt.Printf("- %s\n", reason)
		}
	}

	// Print daily plans
	fmt.Printf("\n=== 每日行程 ===\n")
	for _, day := range plan.DailyPlans {
		fmt.Printf("\n%s:\n", day.Date)

		// Weather info
		if day.Weather != nil {
			fmt.Printf("天气情况:\n")
			if len(day.Weather.Suitable) > 0 {
				fmt.Printf("- 适合活动: %v\n", day.Weather.Suitable)
			}
			if len(day.Weather.Unsuitable) > 0 {
				fmt.Printf("- 不宜活动: %v\n", day.Weather.Unsuitable)
			}
			if len(day.Weather.Precautions) > 0 {
				fmt.Printf("- 注意事项: %v\n", day.Weather.Precautions)
			}
			if len(day.Weather.IndoorOptions) > 0 {
				fmt.Printf("- 室内选项: %v\n", day.Weather.IndoorOptions)
			}
			if len(day.Weather.OutdoorOptions) > 0 {
				fmt.Printf("- 户外选项: %v\n", day.Weather.OutdoorOptions)
			}
		}

		// Activities
		fmt.Printf("\n行程安排:\n")
		for _, activity := range day.Activities {
			fmt.Printf("- %s (%s): %d分钟\n", activity.Time, activity.Type, activity.Duration)
			if len(activity.Notes) > 0 {
				fmt.Printf("  提示: %v\n", activity.Notes)
			}
		}

		// Dining
		if len(day.Dining) > 0 {
			fmt.Printf("\n用餐推荐:\n")
			for _, dining := range day.Dining {
				fmt.Printf("- %s (%s)\n", dining.Restaurant.Name, dining.Restaurant.CuisineType)
				fmt.Printf("  特色: %v\n", dining.ReasonToVisit)
				if len(dining.SpecialNotes) > 0 {
					fmt.Printf("  提示: %v\n", dining.SpecialNotes)
				}
			}
		}
	}

	// Print travel tips
	fmt.Printf("\n=== 出行贴士 ===\n")
	for _, tip := range plan.Tips {
		fmt.Printf("- %s\n", tip)
	}

	// Save plan to file
	planJSON, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal plan: %v", err)
	}
	if err := os.WriteFile("trip_plan.json", planJSON, 0644); err != nil {
		log.Fatalf("Failed to save plan: %v", err)
	}
	fmt.Printf("\n行程已保存至 trip_plan.json\n")
}
