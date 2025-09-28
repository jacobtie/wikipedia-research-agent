package wikipedia

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type QueryPageResult struct {
	PageID  int
	Snippet string
}

type querySearchResult struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			PageID  int    `json:"pageid"`
			Snippet string `json:"snippet"`
		} `json:"search"`
	} `json:"query"`
}

func (c *Client) QueryPages(ctx context.Context, searchTerm string) (map[string]*QueryPageResult, error) {
	endpoint := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&formatversion=2&list=search&srlimit=5&srsearch=%s",
		url.QueryEscape(searchTerm),
	)
	pages, err := makeRequest[querySearchResult](ctx, c.httpClient, endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query pages: %w", err)
	}
	result := make(map[string]*QueryPageResult)
	for _, searchResult := range pages.Query.Search {
		slog.Info("found page", "page", searchResult.Title)
		result[searchResult.Title] = &QueryPageResult{
			PageID:  searchResult.PageID,
			Snippet: searchResult.Snippet,
		}
	}
	return result, nil
}
