// routes/routes.go — Pendaftaran semua route/endpoint API
// Menggabungkan handler, middleware, dan konfigurasi routing Gin
// [UPGRADE v3] — Menambahkan route admin untuk CRUD ai_models.
// [UPGRADE 1 & 3] — Menambahkan route /api/tokens dan /api/admin/configs.

package routes

import (
	"net/http"
	"time"

	"ai-generate-api/config"
	"ai-generate-api/handlers"
	"ai-generate-api/middleware"

	"github.com/gin-contrib/cors"
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

	// Konfigurasi CORS untuk frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Frontend Vite dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Inisialisasi handler dengan dependensi yang dibutuhkan
	authHandler := handlers.NewAuthHandler(env.JWTSecret, env.JWTExpiredHours)
	aiHandler := handlers.NewAIHandler(env)
	modelHandler := handlers.NewModelHandler()
	tokenHandler := handlers.NewTokenHandler()
	configHandler := handlers.NewConfigHandler()

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
	r.GET("/v1/models/*model", aiHandler.GetModelDetail)
	r.GET("/models", aiHandler.GetModels)
	r.GET("/models/*model", aiHandler.GetModelDetail)
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

		// [UPGRADE 1] — Endpoint pengelolaan API token milik user yang login
		protected.GET("/tokens", tokenHandler.List)
		protected.POST("/tokens", tokenHandler.Create)
		protected.GET("/tokens/:id", tokenHandler.Detail)
		protected.PUT("/tokens/:id/status", tokenHandler.UpdateStatus)
		protected.DELETE("/tokens/:id", tokenHandler.Delete)
	}

	// Endpoint admin — membutuhkan JWT token + role admin
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth(env.JWTSecret), middleware.AdminOnly())
	{
		admin.GET("/models", modelHandler.List)
		admin.POST("/models", modelHandler.Create)
		admin.GET("/models/:id", modelHandler.Detail)
		admin.PUT("/models/:id", modelHandler.Update)
		admin.DELETE("/models/:id", modelHandler.Delete)

		// [UPGRADE 3] — Endpoint pengelolaan config HF dari admin panel
		admin.GET("/configs", configHandler.List)
		admin.POST("/configs", configHandler.Create)
		admin.GET("/configs/:id", configHandler.Detail)
		admin.PUT("/configs/:id", configHandler.Update)
		admin.DELETE("/configs/:id", configHandler.Delete)
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
