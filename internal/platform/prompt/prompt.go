package prompt

import "github.com/jacobtie/wikipedia-research-agent/internal/platform/tool"

type MessageRole uint8

const (
	USER_MESSAGE_ROLE = iota
	ASSISTANT_MESSAGE_ROLE
	TOOL_MESSAGE_ROLE
)

type Prompt struct {
	SystemPrompt string
	Messages     []*Message
	Tools        tool.Registry
}

type Message struct {
	Role      MessageRole
	Content   string
	ToolCalls []*tool.ToolCall
	ToolName  string
}

func New(systemPrompt, task string, toolRegistry tool.Registry) *Prompt {
	return &Prompt{
		SystemPrompt: systemPrompt,
		Messages: []*Message{
			{
				Role:    USER_MESSAGE_ROLE,
				Content: task,
			},
		},
		Tools: toolRegistry,
	}
}
