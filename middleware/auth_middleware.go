// middleware/auth_middleware.go — Middleware autentikasi untuk melindungi endpoint
// Berisi 2 middleware: JWTAuth (untuk /api/me) dan APIKeyAuth (untuk /v1/messages)
// [UPGRADE v3] — Menyimpan role dari JWT ke context untuk endpoint admin.
// [UPGRADE 1] — APIKeyAuth mencari token di tabel api_tokens (bukan lagi users.api_token).

package middleware

import (
	"strings"
	"time"

	"ai-generate-api/config"
	"ai-generate-api/models"
	"ai-generate-api/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth adalah middleware yang memvalidasi JWT token dari header Authorization
// Digunakan untuk endpoint yang membutuhkan autentikasi JWT (misal: GET /api/me)
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, 401, "Token tidak ditemukan", "Header Authorization diperlukan")
			c.Abort()
			return
		}

		// Format header harus: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, 401, "Format token tidak valid", "Gunakan format: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi dan parse JWT token
		claims, err := utils.ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			utils.ErrorResponse(c, 401, "Token tidak valid atau sudah kadaluarsa", err.Error())
			c.Abort()
			return
		}

		// Simpan data user ke context untuk digunakan di handler
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// APIKeyAuth adalah middleware yang memvalidasi API Token dari header X-API-Key atau Authorization
// Digunakan khusus untuk endpoint /v1/messages
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiToken string

		// Cek header X-API-Key terlebih dahulu
		apiToken = c.GetHeader("X-API-Key")

		// Jika X-API-Key kosong, coba ambil dari header Authorization: Bearer <token>
		if apiToken == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					apiToken = parts[1]
				}
			}
		}

		// Jika tetap kosong, tolak request
		if apiToken == "" {
			utils.ErrorResponse(c, 401, "API Key tidak ditemukan", "Sertakan header X-API-Key atau Authorization: Bearer <api_token>")
			c.Abort()
			return
		}

		// [UPGRADE 1] — Cari token di tabel api_tokens dengan kondisi:
		// token cocok, masih aktif, dan belum kadaluarsa (expires_at NULL dianggap tidak kadaluarsa).
		var apiTokenRecord models.APIToken
		if err := config.DB.
			Where("token = ? AND is_active = ?", apiToken, true).
			Where("expires_at IS NULL OR expires_at > ?", time.Now()).
			First(&apiTokenRecord).Error; err != nil {
			utils.ErrorResponse(c, 401, "API Key tidak valid", "API Key tidak ditemukan atau sudah tidak aktif")
			c.Abort()
			return
		}

		// Cari user pemilik token untuk melengkapi data di context
		var user models.User
		if err := config.DB.First(&user, apiTokenRecord.UserID).Error; err != nil {
			utils.ErrorResponse(c, 401, "API Key tidak valid", "User pemilik API Key tidak ditemukan")
			c.Abort()
			return
		}

		// [UPGRADE 1] — Catat waktu pemakaian token terakhir.
		now := time.Now()
		if err := config.DB.Model(&apiTokenRecord).Update("last_used_at", now).Error; err != nil {
			utils.ErrorResponse(c, 401, "API Key tidak valid", "Gagal memperbarui status pemakaian token")
			c.Abort()
			return
		}

		// Simpan data user ke context untuk digunakan di handler
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)

		c.Next()
	}
}
