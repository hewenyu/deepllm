package main

import (
	"context"
	"deepllm/components/mock"
	"deepllm/components/ollama"
	"deepllm/components/tools"
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

	// Create tourism tools
	tourismTools := tools.CreateTourismTools(dataQuery)

	// Create system prompt
	systemPrompt := `你是一位专业的杭州旅游助手，可以为游客提供景点、美食、住宿和天气等方面的建议。
你可以使用以下工具来帮助回答问题：
- search_attractions: 搜索景点信息
- search_restaurants: 搜索餐厅信息
- search_hotels: 搜索酒店信息
- get_weather: 查询天气信息

请根据用户的需求，合理使用这些工具来提供专业的建议。回答要详细、准确，并注意以下几点：
1. 推荐时要考虑位置、价格、评分等因素
2. 解释推荐的理由，帮助用户做出选择
3. 如果天气不好，要提供相应的替代建议
4. 注意不同景点的开放时间
5. 所有回答使用中文`

	// Create user prompt
	userPrompt := `我计划明天去西湖游玩，想知道：
1. 天气情况如何？
2. 有哪些值得去的景点？
3. 中午可以在哪里吃饭？
请给我一些建议。`

	// Create messages
	messages := []*mock.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	// Bind tools to chat model
	chatModel.BindTools(tourismTools)

	// Call the model
	ctx := context.Background()
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("调用模型失败: %v", err)
	}

	// Output result
	fmt.Printf("\n=== 旅游助手回答 ===\n")
	if response != nil && response.Content != "" {
		fmt.Printf("%s\n", response.Content)
	} else {
		fmt.Println("抱歉，助手暂时无法回答您的问题。")
	}

	// 数据查询演示
	fmt.Println("\n=== 数据查询演示 ===")

	// 加载并筛选景点
	attractions, err := dataLoader.LoadAttractions()
	if err != nil {
		log.Fatalf("加载景点数据失败: %v", err)
	}

	location := data.Location{
		Name:      "西湖",
		Latitude:  30.2587,
		Longitude: 120.1315,
	}

	// 扩大搜索范围到5公里
	nearbyAttractions := data.FindNearbyAttractions(attractions, location, 5.0)
	fmt.Printf("\n【附近景点】\n")
	fmt.Printf("在西湖5公里范围内找到%d个景点\n", len(nearbyAttractions))

	// 按偏好筛选
	preferences := []string{"自然风光", "人文景观", "历史文化"}
	filteredAttractions := dataQuery.FilterAttractionsByPreferences(nearbyAttractions, preferences)
	fmt.Printf("\n【推荐景点】\n")
	fmt.Printf("符合偏好的景点数量：%d（偏好：%s）\n", len(filteredAttractions), "自然风光、人文景观、历史文化")

	// 按评分排序并展示详情
	sortedAttractions := dataQuery.SortByRating(filteredAttractions)
	fmt.Println("\n【景点排名】（按评分排序）")
	for i, attraction := range sortedAttractions {
		if i >= 5 { // 显示前5个景点
			break
		}
		fmt.Printf("%d. %s\n   评分: %.1f  价格: %.0f元\n   类别: %v\n   描述: %s\n\n",
			i+1,
			attraction.Name,
			attraction.Rating,
			attraction.Price,
			attraction.Category,
			attraction.Description,
		)
	}

	// 加载并显示天气信息
	weatherData, err := dataLoader.LoadWeather()
	if err != nil {
		log.Fatalf("加载天气数据失败: %v", err)
	}

	fmt.Println("【天气信息】")
	today := time.Now()
	weather, found := dataQuery.GetWeatherForDate(weatherData, today, location)
	if found {
		fmt.Printf("西湖今日天气:\n")
		fmt.Printf("- 气温: %.1f°C - %.1f°C\n", weather.Temperature.Min, weather.Temperature.Max)
		fmt.Printf("- 天气: %s\n", translateWeatherCondition(weather.Condition))
		fmt.Printf("- 湿度: %.1f%%\n", weather.Humidity)
		fmt.Printf("- 风速: %.1f公里/小时\n", weather.WindSpeed)
		if weather.Precipitation > 0 {
			fmt.Printf("- 降水量: %.1f毫米\n", weather.Precipitation)
		}
	} else {
		fmt.Println("暂无今日天气数据")
		// 尝试获取最近的天气数据
		for _, w := range weatherData {
			fmt.Printf("\n最近天气预报 (%s):\n", w.Date.Format("2006-01-02"))
			fmt.Printf("- 气温: %.1f°C - %.1f°C\n", w.Temperature.Min, w.Temperature.Max)
			fmt.Printf("- 天气: %s\n", translateWeatherCondition(w.Condition))
			fmt.Printf("- 湿度: %.1f%%\n", w.Humidity)
			fmt.Printf("- 风速: %.1f公里/小时\n", w.WindSpeed)
			if w.Precipitation > 0 {
				fmt.Printf("- 降水量: %.1f毫米\n", w.Precipitation)
			}
			break
		}
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
