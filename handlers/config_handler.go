// handlers/config_handler.go — Handler HTTP untuk pengelolaan config aplikasi (khusus admin)
// [UPGRADE 3] — CRUD config pada tabel configs. Endpoint diproteksi JWT + AdminOnly.

package handlers

import (
	"net/http"
	"strings"

	"ai-generate-api/config"
	"ai-generate-api/dto"
	"ai-generate-api/models"
	"ai-generate-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// ConfigHandler menangani endpoint admin untuk tabel configs.
type ConfigHandler struct{}

// NewConfigHandler membuat instance ConfigHandler baru.
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

// List menangani GET /api/admin/configs.
func (h *ConfigHandler) List(c *gin.Context) {
	var configs []models.Config
	if err := config.DB.Order("id ASC").Find(&configs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar config", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar config berhasil diambil", configs)
}

// Create menangani POST /api/admin/configs.
func (h *ConfigHandler) Create(c *gin.Context) {
	var req dto.CreateConfigRequest
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

	req.Key = strings.TrimSpace(req.Key)
	req.Description = strings.TrimSpace(req.Description)

	var existing models.Config
	if err := config.DB.Where("key = ?", req.Key).First(&existing).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Key config sudah ada", "Gunakan key lain atau update config yang sudah ada")
		return
	}

	cfg := models.Config{
		Key:         req.Key,
		Value:       req.Value,
		Description: req.Description,
	}

	if err := config.DB.Create(&cfg).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan config", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Config berhasil dibuat", cfg)
}

// Detail menangani GET /api/admin/configs/:id.
func (h *ConfigHandler) Detail(c *gin.Context) {
	var cfg models.Config
	if err := config.DB.First(&cfg, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Config tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail config", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail config berhasil diambil", cfg)
}

// Update menangani PUT /api/admin/configs/:id.
func (h *ConfigHandler) Update(c *gin.Context) {
	var req dto.UpdateConfigRequest
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

	var cfg models.Config
	if err := config.DB.First(&cfg, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Config tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil config", err.Error())
		return
	}

	if req.Value != "" {
		cfg.Value = req.Value
	}
	if req.Description != "" {
		cfg.Description = strings.TrimSpace(req.Description)
	}

	if err := config.DB.Save(&cfg).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui config", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Config berhasil diperbarui", cfg)
}

// Delete menangani DELETE /api/admin/configs/:id.
func (h *ConfigHandler) Delete(c *gin.Context) {
	var cfg models.Config
	if err := config.DB.First(&cfg, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Config tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil config", err.Error())
		return
	}

	if err := config.DB.Delete(&cfg).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus config", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Config berhasil dihapus", gin.H{"id": cfg.ID})
}
