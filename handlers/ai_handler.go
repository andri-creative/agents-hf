// handlers/ai_handler.go — Handler HTTP untuk endpoint AI Generation (OpenAI & Anthropic Compatible)
// Mendukung:
// - POST /v1/messages (Dual-mode: Single prompt & Anthropic Messages API untuk 9router dengan full SSE streaming)
// - GET  /v1/models & /models (OpenAI & Anthropic model list untuk 9router "Import from /models")
// - POST /v1/chat/completions & /chat/completions (OpenAI chat completions dengan full SSE streaming untuk Trae/Cursor)
// [UPGRADE v3] — Mengambil model aktif dari database dan memvalidasi slug model sebelum dipakai.

package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ai-generate-api/config"
	"ai-generate-api/dto"
	"ai-generate-api/models"
	"ai-generate-api/services"
	"ai-generate-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
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
// ─────────────────────────────────────────────────────────────────────────────

func (h *AIHandler) Generate(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca request body", err.Error())
		return
	}

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
		if services.IsModelUnavailableError(err) {
			utils.ErrorResponse(c, http.StatusBadRequest, "Model tidak valid", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusBadGateway, "Gagal mendapatkan response dari HuggingFace", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "AI generation berhasil", dto.AIGenerateResponse{
		Model:    model,
		Response: generatedText,
	})
}

