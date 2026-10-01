// handlers/token_handler.go — Handler HTTP untuk pengelolaan API token user
// [UPGRADE 1] — CRUD token pada tabel api_tokens. User hanya bisa mengelola token miliknya sendiri.

package handlers

import (
	"net/http"

	"ai-generate-api/config"
	"ai-generate-api/dto"
	"ai-generate-api/models"
	"ai-generate-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// TokenHandler menangani endpoint pengelolaan API token.
type TokenHandler struct{}

// NewTokenHandler membuat instance TokenHandler baru.
func NewTokenHandler() *TokenHandler {
	return &TokenHandler{}
}

// currentUserID mengambil user_id dari context JWT.
func currentUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}

// List menangani GET /api/tokens — daftar semua token milik user yang login.
func (h *TokenHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User tidak terautentikasi", nil)
		return
	}

	var tokens []models.APIToken
	if err := config.DB.Where("user_id = ?", userID).Order("id DESC").Find(&tokens).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar token berhasil diambil", tokens)
}

// Create menangani POST /api/tokens — generate token baru untuk user yang login.
func (h *TokenHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User tidak terautentikasi", nil)
		return
	}

	var req dto.CreateTokenRequest
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

	tokenString, err := generateAPIToken()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal generate API token", err.Error())
		return
	}

	apiToken := models.APIToken{
		UserID:    userID,
		Token:     tokenString,
		Name:      req.Name,
		IsActive:  true,
		ExpiresAt: req.ExpiresAt,
	}

	if err := config.DB.Create(&apiToken).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Token berhasil dibuat", apiToken)
}

// Detail menangani GET /api/tokens/:id — detail token milik user yang login.
func (h *TokenHandler) Detail(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User tidak terautentikasi", nil)
		return
	}

	var apiToken models.APIToken
	if err := config.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&apiToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Token tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail token berhasil diambil", apiToken)
}

// UpdateStatus menangani PUT /api/tokens/:id/status — aktif/nonaktifkan token.
func (h *TokenHandler) UpdateStatus(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User tidak terautentikasi", nil)
		return
	}

	var req dto.UpdateTokenStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", err.Error())
		return
	}

	var apiToken models.APIToken
	if err := config.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&apiToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Token tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil token", err.Error())
		return
	}

	apiToken.IsActive = req.IsActive
	if err := config.DB.Save(&apiToken).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui status token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Status token berhasil diperbarui", apiToken)
}

// Delete menangani DELETE /api/tokens/:id — hapus token (hard delete).
func (h *TokenHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User tidak terautentikasi", nil)
		return
	}

	var apiToken models.APIToken
	if err := config.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&apiToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Token tidak ditemukan", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil token", err.Error())
		return
	}

	if err := config.DB.Delete(&apiToken).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token berhasil dihapus", gin.H{"id": apiToken.ID})
}
