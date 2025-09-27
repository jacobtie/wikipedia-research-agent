package llm

import (
	"context"

	"github.com/jacobtie/wikipedia-research-agent/internal/platform/prompt"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/tool"
)

type Response struct {
	Content   string
	ToolCalls []*tool.ToolCall
}

type LLM interface {
	Invoke(context.Context, *prompt.Prompt) (*Response, error)
}
