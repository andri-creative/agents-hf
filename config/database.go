// config/database.go — Inisialisasi koneksi database PostgreSQL menggunakan GORM
// Mendukung DATABASE_URL (connection string) ATAU parameter individual DB_HOST, DB_USER, dst.

package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah instance GORM global yang digunakan di seluruh aplikasi
var DB *gorm.DB

// InitDB menginisialisasi koneksi ke database PostgreSQL
// Prioritas: DATABASE_URL (connection string lengkap) > parameter individual
// Akan log.Fatal jika koneksi gagal
func InitDB(env *EnvConfig) {
	var dsn string

	if env.DatabaseURL != "" {
		// Gunakan DATABASE_URL langsung sebagai DSN (dari Prisma, Railway, Supabase, dll.)
		dsn = env.DatabaseURL
		log.Println("[INFO] Menggunakan DATABASE_URL untuk koneksi database")
	} else {
		// Fallback: bangun DSN dari parameter individual
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			env.DBHost,
			env.DBUser,
			env.DBPassword,
			env.DBName,
			env.DBPort,
			env.DBSSLMode,
		)
		log.Printf("[INFO] Menggunakan parameter individual untuk koneksi: %s@%s/%s", env.DBUser, env.DBHost, env.DBName)
	}

	// Konfigurasi level log GORM berdasarkan environment
	var gormLogger logger.Interface
	if env.AppEnv == "development" {
		// Mode development: tampilkan semua query SQL
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		// Mode production: hanya tampilkan error
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Buka koneksi ke PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("[FATAL] Gagal koneksi ke database: %v", err)
	}

	// Konfigurasi connection pool untuk performa optimal
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[FATAL] Gagal mendapatkan instance sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)  // Maksimum 10 koneksi aktif
	sqlDB.SetMaxIdleConns(5)   // Maksimum 5 koneksi idle
	sqlDB.SetConnMaxLifetime(0) // Tanpa batas waktu koneksi

	// Simpan instance ke variabel global
	DB = db
	log.Println("[INFO] Berhasil terhubung ke database PostgreSQL")
}
