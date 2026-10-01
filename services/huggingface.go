// services/huggingface.go — Service untuk integrasi dengan HuggingFace Router API
// Menggunakan format OpenAI-compatible Chat Completions (POST /v1/chat/completions)
// Dokumentasi: https://huggingface.co/docs/inference-providers
//
// Kompatibel dengan OpenAI SDK, 9router, One-API, NextChat, LibreChat, dll.

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ai-generate-api/config"
)

// HuggingFaceService mengelola komunikasi dengan HuggingFace Router API
type HuggingFaceService struct {
	Env        *config.EnvConfig
	HTTPClient *http.Client
}

// NewHuggingFaceService membuat instance service baru dengan timeout 60 detik
func NewHuggingFaceService(env *config.EnvConfig) *HuggingFaceService {
	return &HuggingFaceService{
		Env: env,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Structs untuk OpenAI-compatible Chat Completions API
// ─────────────────────────────────────────────────────────────────────────────

// chatMessage merepresentasikan satu pesan dalam percakapan (role + content)
type chatMessage struct {
	Role    string `json:"role"`    // "system", "user", atau "assistant"
	Content string `json:"content"` // Isi pesan
}

// chatCompletionRequest adalah body request ke /v1/chat/completions
type chatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

// chatChoice adalah satu pilihan dalam array choices dari response OpenAI
type chatChoice struct {
	Message      chatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
}

// chatCompletionResponse adalah struktur response dari /v1/chat/completions
type chatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []chatChoice `json:"choices"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Generate — memanggil HuggingFace Router API dan mengembalikan teks AI (string)
// Digunakan oleh endpoint POST /v1/messages
// ─────────────────────────────────────────────────────────────────────────────

func (s *HuggingFaceService) Generate(model, prompt string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", s.Env.HFBaseURL)

	reqBody := chatCompletionRequest{
		Model: model,
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("gagal marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("gagal membuat HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Env.HFApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal terhubung ke HuggingFace Router API: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal membaca response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HuggingFace Router API error (status %d): %s", resp.StatusCode, string(respBodyBytes))
	}

	var chatResp chatCompletionResponse
	if err := json.Unmarshal(respBodyBytes, &chatResp); err != nil {
		return "", fmt.Errorf("gagal parse response: %w (raw: %s)", err, string(respBodyBytes))
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("HuggingFace mengembalikan response kosong (choices=[])")
	}

	content := chatResp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("konten response kosong")
	}

	return content, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// ForwardChatCompletions — meneruskan request Chat Completions langsung ke HF Router
// Digunakan oleh endpoint POST /v1/chat/completions untuk 9router / OpenAI SDK
// ─────────────────────────────────────────────────────────────────────────────

func (s *HuggingFaceService) ForwardChatCompletions(rawBody []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s/chat/completions", s.Env.HFBaseURL)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawBody))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal membuat request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Env.HFApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("gagal terhubung ke HuggingFace Router: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal membaca response body: %w", err)
	}

	return respBodyBytes, resp.StatusCode, nil
}
