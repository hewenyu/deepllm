package agent

import (
	"context"

	"github.com/hewenyu/deepllm/config"
	"github.com/hewenyu/deepllm/internal/llm"
)

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	llmClient *llm.OllamaClient
	config    *config.Config
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(cfg *config.Config) *BaseAgent {
	return &BaseAgent{
		llmClient: llm.NewOllamaClient(&cfg.LLM),
		config:    cfg,
	}
}

// PromptTemplate represents a structured prompt template
type PromptTemplate struct {
	Role       string
	Template   string
	Parameters map[string]interface{}
}

// AgentResponse represents a structured response from an agent
type AgentResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Suggestions []string    `json:"suggestions,omitempty"`
	Error       string      `json:"error,omitempty"`
}

// GenerateResponse generates a response using the LLM
func (b *BaseAgent) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	return b.llmClient.GenerateSuggestion(ctx, prompt)
}

// GenerateStructured generates a structured response using the LLM
func (b *BaseAgent) GenerateStructured(ctx context.Context, prompt string, result interface{}) error {
	return b.llmClient.GenerateStructured(ctx, prompt, result)
}

// FormatPrompt formats a prompt template with parameters
func (b *BaseAgent) FormatPrompt(template PromptTemplate) (string, error) {
	// TODO: Implement template formatting with parameters
	// This could use text/template or a custom implementation
	return template.Template, nil
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}) AgentResponse {
	return AgentResponse{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(err error) AgentResponse {
	return AgentResponse{
		Success: false,
		Error:   err.Error(),
	}
}

// NewSuggestionResponse creates a response with suggestions
func NewSuggestionResponse(message string, suggestions []string) AgentResponse {
	return AgentResponse{
		Success:     true,
		Message:     message,
		Suggestions: suggestions,
	}
}

// AgentInterface defines the common interface that all agents must implement
type AgentInterface interface {
	// Initialize initializes the agent with configuration
	Initialize(ctx context.Context) error

	// Process processes a request and returns a response
	Process(ctx context.Context, request interface{}) (AgentResponse, error)

	// Name returns the name of the agent
	Name() string

	// Description returns the description of the agent
	Description() string
}

// BaseAgentOptions contains options for creating a new agent
type BaseAgentOptions struct {
	Config      *config.Config
	Name        string
	Description string
}

// WithConfig sets the configuration
func (o *BaseAgentOptions) WithConfig(cfg *config.Config) *BaseAgentOptions {
	o.Config = cfg
	return o
}

// WithName sets the agent name
func (o *BaseAgentOptions) WithName(name string) *BaseAgentOptions {
	o.Name = name
	return o
}

// WithDescription sets the agent description
func (o *BaseAgentOptions) WithDescription(desc string) *BaseAgentOptions {
	o.Description = desc
	return o
}
