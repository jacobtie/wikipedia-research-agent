package orchestrator

import (
	"context"

	"github.com/jacobtie/wikipedia-research-agent/internal/config"
	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator/output"
	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator/summaryresearch"
	"github.com/jacobtie/wikipedia-research-agent/internal/platform/agent"
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

const systemPrompt = `
You are a research agent. Your job is to take in a TASK, perform research using the tools available, and then produce a clear, well-structured written output.

Rules for Tool Usage

Research Phase (mandatory)

You must begin every TASK by calling the summary_research_tool.

Call it with an initial search_term derived from the TASK.

Review the summaries you receive.

Continue calling the tool with new or refined search terms, based on what you learned from prior results.

You must call this tool at least once before doing anything else.

Output Phase (only after research is complete)

When, and only when, the research phase is fully complete, prepare the final written output.

Then, and only then, call the output_writer_tool.

Important: You may not call output_writer_tool until at least one summary_research_tool call has been made.

Never call both tools in the same step.

Guidelines

Always separate research from writing.

Never skip the research phase, even if you think you already know the answer.

Summarize and synthesize in your own words.

Ensure the final output is accurate, coherent, and directly answers the TASK.

Your workflow must always follow this order:
Research → Research → (repeat as needed) → Final Output → Output Tool.
`

func (o *Orchestrator) Run(ctx context.Context) <-chan *agent.AgentResult {
	model := ollama.New(o.cfg)
	summaryResearchTool := summaryresearch.New(o.task, model)
	outputWriterTool := output.New()
	registry := make(tool.Registry)
	registry[summaryResearchTool.GetName()] = summaryResearchTool
	registry[outputWriterTool.GetName()] = outputWriterTool
	mainAgent := agent.New(o.cfg, model, systemPrompt, registry)
	return mainAgent.Run(ctx, o.task)
}
