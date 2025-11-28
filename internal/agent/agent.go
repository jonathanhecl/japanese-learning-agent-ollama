package agent

import (
	"context"
	"fmt"

	"github.com/jonathanhecl/japanese-learning-agent-ollama/internal/profile"
	"github.com/jonathanhecl/japanese-learning-agent-ollama/internal/prompts"

	"github.com/jonathanhecl/gollama"
)

// Agent handles the interaction with the LLM.
type Agent struct {
	llm *gollama.Gollama
}

// New creates a new Agent instance.
func New(model string) (*Agent, error) {
	llm := gollama.New(model)
	if llm == nil {
		return nil, fmt.Errorf("failed to create gollama instance")
	}
	return &Agent{llm: llm}, nil
}

// SetSystemPrompt sets the system prompt for the agent.
func (a *Agent) SetSystemPrompt(prompt string) {
	a.llm.SetSystemPrompt(prompt)
}

// ExtractProfile asks the LLM to extract the user profile from the input.
func (a *Agent) ExtractProfile(ctx context.Context, input string) (*profile.UserProfile, error) {
	prompt := prompts.GetProfileExtractionPrompt(input)
	resp, err := a.llm.Chat(ctx, prompt, gollama.StructToStructuredFormat(profile.UserProfile{}))
	if err != nil {
		return nil, err
	}

	var user profile.UserProfile
	if err := resp.DecodeContent(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Chat sends a message to the LLM and returns the response.
func (a *Agent) Chat(ctx context.Context, input string) (string, error) {
	resp, err := a.llm.Chat(ctx, input)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}
