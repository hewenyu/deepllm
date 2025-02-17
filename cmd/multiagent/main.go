package main

import (
	"context"
	"deepllm/components/agent/coordinator"
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

	// Create tools
	tools := []mock.Tool{
		mock.NewMockTool(
			"search_attractions",
			"搜索附近的景点",
			func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"status": "success",
					"data":   "模拟景点数据",
				}, nil
			},
		),
		mock.NewMockTool(
			"search_restaurants",
			"搜索附近的餐厅",
			func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"status": "success",
					"data":   "模拟餐厅数据",
				}, nil
			},
		),
		mock.NewMockTool(
			"search_hotels",
			"搜索附近的酒店",
			func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"status": "success",
					"data":   "模拟酒店数据",
				}, nil
			},
		),
		mock.NewMockTool(
			"get_weather",
			"获取天气预报",
			func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"status": "success",
					"data":   "模拟天气数据",
				}, nil
			},
		),
	}

	// Create coordinator agent
	coordinatorAgent := coordinator.NewCoordinatorAgent(chatModel, tools, dataQuery)

	// Create sample trip request
	request := &data.TripPlanRequest{
		StartDate: time.Now().Add(24 * time.Hour),
		EndDate:   time.Now().Add(72 * time.Hour),
		Location: data.Location{
			Name:      "西湖",
			Latitude:  30.2587,
			Longitude: 120.1315,
		},
		Budget: struct {
			Total    float64 `json:"total"`
			Hotel    float64 `json:"hotel"`
			Food     float64 `json:"food"`
			Activity float64 `json:"activity"`
		}{
			Total:    5000,
			Hotel:    1000,
			Food:     300,
			Activity: 200,
		},
		Preferences: struct {
			Activities []string `json:"activities"`
			Cuisine    []string `json:"cuisine"`
			Hotel      []string `json:"hotel"`
		}{
			Activities: []string{"观光", "购物", "文化"},
			Cuisine:    []string{"本地菜", "海鲜"},
			Hotel:      []string{"豪华", "湖景"},
		},
		PartySize:    2,
		Requirements: []string{"无障碍设施"},
	}

	// Process the request
	ctx := context.Background()
	result, err := coordinatorAgent.Process(ctx, request)
	if err != nil {
		log.Fatalf("处理请求失败: %v", err)
	}

	// Print the result
	fmt.Printf("行程规划:\n%s\n", result.(*mock.Message).Content)
}
