// dto/ai_dto.go — Data Transfer Objects untuk AI Generation (Custom, OpenAI, dan Anthropic Compatible)
// Mendukung 9router, One-API, NextChat, Claude Code, dll.

package dto

// AIGenerateRequest adalah body request untuk endpoint POST /v1/messages (Single Prompt)
type AIGenerateRequest struct {
	Model  string `json:"model"  validate:"omitempty"`               // Opsional, fallback ke HF_MODEL dari ENV
	Prompt string `json:"prompt" validate:"required,min=1,max=4000"` // Prompt wajib diisi
}

// AIGenerateResponse adalah response sukses dari endpoint POST /v1/messages (Single Prompt format)
type AIGenerateResponse struct {
	Model    string `json:"model"`    // Model yang digunakan
	Response string `json:"response"` // Teks yang dihasilkan oleh AI
}

// ─────────────────────────────────────────────────────────────────────────────
// Model List DTO (Kompatibel dengan OpenAI dan Anthropic / 9router)
// ─────────────────────────────────────────────────────────────────────────────

// ModelItem merepresentasikan satu model dalam format OpenAI & Anthropic
type ModelItem struct {
	ID          string `json:"id"`
	Type        string `json:"type,omitempty"`         // "model" untuk Anthropic
	Object      string `json:"object,omitempty"`       // "model" untuk OpenAI
	DisplayName string `json:"display_name,omitempty"` // untuk Anthropic UI / 9router
	Created     int64  `json:"created"`
	OwnedBy     string `json:"owned_by"`
}

// ModelListResponse adalah format response standar OpenAI & Anthropic untuk GET /v1/models
type ModelListResponse struct {
	Object string      `json:"object,omitempty"`
	Data   []ModelItem `json:"data"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Anthropic Messages API DTOs (Dibutuhkan oleh 9router Anthropic Compatible)
// ─────────────────────────────────────────────────────────────────────────────

// AnthropicContentBlock merepresentasikan blok konten teks dalam pesan Anthropic
type AnthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// AnthropicUsage merepresentasikan info penggunaan token
type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// AnthropicMessageResponse adalah response format resmi Anthropic Messages API
type AnthropicMessageResponse struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"` // "message"
	Role       string                  `json:"role"` // "assistant"
	Model      string                  `json:"model"`
	Content    []AnthropicContentBlock `json:"content"`
	StopReason string                  `json:"stop_reason"`
	Usage      AnthropicUsage          `json:"usage"`
}
