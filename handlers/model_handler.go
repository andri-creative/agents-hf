// [UPGRADE v3] — Menambahkan CRUD ai_models khusus admin.
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ai-generate-api/config"
	"ai-generate-api/dto"
	"ai-generate-api/models"
	"ai-generate-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// ModelHandler menangani endpoint admin untuk ai_models.
type ModelHandler struct{}

func NewModelHandler() *ModelHandler {
	return &ModelHandler{}
}

// List menangani GET /api/admin/models.
func (h *ModelHandler) List(c *gin.Context) {
	var aiModels []models.AIModel
	if err := config.DB.Order("id ASC").Find(&aiModels).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar model", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar model berhasil diambil", aiModels)
}

// Create menangani POST /api/admin/models.
func (h *ModelHandler) Create(c *gin.Context) {
	// Gunakan raw body agar kita bisa mendeteksi apakah is_active dikirim atau tidak.
	rawData, err := c.GetRawData()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca request body", err.Error())
		return
	}

	var req dto.CreateModelRequest
	if err := json.Unmarshal(rawData, &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", err.Error())
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorDetails := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorDetails[fieldErr.Field()] = fieldErr.Tag()
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi input gagal", errorDetails)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Slug = strings.TrimSpace(req.Slug)
	req.Description = strings.TrimSpace(req.Description)

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(rawData, &payload); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", err.Error())
		return
	}

	isActive := true
	if _, exists := payload["is_active"]; exists {
		isActive = req.IsActive
	}

	var existing models.AIModel
	if err := config.DB.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Slug model sudah digunakan", "Gunakan slug lain")
		return
	}

	aiModel := models.AIModel{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IsActive:    isActive,
	}

	if err := config.DB.Create(&aiModel).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan model", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Model berhasil dibuat", aiModel)
}

// Detail menangani GET /api/admin/models/:id.
func (h *ModelHandler) Detail(c *gin.Context) {
	var aiModel models.AIModel
	if err := config.DB.First(&aiModel, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Model tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail model", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail model berhasil diambil", aiModel)
}

// Update menangani PUT /api/admin/models/:id.
func (h *ModelHandler) Update(c *gin.Context) {
	var req dto.UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", err.Error())
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorDetails := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorDetails[fieldErr.Field()] = fieldErr.Tag()
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi input gagal", errorDetails)
		return
	}

	var aiModel models.AIModel
	if err := config.DB.First(&aiModel, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Model tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil model", err.Error())
		return
	}

	if req.Name != "" {
		aiModel.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		aiModel.Description = strings.TrimSpace(req.Description)
	}
	if req.Slug != "" {
		newSlug := strings.TrimSpace(req.Slug)
		if newSlug != aiModel.Slug {
			var existing models.AIModel
			if err := config.DB.Where("slug = ?", newSlug).First(&existing).Error; err == nil {
				utils.ErrorResponse(c, http.StatusConflict, "Slug model sudah digunakan", "Gunakan slug lain")
				return
			}
		}
		aiModel.Slug = newSlug
	}
	if req.IsActive != nil {
		aiModel.IsActive = *req.IsActive
	}

	aiModel.UpdatedAt = time.Now()
	if err := config.DB.Save(&aiModel).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui model", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Model berhasil diperbarui", aiModel)
}

// Delete menangani DELETE /api/admin/models/:id.
func (h *ModelHandler) Delete(c *gin.Context) {
	var aiModel models.AIModel
	if err := config.DB.First(&aiModel, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Model tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil model", err.Error())
		return
	}

	if err := config.DB.Delete(&aiModel).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus model", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Model berhasil dihapus", gin.H{"id": aiModel.ID})
}
