package gemini

import (
	"context"
	"fmt"

	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/llm"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/prompt"
	"google.golang.org/genai"
)

type Client struct {
	geminiClient *genai.Client
}

// Enforces interface implementation
var _ llm.LLM = (*Client)(nil)

func New(ctx context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	return &Client{
		geminiClient: client,
	}, nil
}

func (c *Client) Invoke(ctx context.Context, promptHistory *prompt.Prompt) (*llm.Response, error) {
	// c.geminiClient.Models.GenerateContent()
	return nil, nil
}
