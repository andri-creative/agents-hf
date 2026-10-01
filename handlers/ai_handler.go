// handlers/ai_handler.go — Handler HTTP untuk endpoint AI Generation (OpenAI & Anthropic Compatible)
// Mendukung:
// - POST /v1/messages (Dual-mode: Single prompt & Anthropic Messages API untuk 9router)
// - GET  /v1/models & /models (OpenAI & Anthropic model list untuk 9router "Import from /models")
// - POST /v1/chat/completions & /chat/completions (OpenAI chat completions)

package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

// ─────────────────────────────────────────────────────────────────────────────
// 1. POST /v1/messages — Endpoint Dual-Mode (Single Prompt & Anthropic Format)
// 9router menggunakan mode "Anthropic Compatible" yang mengirim field "messages"
// ─────────────────────────────────────────────────────────────────────────────

func (h *AIHandler) Generate(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca request body", err.Error())
		return
	}

	// Cek apakah request berupa JSON umum
	var genericPayload map[string]interface{}
	if err := json.Unmarshal(rawData, &genericPayload); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format JSON tidak valid", err.Error())
		return
	}

	// ── SKENARIO A: Request dari 9router (Anthropic Compatible dengan "messages") ──
	if rawMessages, hasMessages := genericPayload["messages"]; hasMessages && rawMessages != nil {
		h.handleAnthropicMessages(c, genericPayload)
		return
	}

	// ── SKENARIO B: Request Standar dengan "prompt" ──
	var req dto.AIGenerateRequest
	if err := json.Unmarshal(rawData, &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", err.Error())
		return
	}

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

	model := req.Model
	if model == "" {
		model = h.Env.HFModel
	}

	generatedText, err := h.HFService.Generate(model, req.Prompt)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadGateway, "Gagal mendapatkan response dari HuggingFace", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "AI generation berhasil", dto.AIGenerateResponse{
		Model:    model,
		Response: generatedText,
	})
}

// handleAnthropicMessages menangani request dengan format Anthropic Messages API (dari 9router)
func (h *AIHandler) handleAnthropicMessages(c *gin.Context, payload map[string]interface{}) {
	model, _ := payload["model"].(string)
	if model == "" {
		model = h.Env.HFModel
	}

	// Ekstrak pesan dari payload Anthropic untuk diteruskan ke HuggingFace Router
	rawMessagesList, ok := payload["messages"].([]interface{})
	if !ok || len(rawMessagesList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"type": "error",
			"error": gin.H{
				"type":    "invalid_request_error",
				"message": "Field 'messages' harus berupa array dan tidak boleh kosong",
			},
		})
		return
	}

	// Normalisasi messages agar kompatibel dengan format chat HuggingFace Router
	type simpleChatMsg struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	var normalizedMessages []simpleChatMsg
	for _, m := range rawMessagesList {
		msgMap, ok := m.(map[string]interface{})
		if !ok {
			continue
		}
		role, _ := msgMap["role"].(string)
		if role == "" {
			role = "user"
		}

		var textContent string
		switch v := msgMap["content"].(type) {
		case string:
			textContent = v
		case []interface{}:
			// Format blok Anthropic: [{"type": "text", "text": "..."}]
			for _, block := range v {
				if blockMap, ok := block.(map[string]interface{}); ok {
					if t, ok := blockMap["text"].(string); ok {
						textContent += t
					}
				}
			}
		}

		normalizedMessages = append(normalizedMessages, simpleChatMsg{
			Role:    role,
			Content: textContent,
		})
	}

	// Siapkan request OpenAI-compatible untuk HuggingFace Router
	hfPayload := map[string]interface{}{
		"model":    model,
		"messages": normalizedMessages,
	}
	if maxTokens, ok := payload["max_tokens"].(float64); ok && maxTokens > 0 {
		hfPayload["max_tokens"] = int(maxTokens)
	}

	hfBytes, err := json.Marshal(hfPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type": "error",
			"error": gin.H{
				"type":    "api_error",
				"message": "Gagal menyiapkan payload untuk HuggingFace: " + err.Error(),
			},
		})
		return
	}

	// Panggil HuggingFace Router
	respBytes, statusCode, err := h.HFService.ForwardChatCompletions(hfBytes)
	if err != nil || statusCode != http.StatusOK {
		c.Data(statusCode, "application/json", respBytes)
		return
	}

	// Parse response dari HuggingFace Router untuk dikonversi ke format Anthropic
	var hfResp struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(respBytes, &hfResp); err != nil || len(hfResp.Choices) == 0 {
		// Jika gagal parse, kembalikan response apa adanya
		c.Data(statusCode, "application/json", respBytes)
		return
	}

	// Format response Anthropic resmi yang diharapkan oleh 9router
	randBytes := make([]byte, 12)
	_, _ = rand.Read(randBytes)
	msgID := fmt.Sprintf("msg_%s", hex.EncodeToString(randBytes))

	anthropicResp := dto.AnthropicMessageResponse{
		ID:    msgID,
		Type:  "message",
		Role:  "assistant",
		Model: model,
		Content: []dto.AnthropicContentBlock{
			{
				Type: "text",
				Text: hfResp.Choices[0].Message.Content,
			},
		},
		StopReason: "end_turn",
		Usage: dto.AnthropicUsage{
			InputTokens:  hfResp.Usage.PromptTokens,
			OutputTokens: hfResp.Usage.CompletionTokens,
		},
	}

	c.JSON(http.StatusOK, anthropicResp)
}

