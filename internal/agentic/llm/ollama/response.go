package ollama

import (
	"fmt"

	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/llm"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/tool"
)

type response struct {
	Message struct {
		Content   string      `json:"content"`
		ToolCalls []*toolCall `json:"tool_calls"`
	} `json:"message"`
}

func (c *Client) formatResponse(res *response) (*llm.Response, error) {
	if len(res.Message.ToolCalls) == 0 && res.Message.Content == "" {
		return nil, fmt.Errorf("model returned neither tools nor content")
	}
	r := &llm.Response{Content: res.Message.Content}
	if len(res.Message.ToolCalls) > 0 {
		r.ToolCalls = make([]*tool.ToolCall, 0, len(res.Message.ToolCalls))
		for _, tc := range res.Message.ToolCalls {
			r.ToolCalls = append(r.ToolCalls, &tool.ToolCall{
				Name:   tc.Function.Name,
				KWArgs: tc.Function.Arguments,
			})
		}
	}
	return r, nil
}
