// config/env.go — Membaca dan memvalidasi semua environment variables
// Mendukung DATABASE_URL (connection string lengkap) atau parameter individual DB

package config

import (
	"log"
	"os"
	"strconv"
)

// EnvConfig menyimpan semua konfigurasi yang dibaca dari environment
type EnvConfig struct {
	AppPort         string
	AppEnv          string
	DatabaseURL     string // Prioritas utama — connection string lengkap (misal dari Prisma/Railway/Supabase)
	// Parameter individual (digunakan jika DATABASE_URL tidak ada)
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	JWTSecret       string
	JWTExpiredHours int
	HFApiKey        string
	HFModel         string
	HFBaseURL       string
}

// LoadEnv membaca environment variables dan mengembalikan EnvConfig
// Jika JWT_SECRET kosong, server akan gagal start (log.Fatal)
func LoadEnv() *EnvConfig {
	// JWT_SECRET wajib ada — fatal jika kosong
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("[FATAL] JWT_SECRET tidak boleh kosong. Set di file .env atau environment system.")
	}

	// JWT_EXPIRED_HOURS, default 24 jam
	jwtExpiredHours := 24
	if val := os.Getenv("JWT_EXPIRED_HOURS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			jwtExpiredHours = parsed
		}
	}

	// Ambil DATABASE_URL — bisa dari DATABASE_URL, POSTGRES_URL, atau PRISMA_DATABASE_URL
	// Prioritas: DATABASE_URL > POSTGRES_URL > PRISMA_DATABASE_URL > parameter individual
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("POSTGRES_URL")
	}
	if databaseURL == "" {
		databaseURL = os.Getenv("PRISMA_DATABASE_URL")
	}

	return &EnvConfig{
		AppPort:         getEnvOrDefault("APP_PORT", "8080"),
		AppEnv:          getEnvOrDefault("APP_ENV", "development"),
		DatabaseURL:     databaseURL,
		// Parameter individual sebagai fallback
		DBHost:          getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:          getEnvOrDefault("DB_PORT", "5432"),
		DBUser:          getEnvOrDefault("DB_USER", "postgres"),
		DBPassword:      getEnvOrDefault("DB_PASSWORD", ""),
		DBName:          getEnvOrDefault("DB_NAME", "ai_generate_db"),
		DBSSLMode:       getEnvOrDefault("DB_SSLMODE", "disable"),
		JWTSecret:       jwtSecret,
		JWTExpiredHours: jwtExpiredHours,
		HFApiKey:        getEnvOrDefault("HF_API_KEY", ""),
		HFModel:         getEnvOrDefault("HF_MODEL", "mistralai/Mistral-7B-Instruct-v0.2"),
		HFBaseURL:       getEnvOrDefault("HF_BASE_URL", "https://api-inference.huggingface.co"),
	}
}

// getEnvOrDefault mengambil nilai environment variable atau mengembalikan nilai default
func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
