// services/huggingface.go — Service untuk integrasi dengan HuggingFace Router API
// Menggunakan format OpenAI-compatible Chat Completions (POST /v1/chat/completions)
// Dokumentasi: https://huggingface.co/docs/inference-providers
//
// Kompatibel dengan OpenAI SDK, 9router, One-API, NextChat, dll.
// [UPGRADE v3] — Memvalidasi model aktif dari database sebelum request diteruskan ke HuggingFace.
// [UPGRADE 3] — API key & base URL dibaca dari tabel configs (fallback ke ENV).

package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"ai-generate-api/config"
	"ai-generate-api/models"

	"gorm.io/gorm"
)

var ErrModelUnavailable = errors.New("model tidak tersedia")

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
	if err := validateActiveModel(model); err != nil {
		return "", err
	}

	// [UPGRADE 3] — Ambil config HF dari DB (fallback ENV) setiap request.
	apiKey, _, baseURL := GetHuggingFaceConfig()
	url := fmt.Sprintf("%s/chat/completions", baseURL)

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

	req.Header.Set("Authorization", "Bearer "+apiKey)
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
	var reqBody struct {
		Model string `json:"model"`
	}
	if err := json.Unmarshal(rawBody, &reqBody); err != nil {
		return marshalInvalidModelResponse("Format request ke HuggingFace tidak valid"), http.StatusBadRequest, nil
	}

	if err := validateActiveModel(reqBody.Model); err != nil {
		return marshalInvalidModelResponse(err.Error()), http.StatusBadRequest, nil
	}

	// [UPGRADE 3] — Ambil config HF dari DB (fallback ENV) setiap request.
	apiKey, _, baseURL := GetHuggingFaceConfig()
	url := fmt.Sprintf("%s/chat/completions", baseURL)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawBody))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal membuat request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
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

// IsModelUnavailableError membantu handler membedakan error validasi model dan error upstream.
func IsModelUnavailableError(err error) bool {
	return errors.Is(err, ErrModelUnavailable)
}

func validateActiveModel(model string) error {
	if model == "" {
		return fmt.Errorf("%w: slug model wajib diisi", ErrModelUnavailable)
	}

	var aiModel models.AIModel
	err := config.DB.Where("slug = ? AND is_active = ?", model, true).First(&aiModel).Error
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: model '%s' tidak aktif atau tidak ditemukan", ErrModelUnavailable, model)
	}

	return fmt.Errorf("gagal memvalidasi model '%s': %w", model, err)
}

func marshalInvalidModelResponse(message string) []byte {
	respBytes, err := json.Marshal(map[string]interface{}{
		"error": map[string]interface{}{
			"type":    "invalid_request_error",
			"message": message,
		},
	})
	if err != nil {
		return []byte(`{"error":{"type":"invalid_request_error","message":"Model tidak valid"}}`)
	}

	return respBytes
}
