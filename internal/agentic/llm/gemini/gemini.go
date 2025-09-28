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
	request, err := c.formatRequestContent(promptHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to format gemini request: %w", err)
	}
	tools, err := c.formatRequestTools(promptHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to format gemini tools: %w", err)
	}
	resp, err := c.geminiClient.Models.GenerateContent(ctx, "gemini-2.0-flash", request, &genai.GenerateContentConfig{
		SystemInstruction: genai.Text(promptHistory.SystemPrompt)[0],
		Tools:             tools,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke gemini model: %w", err)
	}
	formattedResp, err := c.formatResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to format response from gemini model: %w", err)
	}
	return formattedResp, nil
}
