package wikipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func makeRequest[T any](ctx context.Context, httpClient *http.Client, endpoint, method string, body io.Reader) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "TestBot")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	if resp.StatusCode != 200 {
		errResp, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response after status code %d: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("received error code %d with body %s", resp.StatusCode, errResp)
	}
	var val T
	if err := json.NewDecoder(resp.Body).Decode(&val); err != nil {
		return nil, fmt.Errorf("failed to read success response from: %w", err)
	}
	return &val, nil
}
