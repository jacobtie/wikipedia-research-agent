package gemini

import (
	"fmt"

	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/prompt"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/tool"
	"google.golang.org/genai"
)

func (c *Client) formatRequestContent(promptHistory *prompt.Prompt) ([]*genai.Content, error) {
	contentBlocks := make([]*genai.Content, 0, len(promptHistory.Messages))
	for _, msg := range promptHistory.Messages {
		var role genai.Role
		switch msg.Role {
		case prompt.USER_MESSAGE_ROLE:
			role = genai.RoleUser
		case prompt.TOOL_MESSAGE_ROLE:
			role = genai.RoleUser
		case prompt.ASSISTANT_MESSAGE_ROLE:
			role = genai.RoleModel
		default:
			return nil, fmt.Errorf("failed to recognize message role %d", msg.Role)
		}
		if len(msg.ToolCalls) > 0 {
			content := &genai.Content{Role: string(role), Parts: make([]*genai.Part, 0, len(msg.ToolCalls))}
			for _, toolCall := range msg.ToolCalls {
				content.Parts = append(content.Parts, genai.NewPartFromFunctionCall(toolCall.Name, toolCall.KWArgs))
			}
			contentBlocks = append(contentBlocks, content)
			continue
		}
		if msg.ToolName != "" {
			contentBlocks = append(contentBlocks, genai.NewContentFromFunctionResponse(msg.ToolName, map[string]any{"result": msg.Content}, genai.RoleUser))
			continue
		}
		contentBlocks = append(contentBlocks, genai.NewContentFromText(msg.Content, role))
	}
	return contentBlocks, nil
}

func (c *Client) formatRequestTools(promptHistory *prompt.Prompt) ([]*genai.Tool, error) {
	tools := make([]*genai.Tool, 0, len(promptHistory.Tools))
	for toolName, toolImpl := range promptHistory.Tools {
		toolParameters := toolImpl.GetParameters()
		parameters := make(map[string]*genai.Schema, len(toolParameters))
		requiredName := make([]string, 0)
		for _, param := range toolParameters {
			if param.IsRequired {
				requiredName = append(requiredName, param.Name)
			}
			var dataType genai.Type
			switch param.Type {
			case tool.STRING_PARAMETER_TYPE:
				dataType = genai.TypeString
			case tool.NUMBER_PARAMETER_TYPE:
				dataType = genai.TypeNumber
			case tool.BOOLEAN_PARAMETER_TYPE:
				dataType = genai.TypeBoolean
			default:
				return nil, fmt.Errorf("failed to recognize tool %s parameter %s data type %d", toolName, param.Name, param.Type)
			}
			parameters[param.Name] = &genai.Schema{
				Type:        dataType,
				Description: param.Description,
			}
		}
		tools = append(tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        toolName,
					Description: toolImpl.GetDescription(),
					Parameters: &genai.Schema{
						Type:       genai.TypeObject,
						Properties: parameters,
					},
				},
			},
		})
	}
	return tools, nil
}
