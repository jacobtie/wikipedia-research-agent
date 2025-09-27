package orchestrator

import (
	"context"

	"github.com/jacobtie/wikipedia-research-agent/internal/platform/agent"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/config"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/llm/ollama"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/tool"
)

type Orchestrator struct {
	cfg  *config.Config
	task string
}

func New(cfg *config.Config, task string) *Orchestrator {
	return &Orchestrator{
		cfg:  cfg,
		task: task,
	}
}

const MAIN_AGENT_SYSTEM_PROMPT = `
	TODO
`

func (o *Orchestrator) Run(ctx context.Context) <-chan *agent.AgentResult {
	model := ollama.New(o.cfg)
	mainAgent := agent.New(o.cfg, model, MAIN_AGENT_SYSTEM_PROMPT, make(tool.Registry))
	return mainAgent.Run(ctx, o.task)
}
