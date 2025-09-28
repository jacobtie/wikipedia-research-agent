package summaryresearch

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator/wikipedia"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/llm"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/tool"
)

type SummaryResearchTool struct {
	task       string
	wikiClient *wikipedia.Client
	model      llm.LLM
}

var _ tool.Tool = (*SummaryResearchTool)(nil)

func New(task string, model llm.LLM) *SummaryResearchTool {
	return &SummaryResearchTool{
		task:       task,
		wikiClient: wikipedia.New(),
		model:      model,
	}
}

func (s *SummaryResearchTool) GetName() string {
	return "summary_research_tool"
}

func (s *SummaryResearchTool) GetDescription() string {
	return "The summary_research_tool uses a search_term to search Wikipedia for relevantly related pages and returns you summaries of those pages"
}

func (s *SummaryResearchTool) GetParameters() []*tool.Parameter {
	return []*tool.Parameter{
		{
			Name:        "search_term",
			Description: "The search_term parameter is the word or short phrase to search for in Wikipedia to pull summaries of relevantly related pages to the search_term",
			Type:        tool.STRING_PARAMETER_TYPE,
			IsRequired:  true,
		},
	}
}

func (s *SummaryResearchTool) Run(ctx context.Context, kwargs map[string]any) (string, error) {
	rawSearchTerm, ok := kwargs["search_term"]
	if !ok {
		return "", fmt.Errorf("failed to run summary_research_tool: missing required parameter 'search_term'")
	}
	searchTerm, ok := rawSearchTerm.(string)
	if !ok {
		return "", fmt.Errorf("failed to run summary_research_tool: search_term parameter must be a string")
	}
	pageResults, err := s.wikiClient.QueryPages(ctx, searchTerm)
	if err != nil {
		return "", fmt.Errorf("failed to run summary_research_tool: %w", err)
	}
	relevantPageResults := s.getRelevantPages(ctx, pageResults)
	if len(relevantPageResults) == 0 {
		return "", fmt.Errorf("no pages were relevant, try using a more general search_term")
	}
	results := make(map[string]string)
	var resultsMu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(relevantPageResults))
	for pageTitle, pageResult := range relevantPageResults {
		go func(title string, pageID int) {
			defer wg.Done()
			pageContent, err := s.wikiClient.ReadContent(ctx, pageID)
			if err != nil {
				slog.Error("failed to read page content", "error", err.Error())
				return
			}
			summary, err := s.summarizeText(ctx, s.task, pageContent)
			if err != nil {
				slog.Error("failed to summarize page content", "error", err.Error())
				return
			}
			resultsMu.Lock()
			defer resultsMu.Unlock()
			results[title] = summary
		}(pageTitle, pageResult.PageID)
	}
	wg.Wait()
	var summariesBuilder strings.Builder
	summariesBuilder.WriteString("Page Summaries\n")
	for pageTitle, summary := range results {
		summariesBuilder.WriteString(fmt.Sprintf(`
Title: %s
Summary: %s
`, pageTitle, summary))
	}
	return summariesBuilder.String(), nil
}
