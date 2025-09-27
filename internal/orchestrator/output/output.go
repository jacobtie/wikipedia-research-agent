package output

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jacobtie/wikipedia-research-agent/internal/platform/tool"
)

type WriterTool struct{}

// Force interface implementation
var _ tool.Tool = (*WriterTool)(nil)

func New() *WriterTool {
	return &WriterTool{}
}

func (w *WriterTool) GetName() string {
	return "output_writer"
}

func (w *WriterTool) GetDescription() string {
	return "The output_writer tool writes the final output to a file once all research is complete"
}

func (w *WriterTool) GetParameters() []*tool.Parameter {
	return []*tool.Parameter{
		{
			Name:        "content",
			Type:        tool.STRING_PARAMETER_TYPE,
			Description: "The content parameter is the content of the final research to write to the file",
			IsRequired:  true,
		},
	}
}

func (w *WriterTool) Run(ctx context.Context, kwargs map[string]any) (string, error) {
	rawContent, ok := kwargs["content"]
	if !ok {
		return "", fmt.Errorf("missing required string parameter 'content'")
	}
	content, ok := rawContent.(string)
	if !ok {
		return "", fmt.Errorf("content parameter must be a string")
	}
	unixTimestamp := time.Now().Unix()
	filePath := path.Join("output", fmt.Sprintf("output_%d", unixTimestamp))
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file to disk: %w", err)
	}
	return fmt.Sprintf("successfully wrote file to path %s", filePath), nil
}
