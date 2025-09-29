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
