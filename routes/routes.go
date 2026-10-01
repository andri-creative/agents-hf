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
	// Endpoint AI — membutuhkan API Key (bukan JWT)
	// ─────────────────────────────────────────────

	v1 := r.Group("/v1")
	v1.Use(middleware.APIKeyAuth())
	{
		v1.POST("/messages", aiHandler.Generate) // POST /v1/messages
	}

	return r
}
