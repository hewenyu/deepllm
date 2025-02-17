package mock

import (
	"context"
)

// Message represents a chat message
type Message struct {
	Role       string                 `json:"role"`
	Content    string                 `json:"content"`
	ToolCalls  []ToolCall             `json:"tool_calls,omitempty"`
	ToolResult map[string]interface{} `json:"tool_result,omitempty"`
}

// ToolCall represents a tool call request
type ToolCall struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Args    map[string]interface{} `json:"args"`
	Results map[string]interface{} `json:"results,omitempty"`
}

// ChatModel represents a chat model interface
type ChatModel interface {
	Generate(ctx context.Context, messages []*Message) (*Message, error)
	BindTools(tools []Tool) error
}

// Tool represents a tool interface
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error)
}

// Runnable represents a runnable interface
type Runnable[I any, O any] interface {
	Invoke(ctx context.Context, input I) (O, error)
}

// MockTool is a mock implementation of Tool
type MockTool struct {
	name        string
	description string
	handler     func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error)
}

// NewMockTool creates a new MockTool instance
func NewMockTool(name, description string, handler func(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error)) *MockTool {
	return &MockTool{
		name:        name,
		description: description,
		handler:     handler,
	}
}

// Name returns the tool's name
func (t *MockTool) Name() string {
	return t.name
}

// Description returns the tool's description
func (t *MockTool) Description() string {
	return t.description
}

// Execute executes the tool with the given arguments
func (t *MockTool) Execute(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
	return t.handler(ctx, args)
}
