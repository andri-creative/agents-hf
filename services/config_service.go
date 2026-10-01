// services/config_service.go — Service untuk membaca konfigurasi aplikasi
// [UPGRADE 3] — Ambil config dari DB, fallback ke ENV bila belum ada di database.

package services

import (
	"ai-generate-api/config"
	"ai-generate-api/models"
)

// GetConfig mengambil nilai config dari database berdasarkan key.
// Jika tidak ditemukan / kosong, fallback ke nilai ENV yang sesuai.
func GetConfig(key string) string {
	var cfg models.Config
	if err := config.DB.Where("key = ?", key).First(&cfg).Error; err == nil && cfg.Value != "" {
		return cfg.Value
	}

	// Fallback ke ENV
	if config.Env == nil {
		return ""
	}

	switch key {
	case "HF_API_KEY":
		return config.Env.HFApiKey
	case "HF_MODEL":
		return config.Env.HFModel
	case "HF_BASE_URL":
		return config.Env.HFBaseURL
	}

	return ""
}

// GetHuggingFaceConfig mengambil 3 konfigurasi HuggingFace sekaligus.
func GetHuggingFaceConfig() (apiKey, model, baseURL string) {
	return GetConfig("HF_API_KEY"),
		GetConfig("HF_MODEL"),
		GetConfig("HF_BASE_URL")
}
