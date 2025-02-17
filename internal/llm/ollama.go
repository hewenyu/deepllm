package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hewenyu/deepllm/config"
	"github.com/pkg/errors"
)

// OllamaClient represents an Ollama API client
type OllamaClient struct {
	config *config.LLMConfig
	client *http.Client
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(cfg *config.LLMConfig) *OllamaClient {
	return &OllamaClient{
		config: cfg,
		client: &http.Client{},
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatResponse struct {
	Model         string `json:"model"`
	Response      string `json:"response"`
	Done          bool   `json:"done"`
	Context       []int  `json:"context,omitempty"`
	TotalDuration int64  `json:"total_duration,omitempty"`
}

// Chat sends a chat request to Ollama
func (c *OllamaClient) Chat(ctx context.Context, messages []ChatMessage) (*ChatResponse, error) {
	url := fmt.Sprintf("%s/api/chat", c.config.BaseURL)

	reqBody := ChatRequest{
		Model:    c.config.Model,
		Messages: messages,
		Stream:   false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &chatResp, nil
}

// GenerateSuggestion generates a suggestion using chat messages
func (c *OllamaClient) GenerateSuggestion(ctx context.Context, prompt string) (string, error) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	resp, err := c.Chat(ctx, messages)
	if err != nil {
		return "", err
	}

	return resp.Response, nil
}

// GenerateStructured generates a structured response using chat messages
func (c *OllamaClient) GenerateStructured(ctx context.Context, prompt string, result interface{}) error {
	messages := []ChatMessage{
		{
			Role:    "system",
			Content: "Please provide response in valid JSON format.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	resp, err := c.Chat(ctx, messages)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(resp.Response), result)
}
