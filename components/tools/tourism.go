package tools

import (
	"context"
	"deepllm/internal/data"
	"encoding/json"
	"fmt"
	"time"
)

// 景点查询参数
type AttractionQueryParams struct {
	Location   *data.Location `json:"location,omitempty" jsonschema:"description=查询位置"`
	Radius     float64        `json:"radius,omitempty" jsonschema:"description=搜索半径（公里）"`
	Categories []string       `json:"categories,omitempty" jsonschema:"description=景点类别，如：自然风光、人文景观等"`
	MaxPrice   float64        `json:"max_price,omitempty" jsonschema:"description=最高门票价格"`
}

// 餐厅查询参数
type RestaurantQueryParams struct {
	Location   *data.Location `json:"location,omitempty" jsonschema:"description=查询位置"`
	Radius     float64        `json:"radius,omitempty" jsonschema:"description=搜索半径（公里）"`
	Cuisines   []string       `json:"cuisines,omitempty" jsonschema:"description=菜系类型，如：杭帮菜、海鲜等"`
	PriceRange string         `json:"price_range,omitempty" jsonschema:"description=价格区间，$-$$$$"`
}

// 酒店查询参数
type HotelQueryParams struct {
	Location      *data.Location `json:"location,omitempty" jsonschema:"description=查询位置"`
	Radius        float64        `json:"radius,omitempty" jsonschema:"description=搜索半径（公里）"`
	MinStars      int            `json:"min_stars,omitempty" jsonschema:"description=最低星级"`
	MaxPrice      float64        `json:"max_price,omitempty" jsonschema:"description=最高房价/晚"`
	RequiredAmens []string       `json:"required_amenities,omitempty" jsonschema:"description=必需设施，如：游泳池、健身房等"`
}

// 天气查询参数
type WeatherQueryParams struct {
	Location *data.Location `json:"location,omitempty" jsonschema:"description=查询位置"`
	Date     string         `json:"date,omitempty" jsonschema:"description=查询日期，格式：2024-02-18"`
}

// TourismTools 提供旅游相关的工具集
type TourismTools struct {
	dataQuery *data.DataQuery
}

// NewTourismTools 创建新的旅游工具集
func NewTourismTools(dataQuery *data.DataQuery) *TourismTools {
	return &TourismTools{
		dataQuery: dataQuery,
	}
}

// SearchAttractions 搜索景点
func (t *TourismTools) SearchAttractions(ctx context.Context, params *AttractionQueryParams) (string, error) {
	attractions, err := t.dataQuery.Loader.LoadAttractions()
	if err != nil {
		return "", fmt.Errorf("加载景点数据失败: %v", err)
	}

	// 按位置筛选
	var filtered []data.Attraction
	if params.Location != nil {
		filtered = data.FindNearbyAttractions(attractions, *params.Location, params.Radius)
	} else {
		filtered = attractions
	}

	// 按类别筛选
	if len(params.Categories) > 0 {
		filtered = t.dataQuery.FilterAttractionsByPreferences(filtered, params.Categories)
	}

	// 按价格筛选
	if params.MaxPrice > 0 {
		filtered = t.dataQuery.FilterByBudget(filtered, params.MaxPrice)
	}

	// 按评分排序
	sorted := t.dataQuery.SortByRating(filtered)

	// 转换为JSON
	result, err := json.Marshal(sorted)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %v", err)
	}

	return string(result), nil
}

// SearchRestaurants 搜索餐厅
func (t *TourismTools) SearchRestaurants(ctx context.Context, params *RestaurantQueryParams) (string, error) {
	restaurants, err := t.dataQuery.Loader.LoadRestaurants()
	if err != nil {
		return "", fmt.Errorf("加载餐厅数据失败: %v", err)
	}

	// 按位置筛选
	var filtered []data.Restaurant
	if params.Location != nil {
		filtered = data.FindNearbyRestaurants(restaurants, *params.Location, params.Radius)
	} else {
		filtered = restaurants
	}

	// 按菜系筛选
	if len(params.Cuisines) > 0 {
		filtered = t.dataQuery.FilterRestaurantsByPreferences(filtered, params.Cuisines)
	}

	// 转换为JSON
	result, err := json.Marshal(filtered)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %v", err)
	}

	return string(result), nil
}

// SearchHotels 搜索酒店
func (t *TourismTools) SearchHotels(ctx context.Context, params *HotelQueryParams) (string, error) {
	hotels, err := t.dataQuery.Loader.LoadHotels()
	if err != nil {
		return "", fmt.Errorf("加载酒店数据失败: %v", err)
	}

	// 按位置筛选
	var filtered []data.Hotel
	if params.Location != nil {
		filtered = data.FindNearbyHotels(hotels, *params.Location, params.Radius)
	} else {
		filtered = hotels
	}

	// 按星级和价格筛选
	var result []data.Hotel
	for _, hotel := range filtered {
		if params.MinStars > 0 && hotel.Stars < params.MinStars {
			continue
		}
		if params.MaxPrice > 0 && hotel.PricePerNight > params.MaxPrice {
			continue
		}
		if len(params.RequiredAmens) > 0 {
			hasAllAmens := true
			for _, required := range params.RequiredAmens {
				found := false
				for _, amen := range hotel.Amenities {
					if amen == required {
						found = true
						break
					}
				}
				if !found {
					hasAllAmens = false
					break
				}
			}
			if !hasAllAmens {
				continue
			}
		}
		result = append(result, hotel)
	}

	// 转换为JSON
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %v", err)
	}

	return string(jsonResult), nil
}

// GetWeather 获取天气信息
func (t *TourismTools) GetWeather(ctx context.Context, params *WeatherQueryParams) (string, error) {
	weatherData, err := t.dataQuery.Loader.LoadWeather()
	if err != nil {
		return "", fmt.Errorf("加载天气数据失败: %v", err)
	}

	date, err := time.Parse("2006-01-02", params.Date)
	if err != nil {
		return "", fmt.Errorf("日期格式错误: %v", err)
	}

	weather, found := t.dataQuery.GetWeatherForDate(weatherData, date, *params.Location)
	if !found {
		return "", fmt.Errorf("未找到指定日期的天气数据")
	}

	// 转换为JSON
	result, err := json.Marshal(weather)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %v", err)
	}

	return string(result), nil
}
