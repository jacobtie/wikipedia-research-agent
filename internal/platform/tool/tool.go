package tool

import "context"

type ParameterType byte

const (
	STRING_PARAMETER_TYPE ParameterType = iota
	NUMBER_PARAMETER_TYPE
	BOOLEAN_PARAMETER_TYPE
)

type Registry map[string]Tool

type Tool interface {
	GetName() string
	GetDescription() string
	GetParameters() []*Parameter
	Run(ctx context.Context, kwargs map[string]any) (string, error)
}

type Parameter struct {
	Name        string
	Type        ParameterType
	Description string
	IsRequired  bool
}

type ToolCall struct {
	Name   string
	KWArgs map[string]any
}
