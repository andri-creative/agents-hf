// api/index.go — Vercel Serverless Function entrypoint
// Vercel menjalankan fungsi ini sebagai serverless handler (package handler, func Handler)

package handler

import (
	"log"
	"net/http"
	"sync"

	"ai-generate-api/config"
	"ai-generate-api/models"
	"ai-generate-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	engine *gin.Engine
	once   sync.Once
)

// initApp menginisialisasi konfigurasi, database, dan routes sekali saja (reused across warm serverless calls)
func initApp() {
	// Di environment lokal/vercel dev, baca file .env jika ada
	_ = godotenv.Load()

	// Load konfigurasi dari environment variables Vercel
	env := config.LoadEnv()

	// Inisialisasi koneksi database PostgreSQL
	config.InitDB(env)

	// Auto migrate tabel User
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Printf("[WARN] Auto migrate error: %v", err)
	}

	// Inisialisasi Gin router
	engine = routes.SetupRoutes(env)
}

// Handler adalah fungsi utama yang dipanggil oleh Vercel untuk setiap request yang masuk
func Handler(w http.ResponseWriter, r *http.Request) {
	// Pastikan inisialisasi hanya dijalankan sekali selama lifecycle instance serverless (warm container reuse)
	once.Do(initApp)

	// Serahkan penanganan request ke Gin engine
	engine.ServeHTTP(w, r)
}
