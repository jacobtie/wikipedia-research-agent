package wikipedia

import (
	"context"
	"fmt"
	"net/http"
)

type parsePageResult struct {
	Parse struct {
		WikiText string `json:"wikitext"`
	} `json:"parse"`
}

func (c *Client) ReadContent(ctx context.Context, pageID int) (string, error) {
	endpoint := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=parse&format=json&formatversion=2&prop=wikitext&pageid=%d",
		pageID,
	)
	resp, err := makeRequest[parsePageResult](ctx, c.httpClient, endpoint, http.MethodGet, nil)
	if err != nil {
		return "", fmt.Errorf("failed to read page content: %w", err)
	}
	return resp.Parse.WikiText, nil
}
