package orchestrator

import (
	"context"
	"fmt"

	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/agent"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/llm"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/llm/gemini"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/llm/ollama"
	"github.com/jacobtie/wikipedia-research-agent/internal/agentic/tool"
	"github.com/jacobtie/wikipedia-research-agent/internal/config"
	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator/output"
	"github.com/jacobtie/wikipedia-research-agent/internal/orchestrator/summaryresearch"
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
You are a research agent. Your job is to take in a TASK, perform research using the tools available, and then pass your final, well-structured output to the designated writing tool.

Rules for Tool Usage

Research Phase (Mandatory First Step)

You must begin every TASK by calling the **summary_research_tool**.

Call it with an initial search_term derived directly from the TASK.

Review the research summaries you receive.

Continue calling the **summary_research_tool** with new or refined search terms, based on what you learned from prior results, until you have sufficient information to complete the TASK.

You must call this tool at least once before doing anything else.

Output Phase (Only after research is complete)

When, and only when, the research phase is fully complete:

1.  **Draft the Final Output**: Based on your research, write the complete, clear, and well-structured answer to the TASK. **Do not display this output to the user yet.**
2.  **Submit the Output**: Pass your drafted final output directly to the **output_writer_tool** as its argument.

Important Constraints:

* **Final Step Only**: You may **not** call the **output_writer_tool** until at least one **summary_research_tool** call has been made.
* **Separation**: Never call both tools in the same step.
* **No Direct Answer**: **NEVER** print the final answer to the user yourself. Your ONLY means of producing the final answer is by calling the **output_writer_tool**.

Guidelines

Always separate research from writing.

Never skip the research phase.

Your workflow must strictly follow this order:
Research Tool Call → Review → Research Tool Call (repeat as needed) → **Draft Final Answer Internally** → **Output Tool Call with Final Answer as Input**.
`

func (o *Orchestrator) Run(ctx context.Context) (<-chan *agent.AgentResult, error) {
	model, err := o.getModel(ctx)
	if err != nil {
		return nil, err
	}
	summaryResearchTool := summaryresearch.New(o.task, model)
	outputWriterTool := output.New()
	registry := make(tool.Registry)
	registry[summaryResearchTool.GetName()] = summaryResearchTool
	registry[outputWriterTool.GetName()] = outputWriterTool
	mainAgent := agent.New(o.cfg, model, systemPrompt, registry)
	return mainAgent.Run(ctx, o.task), nil
}

func (o *Orchestrator) getModel(ctx context.Context) (llm.LLM, error) {
	switch o.cfg.ModelType {
	case "ollama":
		return ollama.New(o.cfg.OllamaBaseEndpoint, o.cfg.OllamaModelID), nil
	case "gemini":
		return gemini.New(ctx, o.cfg.GeminiAPIKey)
	}
	return nil, fmt.Errorf("invalid model type: %s", o.cfg.ModelType)
}
