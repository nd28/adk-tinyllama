# ADK TinyLlama

A Go application demonstrating how to use Google's Agent Development Kit (ADK) with a local Ollama model (TinyLlama).

## Overview

This project shows how to create a local AI assistant using:
- **Google ADK** - Agent Development Kit for building AI agents
- **Ollama** - Run large language models locally
- **TinyLlama** - A compact LLM that runs locally via Ollama

## Prerequisites

1. **Go 1.25+** - Install from [go.dev](https://go.dev/dl/)
2. **Ollama** - Install from [ollama.com](https://ollama.com)
3. **TinyLlama model** - Pull with: `ollama pull tinyllama`

## Setup

1. Install Ollama and start the service:
   ```bash
   ollama serve
   ```

2. In another terminal, pull the TinyLlama model:
   ```bash
   ollama pull tinyllama
   ```

3. Clone this repository and run:
   ```bash
   go run main.go
   ```

## Project Structure

```
.
├── main.go           # Main application entry point
├── go.mod            # Go module dependencies
├── go.sum            # Go module checksums
└── ollama/
    └── model.go     # Ollama model implementation for ADK
```

## How It Works

1. **Ollama Model** (`ollama/model.go`): Implements the ADK's `model.LLM` interface to connect to a local Ollama instance
2. **Agent** (`main.go`): Creates an LLM agent using ADK's `llmagent`
3. **Runner**: Handles running the agent with user prompts

## Example Usage

The application runs three test prompts:
1. "What is 2 + 2?"
2. "Write a one-line haiku about Go programming."
3. "What does IFSC stand for in banking?"

## Using Different Models

### Qwen3.5

To use Qwen3.5 with Ollama:

1. Pull the model:
   ```bash
   ollama pull qwen3.5
   ```

2. Update `main.go` to use the model:
   ```go
   llm := ollama.NewModel("qwen3.5", "http://127.0.0.1:11434")
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

For better results, you can use a larger variant like `qwen3.5:7b` or `qwen3.5:14b`:
```bash
ollama pull qwen3.5:7b
```

Then update `main.go`:
```go
llm := ollama.NewModel("qwen3.5:7b", "http://127.0.0.1:11434")
```

### Other Models

To use any other model, change the model name in `main.go`:
```go
llm := ollama.NewModel("your-model-name", "http://127.0.0.1:11434")
```

Then pull it with: `ollama pull your-model-name`

## Dependencies

- `google.golang.org/adk` - Google Agent Development Kit
- `google.golang.org/genai` - Google GenAI API
