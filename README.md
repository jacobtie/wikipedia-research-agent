# Wikipedia Research Agent

> [!WARNING]
> As an educator, I strictly condemn the use of this repository to complete academic assignments. Beyond the academic integrity violations, this would simply be a bad idea as this agent can only read from Wikipedia. Furthermore, this is a toy repository and would need a lot of work in error handling and prompt engineering before you could seriously rely on it in any way.

The Wikipedia Research Agent is a toy LLM backed agent that performs research on different Wikipedia pages and then writes a paper based on its research. Besides the use of the [Google Gemini SDK for Golang](https://github.com/googleapis/go-genai), there are no third party dependencies and no agent framework used here. Instead, the implementation is written in pure Golang. This was done as an exercise in building a simple agentic system.

## Setup

### Ollama

To run this gemini with Ollama (which does not work terribly well, but is free), you must have ollama running locally with the `llama3.2` model installed.

> [!TIP]
> You can you a different Ollama model, but note that not all models support tool calling.

To run ollama locally, run:

```bash
docker compose -f docker-compose.dev.yml up
```

and then, to install `llama3.2`, run:

```bash
docker exec ollama ollama pull llama3.2
```

This docker compose stack utilizes a volume, so you only need to pull the `llama3.2` model once.

### Gemini

To run this agent with Gemini, you need to export a `GEMINI_API_KEY` environment variable with your Gemini API key. Note that because of the large context used by the Wikipedia articles and the frequent API access, you will need to use a paid tier account as the free tier will not suffice.

### Input File

Put your prompt into the `input/` directory in a `.txt` file. If you do not specify a prompt when running the agent, the agent will use the prompt in the `default.txt` file.

## Running the Agent

This repository provides a `Makefile` to aid in running the agent.

To run the agent with ollama, run:

```bash
make run-ollama file=input-file-here.txt
```

To run the agent with Gemini, run:

```bash
make run-gemini file=input-file-here.txt
```

If the agent finishes without error (error handling should be improved before this is production ready), then it should output its answer into a file with the UNIX timestamp in the `output` directory.

## Architecture

This agent makes use of the following two tools with the following architecture.

![Wikipedia Research Agent Architecture](https://github.com/jacobtie/wikipedia-research-agent/blob/main/assets/wikipedia-research-agent-architecture.png)

### Research Tool

The research tool is itself an LLM workflow. The idea is that using each Wikipedia page in its entirety would be way too much context for the agent. With this workflow, the LLM only sees one full page at a time and is able to build relevant summaries of the pages for the agent to use in its research. The agent can call the research tool multiple times as it sees fit until it is satisfied with its research.

1. Uses the [Wikipedia API](https://en.wikipedia.org/w/api.php) to get a list of pages related to a search term determined by the agent -- each page has a title and potentially relevant snippet of the page
2. Calls the LLM for each result to determine whether it is relevant to the research task based on the title and snippet
3. Gets the wikitext context for each page via the Wikipedia API
4. Calls the LLM for each page to summarize it with respect to the research task

### Output Tool

The output tool is very simple, it simply outputs the results of the research into an output file with the UNIX timestamp in the `output` directory. Note that this tool is not a very realistic tool and instead is used in this toy as an example of a tool that impacts the environment. In a real, production ready system, the agent would likely respond with the raw text and then the caller of the agent would do something with that text, such as write it to an output file
