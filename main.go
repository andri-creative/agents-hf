// main.go — Entry point aplikasi REST API AI Generate
// Bertanggung jawab untuk inisialisasi config, database, dan menjalankan server HTTP

package main

import (
	"log"

	"ai-generate-api/config"
	"ai-generate-api/models"
	"ai-generate-api/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load file .env ke dalam environment variables
	// Jika file .env tidak ada, lanjutkan (ENV mungkin sudah di-set di sistem)
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] File .env tidak ditemukan, menggunakan environment system")
	}

	// Muat dan validasi konfigurasi environment
	env := config.LoadEnv()

	// Inisialisasi koneksi database PostgreSQL
	config.InitDB(env)

	// Auto migrate model User ke database
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("[FATAL] Gagal auto migrate: %v", err)
	}
	log.Println("[INFO] Auto migrate berhasil")

	// Setup router Gin dan daftarkan semua routes
	r := routes.SetupRoutes(env)

	// Jalankan server HTTP
	addr := ":" + env.AppPort
	log.Printf("[INFO] Server berjalan di http://localhost%s (ENV: %s)", addr, env.AppEnv)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[FATAL] Server gagal dijalankan: %v", err)
	}
}
