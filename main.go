package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"

	"adk-tinyllama/ollama"
)

func main() {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("  ADK Go 1.0 + TinyLlama (Local)")
	fmt.Println("========================================")
	fmt.Println()

	// Create Ollama model (TinyLlama running locally on port 11434)
	llm := ollama.NewModel("tinyllama", "http://127.0.0.1:11434")
	fmt.Printf("Model: %s\n\n", llm.Name())

	// Create LLM agent using ADK's llmagent
	myAgent, err := llmagent.New(llmagent.Config{
		Name:        "tinyllama_assistant",
		Model:       llm,
		Description: "A local AI assistant powered by TinyLlama via Ollama",
		Instruction: "You are a helpful assistant. Keep responses brief and concise.",
	})
	if err != nil {
		log.Fatal("Failed to create agent:", err)
	}

	// Create session service
	sessionSvc := session.InMemoryService()

	// Create runner
	r, err := runner.New(runner.Config{
		AppName:        "adk-tinyllama",
		Agent:          myAgent,
		SessionService: sessionSvc,
	})
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	// Create a session
	createResp, err := sessionSvc.Create(ctx, &session.CreateRequest{
		AppName: "adk-tinyllama",
		UserID:  "user-1",
	})
	if err != nil {
		log.Fatal("Failed to create session:", err)
	}
	sessID := createResp.Session.ID()

	// Test prompts
	prompts := []string{
		"What is 2 + 2?",
		"Write a one-line haiku about Go programming.",
		"What does IFSC stand for in banking?",
	}

	for i, prompt := range prompts {
		fmt.Printf("[%d] You: %s\n", i+1, prompt)
		fmt.Printf("    Agent: ")

		userMsg := &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{{Text: prompt}},
		}

		events := r.Run(ctx, "user-1", sessID, userMsg, agent.RunConfig{})

		for event, err := range events {
			if err != nil {
				fmt.Printf("\n    Error: %v\n", err)
				break
			}
			if event != nil && event.Content != nil {
				for _, part := range event.Content.Parts {
					if part.Text != "" {
						fmt.Print(part.Text)
					}
				}
			}
		}
		fmt.Println()
		fmt.Println()
	}

	fmt.Println("Done!")
}
