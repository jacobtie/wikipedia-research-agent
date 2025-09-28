package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/jacobtie/wikipedia-research-agent/internal/config"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/llm"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/prompt"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/tool"
)

type Agent struct {
	model         llm.LLM
	maxIterations int
	systemPrompt  string
	toolRegistry  tool.Registry
}

func New(cfg *config.Config, model llm.LLM, systemPrompt string, toolRegistry tool.Registry) *Agent {
	return &Agent{
		model:         model,
		maxIterations: cfg.MaxIterations,
		systemPrompt:  systemPrompt,
		toolRegistry:  toolRegistry,
	}
}

type AgentResult struct {
	Msg   string
	Error error
}

func (a *Agent) Run(ctx context.Context, task string) <-chan *AgentResult {
	results := make(chan *AgentResult)
	go a.runIterations(ctx, results, task)
	return results
}

func (a *Agent) runIterations(ctx context.Context, results chan<- *AgentResult, task string) {
	defer close(results)
	promptHistory := prompt.New(a.systemPrompt, task, a.toolRegistry)
	results <- &AgentResult{Msg: fmt.Sprintf("Task: %s", task)}
	foundAnswer := false
	for i := 1; !foundAnswer && i <= a.maxIterations; i++ {
		if a.runIteration(ctx, i, results, promptHistory) {
			foundAnswer = true
		}
	}
	if !foundAnswer {
		results <- &AgentResult{Error: fmt.Errorf("Finished iteration without an answer")}
	}
}

func (a *Agent) runIteration(ctx context.Context, i int, results chan<- *AgentResult, promptHistory *prompt.Prompt) bool {
	results <- &AgentResult{Msg: fmt.Sprintf("Iteration #%d: starting", i)}
	modelRes, err := a.model.Invoke(ctx, promptHistory)
	if err != nil {
		results <- &AgentResult{Error: fmt.Errorf("Iteration #%d: failed to invoke model: %w", i, err)}
		return false
	}
	if len(modelRes.ToolCalls) == 0 {
		results <- &AgentResult{Msg: fmt.Sprintf("Iteration #%d: %s", i, modelRes.Content)}
		return true
	}
	promptHistory.Messages = append(promptHistory.Messages, &prompt.Message{
		Role:      prompt.ASSISTANT_MESSAGE_ROLE,
		Content:   modelRes.Content,
		ToolCalls: modelRes.ToolCalls,
	})
	var promptHistoryMu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(modelRes.ToolCalls))
	for _, tc := range modelRes.ToolCalls {
		go func(toolCall *tool.ToolCall) {
			defer wg.Done()
			results <- &AgentResult{Msg: fmt.Sprintf("Iteration #%d: calling tool %s and args %#v", i, toolCall.Name, toolCall.KWArgs)}
			toolRes, err := a.callTool(ctx, toolCall)
			if err != nil {
				results <- &AgentResult{Msg: fmt.Sprintf("Iteration #%d: failed to call tools: %s", i, err.Error())}
				promptHistoryMu.Lock()
				defer promptHistoryMu.Unlock()
				promptHistory.Messages = append(promptHistory.Messages, &prompt.Message{
					Role:    prompt.USER_MESSAGE_ROLE,
					Content: fmt.Sprintf("you failed to call the tool %s because of the following reason: %s", toolCall.Name, err.Error()),
				})
				return
			}
			results <- &AgentResult{Msg: fmt.Sprintf("Iteration #%d: got tool %s result: %s", i, toolCall.Name, toolRes)}
			promptHistoryMu.Lock()
			defer promptHistoryMu.Unlock()
			promptHistory.Messages = append(promptHistory.Messages, &prompt.Message{
				Role:     prompt.TOOL_MESSAGE_ROLE,
				Content:  toolRes,
				ToolName: toolCall.Name,
			})
		}(tc)
	}
	wg.Wait()
	return false
}

func (a *Agent) callTool(ctx context.Context, toolCall *tool.ToolCall) (string, error) {
	tool, ok := a.toolRegistry[toolCall.Name]
	if !ok {
		return "", fmt.Errorf("could not find tool with name %s", toolCall.Name)
	}
	res, err := tool.Run(ctx, toolCall.KWArgs)
	if err != nil {
		return "", fmt.Errorf("failed to run tool: %s: %w", toolCall.Name, err)
	}
	return res, nil
}
