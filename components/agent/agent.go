package agent

import (
	"context"
	"deepllm/components/mock"
	"deepllm/internal/data"
)

// Agent represents a base agent interface
type Agent interface {
	// Process processes the input and returns a response
	Process(ctx context.Context, input interface{}) (interface{}, error)
	// Name returns the agent's name
	Name() string
}

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	name      string
	model     mock.ChatModel
	tools     []mock.Tool
	DataQuery *data.DataQuery
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(name string, model mock.ChatModel, tools []mock.Tool, dataQuery *data.DataQuery) *BaseAgent {
	return &BaseAgent{
		name:      name,
		model:     model,
		tools:     tools,
		DataQuery: dataQuery,
	}
}

// Name returns the agent's name
func (b *BaseAgent) Name() string {
	return b.name
}

// BuildPrompt builds a prompt for the agent
func (b *BaseAgent) BuildPrompt(role string, context string) string {
	return `You are a ${role} for the Hangzhou Tourism Assistant system.
Your goal is to ${context}

Please consider:
1. User preferences and requirements
2. Budget constraints
3. Weather conditions
4. Location and distance
5. Ratings and reviews
6. Time constraints and scheduling

Respond in a clear and organized manner.`
}

// CreateReactAgent creates a ReAct agent with the given tools
func (b *BaseAgent) CreateReactAgent(ctx context.Context, systemPrompt string) (mock.Runnable[[]*mock.Message, *mock.Message], error) {
	// For now, we'll return a simple mock implementation
	return &mockReactAgent{
		model: b.model,
		tools: b.tools,
	}, nil
}

// mockReactAgent is a mock implementation of the ReAct agent
type mockReactAgent struct {
	model mock.ChatModel
	tools []mock.Tool
}

// Invoke processes the input and returns a response
func (m *mockReactAgent) Invoke(ctx context.Context, input []*mock.Message) (*mock.Message, error) {
	return m.model.Generate(ctx, input)
}
