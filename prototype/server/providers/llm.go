package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// LLMProvider defines the interface for language model backends.
type LLMProvider interface {
	Complete(ctx context.Context, prompt, systemPrompt string) (string, error)
}

// NewProvider returns an LLMProvider based on MITRAN_LLM_PROVIDER env.
// Supported values: bedrock (default), ollama, openai.
func NewProvider(name string) (LLMProvider, error) {
	if name == "" {
		name = os.Getenv("MITRAN_LLM_PROVIDER")
	}
	if name == "" {
		name = "bedrock"
	}
	switch name {
	case "bedrock":
		return &BedrockProvider{
			Region: envOrDefault("AWS_REGION", "us-east-1"),
			Model:  envOrDefault("MITRAN_BEDROCK_MODEL", "anthropic.claude-sonnet-4-20250514"),
		}, nil
	case "ollama":
		return &OllamaProvider{
			BaseURL: envOrDefault("MITRAN_OLLAMA_URL", "http://localhost:11434"),
			Model:   envOrDefault("MITRAN_OLLAMA_MODEL", "llama3"),
		}, nil
	case "openai":
		key := os.Getenv("OPENAI_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY required for openai provider")
		}
		return &OpenAIProvider{
			APIKey: key,
			Model:  envOrDefault("MITRAN_OPENAI_MODEL", "gpt-4o"),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", name)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// BedrockProvider uses AWS Bedrock Runtime API.
type BedrockProvider struct {
	Region string
	Model  string
}

func (b *BedrockProvider) Complete(ctx context.Context, prompt, systemPrompt string) (string, error) {
	// TODO: Implement full SigV4-signed Bedrock InvokeModel call.
	// For now, return a placeholder indicating Bedrock is configured.
	_ = ctx
	return fmt.Sprintf("[bedrock:%s] would process: %s", b.Model, truncate(prompt, 100)), nil
}

// OllamaProvider calls local Ollama HTTP API.
type OllamaProvider struct {
	BaseURL string
	Model   string
}

func (o *OllamaProvider) Complete(ctx context.Context, prompt, systemPrompt string) (string, error) {
	body := map[string]interface{}{
		"model":  o.Model,
		"prompt": prompt,
		"system": systemPrompt,
		"stream": false,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("ollama marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.BaseURL+"/api/generate", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ollama decode: %w", err)
	}
	return result.Response, nil
}

// OpenAIProvider calls the OpenAI Chat Completions API.
type OpenAIProvider struct {
	APIKey string
	Model  string
}

func (op *OpenAIProvider) Complete(ctx context.Context, prompt, systemPrompt string) (string, error) {
	messages := []map[string]string{}
	if systemPrompt != "" {
		messages = append(messages, map[string]string{"role": "system", "content": systemPrompt})
	}
	messages = append(messages, map[string]string{"role": "user", "content": prompt})

	body := map[string]interface{}{
		"model":    op.Model,
		"messages": messages,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("openai marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+op.APIKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("openai decode: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai: no choices returned")
	}
	return result.Choices[0].Message.Content, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
