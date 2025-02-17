package coordinator

import (
	"context"
	"deepllm/components/agent"
	"deepllm/components/agent/planner"
	"deepllm/components/mock"
	"deepllm/internal/data"
	"encoding/json"
	"fmt"
)

// CoordinatorAgent manages the multi-agent system
type CoordinatorAgent struct {
	*agent.BaseAgent
	plannerAgent *planner.PlannerAgent
}

// NewCoordinatorAgent creates a new coordinator agent
func NewCoordinatorAgent(model mock.ChatModel, tools []mock.Tool, dataQuery *data.DataQuery) *CoordinatorAgent {
	return &CoordinatorAgent{
		BaseAgent:    agent.NewBaseAgent("coordinator", model, tools, dataQuery),
		plannerAgent: planner.NewPlannerAgent(model, tools, dataQuery),
	}
}

// Process processes the trip planning request
func (c *CoordinatorAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
	request, ok := input.(*data.TripPlanRequest)
	if !ok {
		return nil, fmt.Errorf("invalid input type for coordinator agent")
	}

	// Validate request
	if err := c.validateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}

	// Get trip plan from planner agent
	planResult, err := c.plannerAgent.Process(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create trip plan: %v", err)
	}

	// Use LLM to review and refine the plan
	systemPrompt := c.BuildPrompt(
		"Trip Plan Reviewer",
		"review and refine the trip plan to ensure it meets all requirements and constraints",
	)

	agent, err := c.CreateReactAgent(ctx, systemPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %v", err)
	}

	// Prepare input for LLM
	messages := []*mock.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role: "user",
			Content: fmt.Sprintf("Please review and refine the following trip plan:\n\n"+
				"Original request:\n%+v\n\n"+
				"Generated plan:\n%s",
				request,
				planResult.(*mock.Message).Content,
			),
		},
	}

	// Get the refined plan from LLM
	result, err := agent.Invoke(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to refine trip plan: %v", err)
	}

	return result, nil
}

// validateRequest validates the trip planning request
func (c *CoordinatorAgent) validateRequest(request *data.TripPlanRequest) error {
	// Check dates
	if request.StartDate.IsZero() || request.EndDate.IsZero() {
		return fmt.Errorf("invalid dates")
	}
	if request.EndDate.Before(request.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	// Check location
	if request.Location.Name == "" {
		return fmt.Errorf("location name is required")
	}
	if request.Location.Latitude == 0 && request.Location.Longitude == 0 {
		return fmt.Errorf("invalid location coordinates")
	}

	// Check budget
	if request.Budget.Total <= 0 {
		return fmt.Errorf("total budget must be positive")
	}
	if request.Budget.Hotel <= 0 {
		return fmt.Errorf("hotel budget must be positive")
	}
	if request.Budget.Food <= 0 {
		return fmt.Errorf("food budget must be positive")
	}
	if request.Budget.Activity <= 0 {
		return fmt.Errorf("activity budget must be positive")
	}

	// Check party size
	if request.PartySize <= 0 {
		return fmt.Errorf("party size must be positive")
	}

	return nil
}

// formatTripPlan formats the trip plan as JSON
func (c *CoordinatorAgent) formatTripPlan(plan *data.TripPlan) (string, error) {
	bytes, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal trip plan: %v", err)
	}
	return string(bytes), nil
}
