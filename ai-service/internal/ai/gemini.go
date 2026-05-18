package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openAIURL = "https://api.openai.com/v1/chat/completions"

type GeminiClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (g *GeminiClient) Chat(ctx context.Context, systemPrompt, history, userMessage string) (string, error) {
	messages := []openAIMessage{}

	if systemPrompt != "" {
		messages = append(messages, openAIMessage{Role: "system", Content: systemPrompt})
	}

	if history != "" {
		messages = append(messages, openAIMessage{
			Role:    "user",
			Content: "Previous conversation:\n" + history,
		}, openAIMessage{
			Role:    "assistant",
			Content: "Understood, I have the context.",
		})
	}

	messages = append(messages, openAIMessage{Role: "user", Content: userMessage})

	reqBody := openAIRequest{
		Model:       "gpt-4o-mini",
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var or openAIResponse
	if err := json.Unmarshal(respBody, &or); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if or.Error != nil {
		return "", fmt.Errorf("openai error: %s", or.Error.Message)
	}

	if len(or.Choices) == 0 {
		return "", fmt.Errorf("empty response from openai")
	}

	return or.Choices[0].Message.Content, nil
}
