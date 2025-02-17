package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/hewenyu/deepllm/components/agent"
	"github.com/hewenyu/deepllm/components/agent/tools"
	"github.com/hewenyu/deepllm/internal/data"
)

func main() {
	ctx := context.Background()

	// Initialize data store
	store := data.NewStore("./data")

	// Load all data
	if err := store.LoadAll(ctx); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	// Create tools
	weatherTool, err := tools.NewWeatherTool(store)
	if err != nil {
		log.Fatalf("Failed to create weather tool: %v", err)
	}

	attractionTool, err := tools.NewAttractionTool(store)
	if err != nil {
		log.Fatalf("Failed to create attraction tool: %v", err)
	}

	restaurantTool, err := tools.NewRestaurantTool(store)
	if err != nil {
		log.Fatalf("Failed to create restaurant tool: %v", err)
	}

	// Create a new Ollama agent with tools
	ollamaAgent, err := agent.NewOllamaAgent(
		ctx,
		"http://localhost:11434", // Ollama service address
		"deepseek-r1:14b",        // Model name
		store,                    // Data store
		[]tool.BaseTool{ // Add tools
			weatherTool,
			attractionTool,
			restaurantTool,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Ollama agent: %v", err)
	}

	// Prepare messages for a complex query
	messages := []*schema.Message{
		schema.SystemMessage(`你是一个专业的旅游助手。你可以使用以下工具来帮助用户：
1. get_weather - 获取天气预报
2. search_attractions - 搜索景点（可以通过区域ID或位置搜索）
3. search_restaurants - 搜索餐厅（可以通过区域ID、位置或菜系搜索）

请根据用户的问题，合理使用这些工具来提供专业的建议。`),
		schema.UserMessage("我在西湖附近，想知道今天的天气，以及附近2公里有什么好玩的地方和餐厅推荐。"),
	}

	// Generate response
	fmt.Println("Generating response...")
	response, err := ollamaAgent.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	fmt.Printf("\nAssistant: %s\n", response.Content)

	// Try another query with streaming response
	fmt.Println("\nStreaming response...")
	messages = append(messages, schema.UserMessage("这些景点里面，哪些适合在雨天游玩？另外推荐一些杭帮菜餐厅。"))

	stream, err := ollamaAgent.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.Close()

	fmt.Print("Assistant: ")
	for {
		chunk, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Print(chunk.Content)
	}
	fmt.Println()
}
