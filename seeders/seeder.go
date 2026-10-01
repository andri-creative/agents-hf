// [UPGRADE v3] — Menambahkan seeder idempotent untuk admin default dan ai_models.
// [UPGRADE 3] — Menambahkan seeder config default (HF_API_KEY, HF_MODEL, HF_BASE_URL).
package seeders

import (
	"log"
	"os"

	"ai-generate-api/config"
	"ai-generate-api/models"
	"ai-generate-api/utils"

	"gorm.io/gorm"
)

// SeedAll menjalankan seluruh seeder yang dibutuhkan aplikasi.
func SeedAll(env *config.EnvConfig) {
	seedDefaultAdmin()
	seedDefaultModel(env)
	seedDefaultConfigs(env)
}

func seedDefaultAdmin() {
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminUsername == "" || adminPassword == "" {
		log.Println("[INFO] Seeder admin dilewati karena ADMIN_USERNAME atau ADMIN_PASSWORD belum di-set")
		return
	}

	var admin models.User
	err := config.DB.Where("username = ?", adminUsername).First(&admin).Error
	if err == nil {
		if admin.Role != "admin" {
			if updateErr := config.DB.Model(&admin).Update("role", "admin").Error; updateErr != nil {
				log.Printf("[WARN] Gagal mengubah role user %s menjadi admin: %v", adminUsername, updateErr)
			}
		}
		return
	}

	if err != gorm.ErrRecordNotFound {
		log.Printf("[WARN] Gagal mengecek admin default: %v", err)
		return
	}

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		log.Printf("[WARN] Gagal hash password admin default: %v", err)
		return
	}

	admin = models.User{
		FullName: "Administrator",
		Username: adminUsername,
		Password: hashedPassword,
		Role:     "admin",
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		log.Printf("[WARN] Gagal membuat admin default: %v", err)
		return
	}

	log.Printf("[INFO] Admin default berhasil dibuat untuk username %s", adminUsername)
}

func seedDefaultModel(env *config.EnvConfig) {
	if env == nil || env.HFModel == "" {
		return
	}

	defaultModel := models.AIModel{
		Name:        "Default Model",
		Slug:        env.HFModel,
		Description: "Model default dari environment aplikasi",
		IsActive:    true,
	}

	if err := config.DB.Where("slug = ?", env.HFModel).FirstOrCreate(&defaultModel).Error; err != nil {
		log.Printf("[WARN] Gagal melakukan seed model default: %v", err)
	}
}

// seedDefaultConfigs menyisipkan config default dari ENV.
// [UPGRADE 3] — Idempotent: hanya insert kalau key belum ada (FirstOrCreate).
func seedDefaultConfigs(env *config.EnvConfig) {
	if env == nil {
		return
	}

	defaults := []models.Config{
		{Key: "HF_API_KEY", Value: env.HFApiKey, Description: "API Key HuggingFace"},
		{Key: "HF_MODEL", Value: env.HFModel, Description: "Model HuggingFace default"},
		{Key: "HF_BASE_URL", Value: env.HFBaseURL, Description: "Base URL HuggingFace Router"},
	}

	for _, cfg := range defaults {
		// Skip config dengan value kosong agar tidak menimpa fallback ENV dengan string kosong.
		if cfg.Value == "" {
			continue
		}

		// FirstOrCreate berdasarkan key → tidak menimpa value yang sudah diubah admin.
		if err := config.DB.Where("key = ?", cfg.Key).FirstOrCreate(&cfg).Error; err != nil {
			log.Printf("[WARN] Gagal melakukan seed config %s: %v", cfg.Key, err)
		}
	}
}
