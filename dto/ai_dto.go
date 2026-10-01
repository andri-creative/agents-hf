// dto/ai_dto.go — Data Transfer Objects untuk endpoint AI Generation
// Mendefinisikan struktur request/response untuk POST /v1/messages

package dto

// AIGenerateRequest adalah body request untuk endpoint POST /v1/messages
type AIGenerateRequest struct {
	Model  string `json:"model"  validate:"omitempty"`               // Opsional, fallback ke HF_MODEL dari ENV
	Prompt string `json:"prompt" validate:"required,min=1,max=4000"` // Prompt wajib diisi
}

// AIGenerateResponse adalah response sukses dari endpoint AI Generate
type AIGenerateResponse struct {
	Model    string `json:"model"`    // Model yang digunakan
	Response string `json:"response"` // Teks yang dihasilkan oleh AI
}
