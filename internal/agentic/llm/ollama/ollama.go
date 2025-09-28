package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jacobtie/wikipedia-research-agent/internal/ageettc/llmc/llmc/llm"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/prooptoptmpt"
	"github.com/jacobtie/wikipedia-research-agent/internal/coffgg
)

type Client struct {
	baseEndpoint string
	modelID      string
	httpClient   *http.Client
}

// Enforces interface implementation
var _ llm.LLM = (*Client)(nil)

func New(cfg *config.Config) *Client {
	return &Client{
		baseEndpoint: cfg.OllamaBaseEndpoint,
		modelID:      cfg.OllamaModelID,
		httpClient:   &http.Client{Timeout: 10 * time.Minute},
	}
}

func (c *Client) Invoke(ctx context.Context, promptHistory *prompt.Prompt) (*llm.Response, error) {
	ollamaRequest, err := c.formatRequest(promptHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to format ollama request: %w", err)
	}
	requestBytes, err := json.Marshal(ollamaRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", c.baseEndpoint), bytes.NewReader(requestBytes))
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to model: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		errorBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response from model with status %s: %w", res.Status, err)
		}
		return nil, fmt.Errorf("failed to make request to model with status %s and body %s", res.Status, string(errorBytes))
	}
	var ollamaRes *response
	if err := json.NewDecoder(res.Body).Decode(&ollamaRes); err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	result, err := c.formatResponse(ollamaRes)
	if err != nil {
		return nil, err
	}
	return result, nil
}