// handleAnthropicMessages menangani request Anthropic Messages API (dari 9router / Trae)
// Mendukung request streaming (stream: true) dengan format Server-Sent Events (SSE)
func (h *AIHandler) handleAnthropicMessages(c *gin.Context, payload map[string]interface{}) {
	model, _ := payload["model"].(string)
	if model == "" {
		model = h.Env.HFModel
	}

	// Cek apakah client meminta streaming (Trae dan 9router secara default meminta stream: true)
	isStream, _ := payload["stream"].(bool)

	// Ekstrak pesan dari payload Anthropic
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

	// Selalu minta respons lengkap (stream: false) dari HuggingFace Router
	hfPayload := map[string]interface{}{
		"model":    model,
		"messages": normalizedMessages,
		"stream":   false,
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
		c.Data(statusCode, "application/json", respBytes)
		return
	}

	answerText := hfResp.Choices[0].Message.Content

	randBytes := make([]byte, 12)
	_, _ = rand.Read(randBytes)
	msgID := fmt.Sprintf("msg_%s", hex.EncodeToString(randBytes))

	// ── JIKA CLIENT MEMINTA STREAMING (stream: true) ──
	// Mengirimkan Server-Sent Events (SSE) format Anthropic yang dinantikan 9router / Trae
	if isStream {
		c.Header("Content-Type", "text/event-stream; charset=utf-8")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")

		// 1. event: message_start
		msgStart, _ := json.Marshal(map[string]interface{}{
			"type": "message_start",
			"message": map[string]interface{}{
				"id":            msgID,
				"type":          "message",
				"role":          "assistant",
				"model":         model,
				"content":       []interface{}{},
				"stop_reason":   nil,
				"stop_sequence": nil,
				"usage": map[string]interface{}{
					"input_tokens":  hfResp.Usage.PromptTokens,
					"output_tokens": 1,
				},
			},
		})
		fmt.Fprintf(c.Writer, "event: message_start\ndata: %s\n\n", msgStart)
		c.Writer.Flush()

		// 2. event: content_block_start
		blockStart, _ := json.Marshal(map[string]interface{}{
			"type":  "content_block_start",
			"index": 0,
			"content_block": map[string]interface{}{
				"type": "text",
				"text": "",
			},
		})
		fmt.Fprintf(c.Writer, "event: content_block_start\ndata: %s\n\n", blockStart)
		c.Writer.Flush()

		// 3. event: content_block_delta (pecah teks dalam potongan kecil)
		words := strings.Fields(answerText)
		if len(words) == 0 {
			words = []string{answerText}
		}

		chunkSize := 3
		for i := 0; i < len(words); i += chunkSize {
			end := i + chunkSize
			if end > len(words) {
				end = len(words)
			}
			chunkText := strings.Join(words[i:end], " ")
			if end < len(words) {
				chunkText += " "
			}

			delta, _ := json.Marshal(map[string]interface{}{
				"type":  "content_block_delta",
				"index": 0,
				"delta": map[string]interface{}{
					"type": "text_delta",
					"text": chunkText,
				},
			})
			fmt.Fprintf(c.Writer, "event: content_block_delta\ndata: %s\n\n", delta)
			c.Writer.Flush()
			time.Sleep(5 * time.Millisecond)
		}

		// 4. event: content_block_stop
		blockStop, _ := json.Marshal(map[string]interface{}{
			"type":  "content_block_stop",
			"index": 0,
		})
		fmt.Fprintf(c.Writer, "event: content_block_stop\ndata: %s\n\n", blockStop)
		c.Writer.Flush()

		// 5. event: message_delta
		msgDelta, _ := json.Marshal(map[string]interface{}{
			"type": "message_delta",
			"delta": map[string]interface{}{
				"stop_reason":   "end_turn",
				"stop_sequence": nil,
			},
			"usage": map[string]interface{}{
				"output_tokens": hfResp.Usage.CompletionTokens,
			},
		})
		fmt.Fprintf(c.Writer, "event: message_delta\ndata: %s\n\n", msgDelta)
		c.Writer.Flush()

		// 6. event: message_stop
		msgStop, _ := json.Marshal(map[string]interface{}{
			"type": "message_stop",
		})
		fmt.Fprintf(c.Writer, "event: message_stop\ndata: %s\n\n", msgStop)
		c.Writer.Flush()
		return
	}

	// ── JIKA NON-STREAMING (stream: false) ──
	anthropicResp := dto.AnthropicMessageResponse{
		ID:    msgID,
		Type:  "message",
		Role:  "assistant",
		Model: model,
		Content: []dto.AnthropicContentBlock{
			{
				Type: "text",
				Text: answerText,
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
	var aiModels []models.AIModel
	if err := config.DB.Where("is_active = ?", true).Order("id ASC").Find(&aiModels).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar model", err.Error())
		return
	}

	respItems := make([]dto.ModelItem, 0, len(aiModels))
	for _, model := range aiModels {
		respItems = append(respItems, dto.ModelItem{
			ID:          model.Slug,
			Type:        "model",
			Object:      "model",
			DisplayName: model.Name,
			Created:     model.CreatedAt.Unix(),
			OwnedBy:     "huggingface",
		})
	}

	c.JSON(http.StatusOK, dto.ModelListResponse{
		Object: "list",
		Data:   respItems,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// 3. GET /v1/models/:model & GET /models/:model
// ─────────────────────────────────────────────────────────────────────────────

func (h *AIHandler) GetModelDetail(c *gin.Context) {
	modelSlug := strings.TrimPrefix(c.Param("model"), "/")
	if modelSlug == "" {
		modelSlug = h.Env.HFModel
	}

	var aiModel models.AIModel
	if err := config.DB.Where("slug = ? AND is_active = ?", modelSlug, true).First(&aiModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Model tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail model", err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.ModelItem{
		ID:          aiModel.Slug,
		Type:        "model",
		Object:      "model",
		DisplayName: aiModel.Name,
		Created:     aiModel.CreatedAt.Unix(),
		OwnedBy:     "huggingface",
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// 4. POST /v1/chat/completions & /chat/completions — OpenAI Standard Chat
// Mendukung streaming (stream: true) untuk aplikasi yang memanggil via format OpenAI
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
	if err := json.Unmarshal(rawData, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Format JSON tidak valid: " + err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// Model default
	model, _ := payload["model"].(string)
	if model == "" {
		model = h.Env.HFModel
		payload["model"] = model
	}

	isStream, _ := payload["stream"].(bool)

	// Pastikan ke HuggingFace selalu non-stream agar dapat respons lengkap dengan stabil
	payload["stream"] = false
	if modifiedData, err := json.Marshal(payload); err == nil {
		rawData = modifiedData
	}

	respBytes, statusCode, err := h.HFService.ForwardChatCompletions(rawData)
	if err != nil || statusCode != http.StatusOK {
		c.Data(statusCode, "application/json", respBytes)
		return
	}

	// Jika non-stream, kembalikan JSON standar OpenAI
	if !isStream {
		c.Data(statusCode, "application/json; charset=utf-8", respBytes)
		return
	}

	// Jika stream: true, kirim SSE chunk format OpenAI
	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBytes, &chatResp); err != nil || len(chatResp.Choices) == 0 {
		c.Data(statusCode, "application/json", respBytes)
		return
	}

	answerText := chatResp.Choices[0].Message.Content

	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	now := time.Now().Unix()
	chatCmplID := fmt.Sprintf("chatcmpl-%d", now)

	// Chunk 1: Role
	roleChunk, _ := json.Marshal(map[string]interface{}{
		"id":      chatCmplID,
		"object":  "chat.completion.chunk",
		"created": now,
		"model":   model,
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"delta": map[string]interface{}{
					"role": "assistant",
				},
				"finish_reason": nil,
			},
		},
	})
	fmt.Fprintf(c.Writer, "data: %s\n\n", roleChunk)
	c.Writer.Flush()

	// Chunk 2+: Content in words
	words := strings.Fields(answerText)
	if len(words) == 0 {
		words = []string{answerText}
	}
	chunkSize := 3
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunkText := strings.Join(words[i:end], " ")
		if end < len(words) {
			chunkText += " "
		}

		contentChunk, _ := json.Marshal(map[string]interface{}{
			"id":      chatCmplID,
			"object":  "chat.completion.chunk",
			"created": now,
			"model":   model,
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"delta": map[string]interface{}{
						"content": chunkText,
					},
					"finish_reason": nil,
				},
			},
		})
		fmt.Fprintf(c.Writer, "data: %s\n\n", contentChunk)
		c.Writer.Flush()
		time.Sleep(5 * time.Millisecond)
	}

	// Final chunk: stop
	stopChunk, _ := json.Marshal(map[string]interface{}{
		"id":      chatCmplID,
		"object":  "chat.completion.chunk",
		"created": now,
		"model":   model,
		"choices": []map[string]interface{}{
			{
				"index":         0,
				"delta":         map[string]interface{}{},
				"finish_reason": "stop",
			},
		},
	})
	fmt.Fprintf(c.Writer, "data: %s\n\n", stopChunk)
	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	c.Writer.Flush()
}
