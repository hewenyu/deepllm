package ollama

import (
	"bytes"
	"context"
	"deepllm/components/mock"
	"encoding/json"
	"fmt"
	"net/http"
)

// ChatModel is an implementation of mock.ChatModel using Ollama
type ChatModel struct {
	baseURL string
	model   string
	tools   []mock.Tool
}

// NewChatModel creates a new ChatModel instance
func NewChatModel(baseURL, model string) *ChatModel {
	return &ChatModel{
		baseURL: baseURL,
		model:   model,
	}
}

// Request represents a request to the Ollama API
type Request struct {
	Model    string                   `json:"model"`
	Messages []Message                `json:"messages"`
	Stream   bool                     `json:"stream"`
	Options  map[string]interface{}   `json:"options,omitempty"`
	Tools    []map[string]interface{} `json:"tools,omitempty"`
}

// Message represents a message in the Ollama API
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents a response from the Ollama API
type Response struct {
	Model     string          `json:"model"`
	Response  string          `json:"response"`
	Done      bool            `json:"done"`
	ToolCalls []mock.ToolCall `json:"tool_calls,omitempty"`
	Context   []interface{}   `json:"context,omitempty"`
}

// Generate generates a response from the Ollama model
func (m *ChatModel) Generate(ctx context.Context, messages []*mock.Message) (*mock.Message, error) {
	// Convert messages to Ollama format
	ollamaMessages := make([]Message, len(messages))
	for i, msg := range messages {
		ollamaMessages[i] = Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Prepare request
	reqBody := Request{
		Model:    m.model,
		Messages: ollamaMessages,
		Stream:   false,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/api/chat", m.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var ollamaResp Response
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Convert response to Message format
	return &mock.Message{
		Role:      "assistant",
		Content:   ollamaResp.Response,
		ToolCalls: ollamaResp.ToolCalls,
	}, nil
}

// BindTools binds tools to the chat model
func (m *ChatModel) BindTools(tools []mock.Tool) error {
	// Convert tools to Ollama format
	ollamaTools := make([]map[string]interface{}, len(tools))
	for i, tool := range tools {
		ollamaTools[i] = map[string]interface{}{
			"name":        tool.Name(),
			"description": tool.Description(),
		}
	}

	// Store tools for later use
	m.tools = tools
	return nil
}

// ExecuteTool executes a tool by name with the given arguments
func (m *ChatModel) ExecuteTool(ctx context.Context, name string, args map[string]interface{}) (map[string]interface{}, error) {
	// Find the tool by name
	for _, tool := range m.tools {
		if tool.Name() == name {
			return tool.Execute(ctx, args)
		}
	}
	return nil, fmt.Errorf("tool not found: %s", name)
}
