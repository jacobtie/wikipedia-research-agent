package summaryresearch

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/prompt"
	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator/wikipedia"
)

func (s *SummaryResearchTool) getRelevantPages(ctx context.Context, pageResults map[string]*wikipedia.QueryPageResult) map[string]*wikipedia.QueryPageResult {
	relevantPageResults := make(map[string]*wikipedia.QueryPageResult)
	var resultsMu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(pageResults))
	for pageTitle, pageResult := range pageResults {
		go func(title string, result *wikipedia.QueryPageResult) {
			defer wg.Done()
			isRelevant, err := s.isPageRelevant(ctx, title, result.Snippet)
			if err != nil {
				slog.Error("failed to determine if page is relevant", "error", err.Error())
				return
			}
			slog.Info("determined page relevancy", "page", title, "isRelevant", isRelevant)
			if !isRelevant {
				return
			}
			resultsMu.Lock()
			defer resultsMu.Unlock()
			relevantPageResults[title] = result
		}(pageTitle, pageResult)
	}
	wg.Wait()
	return relevantPageResults
}

func (s *SummaryResearchTool) isPageRelevant(ctx context.Context, pageTitle, snippet string) (bool, error) {
	messages := []*prompt.Message{
		{
			Role: prompt.USER_MESSAGE_ROLE,
			Content: fmt.Sprintf(`
Evaluate whether the document, described by the TITLE and SNIPPET is related to the research TASK by answering one word of either "YES" or "NO".
TASK: %s
TITLE: %s
SNIPPET: %s
`, s.task, pageTitle, snippet),
		},
	}
	for retries := 5; retries > 0; retries-- {
		resp, err := s.model.Invoke(ctx, &prompt.Prompt{
			SystemPrompt: `
You are an evaluator whose task is to decide whether a document is relevant to a research topic.

You will be given three inputs:

TASK – the overarching research topic.

TITLE – the document title.

SNIPPET – a short excerpt from the document.

Instructions:

First, read the TASK carefully to understand the research goal.

Then, read the TITLE and SNIPPET.

If the document contains information that directly relates to, supports, or provides background context for the TASK, it is relevant.

If the document is unrelated or off-topic, it is not relevant.

Output Rules:

Respond with exactly one word:

"YES" if the document is relevant.

"NO" if the document is not relevant.

Do not explain your reasoning.

Do not output anything except "YES" or "NO".

Examples:

Example 1
TASK: Write a research paper about Vlad the Impaler and how his policies impacted his kingdom
TITLE: Vlad II Dracul
SNIPPET: Internationally known as the father of Vlad the Impaler, or Dracula. Born an illegitimate son of Mircea I of Wallachia…
Answer: YES

Example 2
TASK: Write a research paper about Vlad the Impaler and how his policies impacted his kingdom
TITLE: Modern Romanian Cuisine
SNIPPET: Traditional dishes include mămăligă, sarmale, and mici…
Answer: NO

Example 3
TASK: Research the life of Albert Einstein with a focus on his contributions to physics
TITLE: Theory of Relativity
SNIPPET: The special and general theories of relativity were developed by Albert Einstein in the early 20th century…
Answer: YES

Example 4
TASK: Research the life of Albert Einstein with a focus on his contributions to physics
TITLE: Marie Curie
SNIPPET: She pioneered research on radioactivity and was the first woman to win a Nobel Prize…
Answer: NO
`,
			Messages: messages,
		})
		if err != nil {
			slog.Error("failed to determine whether page was relevant", "error", err.Error())
			continue
		}
		content := strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(resp.Content), "\""), "\"")
		if content != "yes" && content != "no" {
			messages = append(messages, &prompt.Message{
				Role:    prompt.ASSISTANT_MESSAGE_ROLE,
				Content: resp.Content,
			}, &prompt.Message{
				Role: prompt.USER_MESSAGE_ROLE,
				Content: fmt.Sprintf(
					`You are instructed to answer with one word, either "YES" or "NO" but you answered with "%s" which is not one of those two words exactly. Try again and answer only with one word either "YES" or "NO".`,
					resp.Content,
				),
			})
		}
		return content == "yes", nil
	}
	return false, fmt.Errorf("failed to get relevant page within retry limit")
}
