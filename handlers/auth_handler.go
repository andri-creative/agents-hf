// handlers/auth_handler.go — Handler HTTP untuk endpoint autentikasi
// Menangani: Register (/api/register), Login (/api/login), dan Profile (/api/me)
// Logic berat didelegasikan ke layer service/utils, handler hanya urus HTTP

package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"ai-generate-api/config"
	"ai-generate-api/dto"
	"ai-generate-api/models"
	"ai-generate-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// validate adalah instance validator yang digunakan di seluruh handler auth
var validate = validator.New()

// AuthHandler menyimpan dependensi yang dibutuhkan oleh handler auth
type AuthHandler struct {
	JWTSecret       string
	JWTExpiredHours int
}

// NewAuthHandler membuat instance AuthHandler baru
func NewAuthHandler(jwtSecret string, jwtExpiredHours int) *AuthHandler {
	return &AuthHandler{
		JWTSecret:       jwtSecret,
		JWTExpiredHours: jwtExpiredHours,
	}
}

// Register menangani POST /api/register
// Membuat akun baru dengan hash password dan generate API token unik
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", err.Error())
		return
	}

	// Validasi input menggunakan go-playground/validator
	if err := validate.Struct(req); err != nil {
		// Ekstrak detail error validasi per field
		validationErrors := err.(validator.ValidationErrors)
		errorDetails := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorDetails[fieldErr.Field()] = fieldErr.Tag()
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi input gagal", errorDetails)
		return
	}

	// Cek apakah username sudah digunakan
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Username sudah digunakan", "Pilih username lain")
		return
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses password", err.Error())
		return
	}

	// Generate API Token unik: "hf_" + 64 karakter hex random
	// Menggunakan crypto/rand untuk keamanan kriptografis
	apiToken, err := generateAPIToken()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal generate API token", err.Error())
		return
	}

	// Buat record user baru
	user := models.User{
		FullName: req.FullName,
		Username: req.Username,
		Password: hashedPassword,
		APIToken: apiToken,
	}

	// Simpan ke database
	if err := config.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan user", err.Error())
		return
	}

	// Susun response sukses (tanpa password)
	userResp := dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	utils.SuccessResponse(c, http.StatusCreated, "Registrasi berhasil", dto.RegisterResponse{
		APIToken: apiToken,
		User:     userResp,
	})
}

// Login menangani POST /api/login
// Memverifikasi kredensial dan mengembalikan JWT token + API token
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", err.Error())
		return
	}

	// Validasi input
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorDetails := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorDetails[fieldErr.Field()] = fieldErr.Tag()
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi input gagal", errorDetails)
		return
	}

	// Cari user berdasarkan username
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		// Jangan berikan info apakah username atau password yang salah (keamanan)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Username atau password salah", nil)
		return
	}

	// Verifikasi password dengan bcrypt
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Username atau password salah", nil)
		return
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateJWT(user.ID, user.Username, h.JWTSecret, h.JWTExpiredHours)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal generate token", err.Error())
		return
	}

	// Susun response
	userResp := dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", dto.LoginResponse{
		JWTToken: jwtToken,
		APIToken: user.APIToken,
		User:     userResp,
	})
}

// GetProfile menangani GET /api/me (Protected — JWT required)
// Mengembalikan data user yang sedang login berdasarkan user_id di JWT claims
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Ambil user_id dari context (di-set oleh JWTAuth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User tidak terautentikasi", nil)
		return
	}

	// Cari user di database
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User tidak ditemukan", nil)
		return
	}

	// Return data user (password otomatis tersembunyi karena json:"-")
	userResp := dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	utils.SuccessResponse(c, http.StatusOK, "Data profile berhasil diambil", userResp)
}

// generateAPIToken membuat API token unik dengan format: hf_ + 64 karakter hex random
// Menggunakan crypto/rand untuk keamanan kriptografis (BUKAN math/rand)
func generateAPIToken() (string, error) {
	// 32 bytes random → 64 karakter hex
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "hf_" + hex.EncodeToString(bytes), nil
}
