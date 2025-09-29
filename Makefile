.PHONY: run-ollama run-gemini

file ?= default.txt

run-ollama:
	 MODEL_TYPE=ollama go run ./... -file $(file)

run-gemini:
	 MODEL_TYPE=gemini go run ./... -file $(file)
