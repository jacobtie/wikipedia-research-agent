package gemini

import (
	"fmt"
	"strings"

	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/llm"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/tool"
	"google.golang.org/genai"
)

func (c *Client) formatResponse(resp *genai.GenerateContentResponse) (*llm.Response, error) {
	if len(resp.Candidates) != 1 {
		return nil, fmt.Errorf("candidates length was %d when it should be 1", len(resp.Candidates))
	}
	candidate := resp.Candidates[0]
	r := &llm.Response{}
	textParts := make([]string, 0)
	r.ToolCalls = make([]*tool.ToolCall, 0)
	for _, part := range candidate.Content.Parts {
		if part.FunctionCall != nil {
			r.ToolCalls = append(r.ToolCalls, &tool.ToolCall{
				Name:   part.FunctionCall.Name,
				KWArgs: part.FunctionCall.Args,
			})
			continue
		}
		textParts = append(textParts, part.Text)
	}
	r.Content = strings.Join(textParts, "\n")
	if r.Content == "" && len(r.ToolCalls) == 0 {
		return nil, fmt.Errorf(
			"model unexpectedly returned no content and no tool calls with finish reason '%s' and finish message '%s'",
			candidate.FinishReason,
			candidate.FinishMessage,
		)
	}
	return r, nil
}
