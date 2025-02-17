package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// WeatherParams represents the parameters for the weather tool
type WeatherParams struct {
	City string `json:"city" jsonschema:"description=The city to get weather for"`
}

// WeatherTool is a simple mock weather tool
type WeatherTool struct{}

// NewWeatherTool creates a new weather tool
func NewWeatherTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"get_weather",
		"Get the current weather for a city",
		func(_ context.Context, params *WeatherParams) (string, error) {
			// This is a mock implementation
			response := map[string]interface{}{
				"city":        params.City,
				"temperature": "22°C",
				"condition":   "晴朗",
				"humidity":    "65%",
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return "", fmt.Errorf("failed to marshal response: %v", err)
			}

			return string(jsonResponse), nil
		},
	)
}

// Info returns the tool info
func (w *WeatherTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_weather",
		Desc: "Get the current weather for a city",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"city": {
				Desc:     "The city to get weather for",
				Type:     schema.String,
				Required: true,
			},
		}),
	}, nil
}

// InvokableRun runs the weather tool
func (w *WeatherTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params WeatherParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", fmt.Errorf("failed to unmarshal arguments: %v", err)
	}

	// This is a mock implementation
	response := map[string]interface{}{
		"city":        params.City,
		"temperature": "22°C",
		"condition":   "晴朗",
		"humidity":    "65%",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %v", err)
	}

	return string(jsonResponse), nil
}
