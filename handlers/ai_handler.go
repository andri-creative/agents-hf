// handlers/ai_handler.go — Handler HTTP untuk endpoint AI Generation
// Menangani: POST /v1/messages (Protected — API Key required)
// Mendelegasikan pemanggilan HuggingFace ke services/huggingface.go

package handlers

import (
	"net/http"

	"ai-generate-api/config"
	"ai-generate-api/dto"
	"ai-generate-api/services"
	"ai-generate-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AIHandler menyimpan dependensi yang dibutuhkan oleh handler AI
type AIHandler struct {
	HFService *services.HuggingFaceService
	Env       *config.EnvConfig
}

// NewAIHandler membuat instance AIHandler baru
func NewAIHandler(env *config.EnvConfig) *AIHandler {
	return &AIHandler{
		HFService: services.NewHuggingFaceService(env),
		Env:       env,
	}
}

// Generate menangani POST /v1/messages (Protected — API Key required)
// Menerima model dan prompt, meneruskan ke HuggingFace, dan mengembalikan response AI
func (h *AIHandler) Generate(c *gin.Context) {
	var req dto.AIGenerateRequest

	// Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", err.Error())
		return
	}

	// Validasi input
	aiValidate := validator.New()
	if err := aiValidate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorDetails := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorDetails[fieldErr.Field()] = fieldErr.Tag()
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi input gagal", errorDetails)
		return
	}

	// Jika model tidak disertakan, gunakan model default dari ENV
	model := req.Model
	if model == "" {
		model = h.Env.HFModel
	}

	// Panggil HuggingFace service untuk generate teks
	generatedText, err := h.HFService.Generate(model, req.Prompt)
	if err != nil {
		// Periksa apakah error dari HuggingFace API (502 Bad Gateway)
		utils.ErrorResponse(c, http.StatusBadGateway, "Gagal mendapatkan response dari HuggingFace", err.Error())
		return
	}

	// Return response sukses
	utils.SuccessResponse(c, http.StatusOK, "AI generation berhasil", dto.AIGenerateResponse{
		Model:    model,
		Response: generatedText,
	})
}
