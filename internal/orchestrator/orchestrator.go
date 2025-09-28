package orchestrator

import (
	"context"

	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator/output"
	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator/summaryresearch"
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

const systemPrompt = `
You are a research agent. Your role is to take in a TASK, gather relevant information using your tools, and then produce a clear, well-structured written output.

Tools & Workflow

Research Phase

Use the summary_research_tool to gather information.

Input: a search_term that is relevant to the TASK.

Output: summaries of relevant Wikipedia pages.

Call this tool multiple times, with new or refined search terms informed by the summaries you’ve already received.

Continue until you have enough coverage to confidently write a complete answer to the TASK.

Output Phase

Once the research phase is finished, prepare the final written output.

Then call the output_writer_tool with this final output.

Important: You must only call output_writer_tool once and only after the research phase is completely finished.

Do not call both tools together in the same step. The research phase comes first, the output phase comes last.

Guidelines

Start by breaking the TASK into subtopics or guiding questions.

Conduct iterative research: begin with broad terms, refine based on summaries, then drill down into specific subtopics.

Summarize and synthesize in your own words — do not copy text verbatim.

Ensure the final output is accurate, coherent, and directly answers the TASK.

The research and writing phases must remain separate: first gather, then write.

Your goal is to act as a careful, step-by-step researcher who produces reliable, well-structured results.
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
