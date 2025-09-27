package ollama

import (
	"fmt"
	"strings"

	"github.com/jacobtie/wikipedia-research-agent/internal/platform/prompt"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/tool"
)

type request struct {
	Model    string            `json:"model"`
	Messages []*requestMessage `json:"messages"`
	Stream   bool              `json:"stream"`
	Tools    []*requestTool    `json:"tools"`
}

type requestMessage struct {
	Role      string      `json:"role"`
	Content   string      `json:"content"`
	ToolCalls []*toolCall `json:"tool_calls,omitempty"`
	ToolName  string      `json:"tool_name,omitempty"`
}

type toolCall struct {
	Function *toolCallFunction `json:"function"`
}

type toolCallFunction struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type requestTool struct {
	Type     string        `json:"type"`
	Function *toolFunction `json:"function"`
}

type toolFunction struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Parameters  *toolFunctionParameters `json:"parameters"`
}

type toolFunctionParameters struct {
	Type       string                                     `json:"type"`
	Properties map[string]*toolFunctionParametersProperty `json:"properties"`
	Required   []string                                   `json:"required"`
}

type toolFunctionParametersProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (c *Client) formatRequest(promptHistory *prompt.Prompt) (*request, error) {
	r := &request{
		Model:    c.modelID,
		Messages: make([]*requestMessage, 0, len(promptHistory.Messages)+1),
		Stream:   false,
		Tools:    make([]*requestTool, 0, len(promptHistory.Tools)),
	}
	r.Messages = append(r.Messages, &requestMessage{
		Role:    "system",
		Content: strings.TrimSpace(promptHistory.SystemPrompt),
	})
	requestMessages, err := formatMessages(promptHistory.Messages)
	if err != nil {
		return nil, err
	}
	r.Messages = append(r.Messages, requestMessages...)
	requestTools, err := formatTools(promptHistory.Tools)
	if err != nil {
		return nil, err
	}
	r.Tools = append(r.Tools, requestTools...)
	return r, nil
}

func formatMessages(messages []*prompt.Message) ([]*requestMessage, error) {
	requestMessages := make([]*requestMessage, 0, len(messages))
	for _, msg := range messages {
		var role string
		switch msg.Role {
		case prompt.USER_MESSAGE_ROLE:
			role = "user"
		case prompt.TOOL_MESSAGE_ROLE:
			role = "tool"
		case prompt.ASSISTANT_MESSAGE_ROLE:
			role = "assistant"
		default:
			return nil, fmt.Errorf("failed to recognize message role %d", msg.Role)
		}
		var toolCalls []*toolCall
		if len(msg.ToolCalls) > 0 {
			toolCalls = make([]*toolCall, 0, len(msg.ToolCalls))
			for _, tc := range msg.ToolCalls {
				toolCalls = append(toolCalls, &toolCall{
					Function: &toolCallFunction{
						Name:      tc.Name,
						Arguments: tc.KWArgs,
					},
				})
			}
		}
		requestMessages = append(requestMessages, &requestMessage{
			Role:      role,
			Content:   msg.Content,
			ToolCalls: toolCalls,
			ToolName:  msg.ToolName,
		})
	}
	return requestMessages, nil
}

func formatTools(tools tool.Registry) ([]*requestTool, error) {
	requestTools := make([]*requestTool, 0, len(tools))
	for tName, t := range tools {
		toolProperties := make(map[string]*toolFunctionParametersProperty)
		requiredProperties := make([]string, 0)
		for _, tProp := range t.GetParameters() {
			if tProp.IsRequired {
				requiredProperties = append(requiredProperties, tProp.Name)
			}
			var dataType string
			switch tProp.Type {
			case tool.STRING_PARAMETER_TYPE:
				dataType = "string"
			case tool.NUMBER_PARAMETER_TYPE:
				dataType = "number"
			case tool.BOOLEAN_PARAMETER_TYPE:
				dataType = "boolean"
			default:
				return nil, fmt.Errorf("failed to recognize tool %s parameter %s data type %d", tName, tProp.Name, tProp.Type)
			}
			toolProperties[tProp.Name] = &toolFunctionParametersProperty{
				Type:        dataType,
				Description: tProp.Description,
			}
		}
		requestTools = append(requestTools, &requestTool{
			Type: "function",
			Function: &toolFunction{
				Name:        tName,
				Description: t.GetDescription(),
				Parameters: &toolFunctionParameters{
					Type:       "object",
					Properties: toolProperties,
					Required:   requiredProperties,
				},
			},
		})
	}
	return requestTools, nil
}
