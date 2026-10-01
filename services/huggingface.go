// services/huggingface.go — Service untuk integrasi dengan HuggingFace Router API
// Menggunakan format OpenAI-compatible Chat Completions (POST /v1/chat/completions)
// Dokumentasi: https://huggingface.co/docs/inference-providers
//
// Contoh Python equivalent:
//   client = OpenAI(base_url="https://router.huggingface.co/v1", api_key=HF_TOKEN)
//   completion = client.chat.completions.create(model="zai-org/GLM-5.3:novita", messages=[...])

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
// Generate — memanggil HuggingFace Router API dan mengembalikan teks AI
// ─────────────────────────────────────────────────────────────────────────────

// Generate memanggil POST /v1/chat/completions ke HuggingFace Router
// model: nama model (misal: "zai-org/GLM-5.3:novita")
// prompt: teks input dari user — dikemas sebagai pesan role "user"
func (s *HuggingFaceService) Generate(model, prompt string) (string, error) {
	// Endpoint OpenAI-compatible: {HF_BASE_URL}/chat/completions
	// HF_BASE_URL sudah mengandung /v1, jadi: https://router.huggingface.co/v1/chat/completions
	url := fmt.Sprintf("%s/chat/completions", s.Env.HFBaseURL)

	// Susun request body dalam format OpenAI Chat Completions
	reqBody := chatCompletionRequest{
		Model: model,
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Serialize request body ke JSON
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("gagal marshal request body: %w", err)
	}

	// Buat HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("gagal membuat HTTP request: %w", err)
	}

	// Set headers — sama persis dengan Python OpenAI client
	req.Header.Set("Authorization", "Bearer "+s.Env.HFApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Kirim request ke HuggingFace Router
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal terhubung ke HuggingFace Router API: %w", err)
	}
	defer resp.Body.Close()

	// Baca response body
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal membaca response body: %w", err)
	}

	// Cek status code — jika bukan 200, kembalikan error dengan detail dari HF
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HuggingFace Router API error (status %d): %s", resp.StatusCode, string(respBodyBytes))
	}

	// Parse response dalam format OpenAI Chat Completions
	var chatResp chatCompletionResponse
	if err := json.Unmarshal(respBodyBytes, &chatResp); err != nil {
		return "", fmt.Errorf("gagal parse response: %w (raw: %s)", err, string(respBodyBytes))
	}

	// Pastikan ada minimal 1 choice dalam response
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("HuggingFace mengembalikan response kosong (choices=[])")
	}

	// Ambil konten pesan dari choice pertama (setara dengan choices[0].message.content di Python)
	content := chatResp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("konten response kosong")
	}

	return content, nil
}
