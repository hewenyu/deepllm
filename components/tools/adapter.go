package tools

import (
	"context"
	"deepllm/components/mock"
	"deepllm/internal/data"
	"encoding/json"
	"fmt"
)

// CreateTourismTools 创建旅游相关的工具集
func CreateTourismTools(dataQuery *data.DataQuery) []mock.Tool {
	tourismTools := NewTourismTools(dataQuery)
	return []mock.Tool{
		NewSearchAttractionsTool(tourismTools),
		NewSearchRestaurantsTool(tourismTools),
		NewSearchHotelsTool(tourismTools),
		NewGetWeatherTool(tourismTools),
	}
}

// BaseTool 提供基础工具实现
type BaseTool struct {
	name        string
	description string
	handler     func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error)
}

func (t *BaseTool) Name() string {
	return t.name
}

func (t *BaseTool) Description() string {
	return t.description
}

func (t *BaseTool) Execute(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
	return t.handler(ctx, args)
}

// NewSearchAttractionsTool 创建景点搜索工具
func NewSearchAttractionsTool(t *TourismTools) mock.Tool {
	return &BaseTool{
		name:        "search_attractions",
		description: "搜索景点信息，支持按位置、类别、价格等条件筛选",
		handler: func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
			// 解析参数
			params := &AttractionQueryParams{}
			if loc, ok := args["location"].(map[string]interface{}); ok {
				params.Location = &data.Location{
					Name:      loc["name"].(string),
					Latitude:  loc["latitude"].(float64),
					Longitude: loc["longitude"].(float64),
				}
			}
			if radius, ok := args["radius"].(float64); ok {
				params.Radius = radius
			}
			if categories, ok := args["categories"].([]interface{}); ok {
				for _, cat := range categories {
					params.Categories = append(params.Categories, cat.(string))
				}
			}
			if maxPrice, ok := args["max_price"].(float64); ok {
				params.MaxPrice = maxPrice
			}

			// 执行搜索
			result, err := t.SearchAttractions(ctx, params)
			if err != nil {
				return nil, err
			}

			// 解析JSON结果
			var attractions []data.Attraction
			if err := json.Unmarshal([]byte(result), &attractions); err != nil {
				return nil, fmt.Errorf("解析结果失败: %v", err)
			}

			return map[string]interface{}{
				"attractions": attractions,
			}, nil
		},
	}
}

// NewSearchRestaurantsTool 创建餐厅搜索工具
func NewSearchRestaurantsTool(t *TourismTools) mock.Tool {
	return &BaseTool{
		name:        "search_restaurants",
		description: "搜索餐厅信息，支持按位置、菜系、价格区间等条件筛选",
		handler: func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
			// 解析参数
			params := &RestaurantQueryParams{}
			if loc, ok := args["location"].(map[string]interface{}); ok {
				params.Location = &data.Location{
					Name:      loc["name"].(string),
					Latitude:  loc["latitude"].(float64),
					Longitude: loc["longitude"].(float64),
				}
			}
			if radius, ok := args["radius"].(float64); ok {
				params.Radius = radius
			}
			if cuisines, ok := args["cuisines"].([]interface{}); ok {
				for _, cuisine := range cuisines {
					params.Cuisines = append(params.Cuisines, cuisine.(string))
				}
			}
			if priceRange, ok := args["price_range"].(string); ok {
				params.PriceRange = priceRange
			}

			// 执行搜索
			result, err := t.SearchRestaurants(ctx, params)
			if err != nil {
				return nil, err
			}

			// 解析JSON结果
			var restaurants []data.Restaurant
			if err := json.Unmarshal([]byte(result), &restaurants); err != nil {
				return nil, fmt.Errorf("解析结果失败: %v", err)
			}

			return map[string]interface{}{
				"restaurants": restaurants,
			}, nil
		},
	}
}

// NewSearchHotelsTool 创建酒店搜索工具
func NewSearchHotelsTool(t *TourismTools) mock.Tool {
	return &BaseTool{
		name:        "search_hotels",
		description: "搜索酒店信息，支持按位置、星级、价格、设施等条件筛选",
		handler: func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
			// 解析参数
			params := &HotelQueryParams{}
			if loc, ok := args["location"].(map[string]interface{}); ok {
				params.Location = &data.Location{
					Name:      loc["name"].(string),
					Latitude:  loc["latitude"].(float64),
					Longitude: loc["longitude"].(float64),
				}
			}
			if radius, ok := args["radius"].(float64); ok {
				params.Radius = radius
			}
			if minStars, ok := args["min_stars"].(float64); ok {
				params.MinStars = int(minStars)
			}
			if maxPrice, ok := args["max_price"].(float64); ok {
				params.MaxPrice = maxPrice
			}
			if amenities, ok := args["required_amenities"].([]interface{}); ok {
				for _, amen := range amenities {
					params.RequiredAmens = append(params.RequiredAmens, amen.(string))
				}
			}

			// 执行搜索
			result, err := t.SearchHotels(ctx, params)
			if err != nil {
				return nil, err
			}

			// 解析JSON结果
			var hotels []data.Hotel
			if err := json.Unmarshal([]byte(result), &hotels); err != nil {
				return nil, fmt.Errorf("解析结果失败: %v", err)
			}

			return map[string]interface{}{
				"hotels": hotels,
			}, nil
		},
	}
}

// NewGetWeatherTool 创建天气查询工具
func NewGetWeatherTool(t *TourismTools) mock.Tool {
	return &BaseTool{
		name:        "get_weather",
		description: "获取指定日期和位置的天气信息",
		handler: func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
			// 解析参数
			params := &WeatherQueryParams{}
			if loc, ok := args["location"].(map[string]interface{}); ok {
				params.Location = &data.Location{
					Name:      loc["name"].(string),
					Latitude:  loc["latitude"].(float64),
					Longitude: loc["longitude"].(float64),
				}
			}
			if date, ok := args["date"].(string); ok {
				params.Date = date
			}

			// 执行查询
			result, err := t.GetWeather(ctx, params)
			if err != nil {
				return nil, err
			}

			// 解析JSON结果
			var weather data.Weather
			if err := json.Unmarshal([]byte(result), &weather); err != nil {
				return nil, fmt.Errorf("解析结果失败: %v", err)
			}

			return map[string]interface{}{
				"weather": weather,
			}, nil
		},
	}
}
