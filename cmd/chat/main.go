package main

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"

	"github.com/hewenyu/deepllm/internal/chat"
)

func createOllamaChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434", // Ollama 服务地址
		Model:   "llama2",                 // 模型名称
	})
	if err != nil {
		log.Fatalf("create ollama chat model failed: %v", err)
	}
	return chatModel
}

func main() {
	ctx := context.Background()

	// 使用模版创建messages
	log.Printf("===create messages===\n")
	messages := chat.CreateMessagesFromTemplate()
	log.Printf("messages: %+v\n\n", messages)

	// 创建llm
	log.Printf("===create llm===\n")
	cm := createOllamaChatModel(ctx)
	log.Printf("create llm success\n\n")

	log.Printf("===llm generate===\n")
	result := chat.Generate(ctx, cm, messages)
	log.Printf("result: %+v\n\n", result)

	log.Printf("===llm stream generate===\n")
	streamResult := chat.Stream(ctx, cm, messages)
	chat.ReportStream(streamResult)
}
