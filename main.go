// main.go — Entry point aplikasi REST API AI Generate
// Bertanggung jawab untuk inisialisasi config, database, dan menjalankan server HTTP
// [UPGRADE v3] — Menambahkan migrasi ai_models dan seeding data awal admin/model.
// [UPGRADE 1 & 3] — Menambahkan migrasi api_tokens & configs serta migrasi data token lama.

package main

import (
	"log"

	"ai-generate-api/config"
	"ai-generate-api/models"
	"ai-generate-api/routes"
	"ai-generate-api/seeders"

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

	// Auto migrate semua model aplikasi ke database
	// [UPGRADE 1] — Menambahkan &models.APIToken{} (token dipisah dari tabel users)
	// [UPGRADE 3] — Menambahkan &models.Config{}
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.APIToken{},
		&models.AIModel{},
		&models.Config{},
	); err != nil {
		log.Fatalf("[FATAL] Gagal auto migrate: %v", err)
	}
	log.Println("[INFO] Auto migrate berhasil")

	// [UPGRADE 1] — Migrasi data token lama dari users.api_token ke tabel api_tokens.
	// Dijalankan setelah AutoMigrate (tabel api_tokens sudah ada) dan bersifat idempotent.
	seeders.MigrateLegacyAPITokens()

	// Seed data awal secara idempotent.
	seeders.SeedAll(env)

	// Setup router Gin dan daftarkan semua routes
	r := routes.SetupRoutes(env)

	// Jalankan server HTTP
	addr := ":" + env.AppPort
	log.Printf("[INFO] Server berjalan di http://localhost%s (ENV: %s)", addr, env.AppEnv)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[FATAL] Server gagal dijalankan: %v", err)
	}
}
