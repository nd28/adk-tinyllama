package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"strings"

	"google.golang.org/genai"

	"google.golang.org/adk/model"
)

type ollamaModel struct {
	name    string
	baseURL string
	client  *http.Client
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system,omitempty"`
}

type ollamaResponse struct {
	Model         string `json:"model"`
	CreatedAt     string `json:"created_at"`
	Response      string `json:"response"`
	Done          bool   `json:"done"`
	TotalDuration int64  `json:"total_duration,omitempty"`
	EvalCount     int    `json:"eval_count,omitempty"`
	EvalDuration  int64  `json:"eval_duration,omitempty"`
}

func NewModel(modelName string, baseURL string) model.LLM {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:11434"
	}
	return &ollamaModel{
		name:    modelName,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (m *ollamaModel) Name() string {
	return m.name
}

func (m *ollamaModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		prompt := extractPrompt(req)

		body := ollamaRequest{
			Model:  m.name,
			Prompt: prompt,
			Stream: stream,
		}

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			yield(nil, fmt.Errorf("marshal request: %w", err))
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/api/generate", bytes.NewReader(bodyBytes))
		if err != nil {
			yield(nil, fmt.Errorf("create request: %w", err))
			return
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := m.client.Do(httpReq)
		if err != nil {
			yield(nil, fmt.Errorf("send request: %w", err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			yield(nil, fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(body)))
			return
		}

		if stream {
			m.handleStream(resp, yield)
		} else {
			m.handleNonStream(resp, yield)
		}
	}
}

func (m *ollamaModel) handleNonStream(resp *http.Response, yield func(*model.LLMResponse, error) bool) {
	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		yield(nil, fmt.Errorf("decode response: %w", err))
		return
	}

	yield(&model.LLMResponse{
		Content: &genai.Content{
			Parts: []*genai.Part{
				{Text: ollamaResp.Response},
			},
			Role: "model",
		},
	}, nil)
}

func (m *ollamaModel) handleStream(resp *http.Response, yield func(*model.LLMResponse, error) bool) {
	decoder := json.NewDecoder(resp.Body)
	var fullText strings.Builder

	for {
		var ollamaResp ollamaResponse
		if err := decoder.Decode(&ollamaResp); err != nil {
			if err == io.EOF {
				break
			}
			yield(nil, fmt.Errorf("decode stream: %w", err))
			return
		}

		fullText.WriteString(ollamaResp.Response)

		yield(&model.LLMResponse{
			Content: &genai.Content{
				Parts: []*genai.Part{
					{Text: ollamaResp.Response},
				},
				Role: "model",
			},
			Partial: !ollamaResp.Done,
		}, nil)

		if ollamaResp.Done {
			return
		}
	}
}

func extractPrompt(req *model.LLMRequest) string {
	if req == nil || req.Contents == nil {
		return ""
	}
	var parts []string
	for _, content := range req.Contents {
		for _, part := range content.Parts {
			if part.Text != "" {
				parts = append(parts, part.Text)
			}
		}
	}
	return strings.Join(parts, "\n")
}