// ─────────────────────────────────────────────────────────────────────────────
// 2. GET /v1/models & GET /models — Model List untuk 9router "Import from /models"
// ─────────────────────────────────────────────────────────────────────────────

func (h *AIHandler) GetModels(c *gin.Context) {
	now := time.Now().Unix()
	defaultModel := h.Env.HFModel

	models := []dto.ModelItem{
		{
			ID:          defaultModel,
			Type:        "model",
			Object:      "model",
			DisplayName: defaultModel,
			Created:     now,
			OwnedBy:     "huggingface",
		},
	}

	alternatives := []string{
		"zai-org/GLM-5.3:novita",
		"mistralai/Mistral-7B-Instruct-v0.2",
		"meta-llama/Meta-Llama-3-8B-Instruct",
		"Qwen/Qwen2.5-72B-Instruct",
	}

	for _, alt := range alternatives {
		if alt != defaultModel {
			models = append(models, dto.ModelItem{
				ID:          alt,
				Type:        "model",
				Object:      "model",
				DisplayName: alt,
				Created:     now,
				OwnedBy:     "huggingface",
			})
		}
	}

	c.JSON(http.StatusOK, dto.ModelListResponse{
		Object: "list",
		Data:   models,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// 3. GET /v1/models/:model & GET /models/:model
// ─────────────────────────────────────────────────────────────────────────────

func (h *AIHandler) GetModelDetail(c *gin.Context) {
	modelID := c.Param("model")
	if modelID == "" {
		modelID = h.Env.HFModel
	}

	c.JSON(http.StatusOK, dto.ModelItem{
		ID:          modelID,
		Type:        "model",
		Object:      "model",
		DisplayName: modelID,
		Created:     time.Now().Unix(),
		OwnedBy:     "huggingface",
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// 4. POST /v1/chat/completions & /chat/completions — OpenAI Standard Chat
// ─────────────────────────────────────────────────────────────────────────────

func (h *AIHandler) ChatCompletions(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Gagal membaca request body: " + err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(rawData, &payload); err == nil {
		if m, ok := payload["model"].(string); !ok || m == "" {
			payload["model"] = h.Env.HFModel
			if modifiedData, err := json.Marshal(payload); err == nil {
				rawData = modifiedData
			}
		}
	}

	respBytes, statusCode, err := h.HFService.ForwardChatCompletions(rawData)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "api_error",
			},
		})
		return
	}

	c.Data(statusCode, "application/json; charset=utf-8", respBytes)
}
