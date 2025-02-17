package agent

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/hewenyu/deepllm/internal/data"
)

// OllamaAgent represents an agent powered by local Ollama model
type OllamaAgent struct {
	chatModel model.ChatModel
	chain     compose.Runnable[[]*schema.Message, *schema.Message]
	store     *data.Store
}

// NewOllamaAgent creates a new Ollama-powered agent
func NewOllamaAgent(ctx context.Context, baseURL string, modelName string, store *data.Store, tools []tool.BaseTool) (*OllamaAgent, error) {
	// Initialize chat model
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: baseURL,   // Ollama service address
		Model:   modelName, // Model name
	})
	if err != nil {
		return nil, err
	}

	// Bind tools to chat model if tools are provided
	if len(tools) > 0 {
		toolInfos := make([]*schema.ToolInfo, 0, len(tools))
		for _, t := range tools {
			info, err := t.Info(ctx)
			if err != nil {
				log.Printf("Failed to get tool info: %v", err)
				continue
			}
			toolInfos = append(toolInfos, info)
		}
		if err := chatModel.BindTools(toolInfos); err != nil {
			return nil, err
		}
	}

	// Create tools node if tools are provided
	var toolsNode compose.Runnable[[]*schema.Message, *schema.Message]
	if len(tools) > 0 {
		tn, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
			Tools: tools,
		})
		if err != nil {
			return nil, err
		}
		toolsNode = tn
	}

	// Build the chain
	chain := compose.NewChain[[]*schema.Message, *schema.Message]()
	chain.AppendChatModel(chatModel)
	if toolsNode != nil {
		chain.AppendToolsNode(toolsNode)
	}

	// Compile the chain
	runnable, err := chain.Compile(ctx)
	if err != nil {
		return nil, err
	}

	return &OllamaAgent{
		chatModel: chatModel,
		chain:     runnable,
		store:     store,
	}, nil
}

// Generate generates a response for the given messages
func (a *OllamaAgent) Generate(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	return a.chain.Invoke(ctx, messages)
}

// Stream generates a streaming response for the given messages
func (a *OllamaAgent) Stream(ctx context.Context, messages []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	return a.chain.Stream(ctx, messages)
}

// GetStore returns the data store
func (a *OllamaAgent) GetStore() *data.Store {
	return a.store
}
