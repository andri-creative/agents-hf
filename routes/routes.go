// routes/routes.go — Pendaftaran semua route/endpoint API
// Menggabungkan handler, middleware, dan konfigurasi routing Gin

package routes

import (
	"net/http"

	"ai-generate-api/config"
	"ai-generate-api/handlers"
	"ai-generate-api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes menginisialisasi router Gin dan mendaftarkan semua endpoint
// Mengembalikan *gin.Engine yang siap digunakan untuk menjalankan server
func SetupRoutes(env *config.EnvConfig) *gin.Engine {
	// Set mode Gin berdasarkan APP_ENV
	if env.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Buat router Gin dengan default middleware (Logger + Recovery)
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	// Inisialisasi handler dengan dependensi yang dibutuhkan
	authHandler := handlers.NewAuthHandler(env.JWTSecret, env.JWTExpiredHours)
	aiHandler := handlers.NewAIHandler(env)

	// ─────────────────────────────────────────────
	// Endpoint publik — tidak membutuhkan autentikasi
	// ─────────────────────────────────────────────

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ─────────────────────────────────────────────
	// Endpoint Models (OpenAI-compatible)
	// Dibutuhkan oleh 9router, One-API, NextChat saat "Fetch Models"
	// Mendukung akses via /v1/models dan /models
	// ─────────────────────────────────────────────
	r.GET("/v1/models", aiHandler.GetModels)
	r.GET("/v1/models/:model", aiHandler.GetModelDetail)
	r.GET("/models", aiHandler.GetModels)
	r.GET("/models/:model", aiHandler.GetModelDetail)
	r.GET("/v1/messages/models", aiHandler.GetModels)

	// Group endpoint /api untuk autentikasi
	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register) // POST /api/register
		api.POST("/login", authHandler.Login)       // POST /api/login
	}

	// ─────────────────────────────────────────────
	// Endpoint protected — membutuhkan JWT token
	// ─────────────────────────────────────────────

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth(env.JWTSecret))
	{
		protected.GET("/me", authHandler.GetProfile) // GET /api/me
	}

	// ─────────────────────────────────────────────
	// Endpoint AI / OpenAI-compatible (Protected API Key)
	// Mendukung:
	// - POST /v1/messages         (Format custom prompt)
	// - POST /v1/chat/completions (Format standar OpenAI/9router)
	// - POST /chat/completions    (Format standar OpenAI/9router tanpa prefix)
	// ─────────────────────────────────────────────

	v1 := r.Group("/v1")
	v1.Use(middleware.APIKeyAuth())
	{
		v1.POST("/messages", aiHandler.Generate)                // POST /v1/messages
		v1.POST("/chat/completions", aiHandler.ChatCompletions) // POST /v1/chat/completions
	}

	// Router root tanpa prefix /v1 untuk client yang baseURL-nya sudah /v1
	rootChat := r.Group("")
	rootChat.Use(middleware.APIKeyAuth())
	{
		rootChat.POST("/chat/completions", aiHandler.ChatCompletions)
	}

	return r
}
