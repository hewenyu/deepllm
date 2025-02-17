package main

import (
	"context"
	"deepllm/components/mock"
	"deepllm/components/ollama"
	"deepllm/internal/data"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Initialize data loader and query
	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "./data"
	}
	dataLoader := data.NewDataLoader(dataPath)
	dataQuery := data.NewDataQuery(dataLoader)

	// Initialize Ollama chat model
	ollamaBaseURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaBaseURL == "" {
		ollamaBaseURL = "http://localhost:11434"
	}
	ollamaModel := os.Getenv("OLLAMA_MODEL")
	if ollamaModel == "" {
		ollamaModel = "deepseek-r1:14b"
	}
	chatModel := ollama.NewChatModel(ollamaBaseURL, ollamaModel)

	// Create a simple tool
	weatherTool := mock.NewMockTool(
		"get_weather",
		"获取指定位置的天气预报",
		func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
			// Mock implementation
			return map[string]interface{}{
				"temperature": 25.0,
				"condition":   "晴朗",
				"humidity":    60.0,
			}, nil
		},
	)

	// Create a simple prompt
	messages := []*mock.Message{
		{
			Role:    "system",
			Content: "你是一位专业的杭州旅游助手，可以为游客提供天气、景点、美食等方面的建议。请用中文回答用户的问题。",
		},
		{
			Role:    "user",
			Content: "西湖今天天气怎么样？适合去划船吗？",
		},
	}

	// Call the model
	ctx := context.Background()
	chatModel.BindTools([]mock.Tool{weatherTool})
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("调用模型失败: %v", err)
	}

	fmt.Printf("助手回答:\n%s\n", response.Content)

	// Demonstrate data querying
	fmt.Println("\n数据查询演示:")

	// Load and filter attractions
	attractions, err := dataLoader.LoadAttractions()
	if err != nil {
		log.Fatalf("加载景点数据失败: %v", err)
	}

	location := data.Location{
		Name:      "西湖",
		Latitude:  30.2587,
		Longitude: 120.1315,
	}

	nearbyAttractions := data.FindNearbyAttractions(attractions, location, 2.0) // Within 2km
	fmt.Printf("\n在西湖2公里范围内找到%d个景点\n", len(nearbyAttractions))

	// Filter by preferences
	preferences := []string{"cultural", "historical"}
	filteredAttractions := dataQuery.FilterAttractionsByPreferences(nearbyAttractions, preferences)
	fmt.Printf("\n找到%d个符合偏好的景点（文化、历史）\n", len(filteredAttractions))

	// Sort by rating
	sortedAttractions := dataQuery.SortByRating(filteredAttractions)
	fmt.Println("\n评分最高的景点:")
	for i, attraction := range sortedAttractions {
		if i >= 3 {
			break
		}
		fmt.Printf("%d. %s (评分: %.1f)\n", i+1, attraction.Name, attraction.Rating)
	}

	// Demonstrate weather querying
	weatherData, err := dataLoader.LoadWeather()
	if err != nil {
		log.Fatalf("加载天气数据失败: %v", err)
	}

	weather, found := dataQuery.GetWeatherForDate(weatherData, time.Now(), location)
	if found {
		fmt.Printf("\n%s当前天气:\n", location.Name)
		fmt.Printf("气温: %.1f°C - %.1f°C\n", weather.Temperature.Min, weather.Temperature.Max)
		fmt.Printf("天气状况: %s\n", translateWeatherCondition(weather.Condition))
		fmt.Printf("湿度: %.1f%%\n", weather.Humidity)
		fmt.Printf("风速: %.1f公里/小时\n", weather.WindSpeed)
	} else {
		fmt.Println("\n当前日期没有可用的天气数据")
	}
}

// translateWeatherCondition 将英文天气状况翻译为中文
func translateWeatherCondition(condition string) string {
	translations := map[string]string{
		"sunny":         "晴朗",
		"partly cloudy": "多云",
		"cloudy":        "阴天",
		"light rain":    "小雨",
		"rain":          "雨天",
		"heavy rain":    "大雨",
		"thunderstorm":  "雷雨",
		"snow":          "雪",
		"foggy":         "雾",
	}

	if translation, ok := translations[condition]; ok {
		return translation
	}
	return condition
}
