package summaryresearch

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jacobtie/wikipedia-research-agent/internal/platform/prompt"
)

func (s *SummaryResearchTool) summarizeText(ctx context.Context, task, title, content string) (string, error) {
	messages := []*prompt.Message{
		{
			Role: prompt.USER_MESSAGE_ROLE,
			Content: fmt.Sprintf(`
Summarize the following CONTENT, paying attention to the researcher's given TASK so that relevant parts of the CONTENT relative to the TASK are summarized.
TASK: %s
CONTENT: %s
`, task, content),
		},
	}
	slog.Info("summarizing page", "page", title)
	for retries := 5; retries > 0; retries-- {
		resp, err := s.model.Invoke(ctx, &prompt.Prompt{
			SystemPrompt: `
Your role is to summarize an article for use in a piece of research. You are given a TASK which explains the original research topic and the CONTENT which is the body of an article in the research.
Do not attempt to solve the TASK. Instead, you should summarize the CONTENT. Use the TASK to understand what may be important for the researcher when deciding which parts of the CONTENT to summarize.
`,
			Messages: messages,
		})
		if err != nil {
			slog.Error("failed summarize content", "error", err)
			continue
		}
		return resp.Content, nil
	}
	return "", fmt.Errorf("failed to summarize content within retry limit")
}
