// seeders/migration.go — Migrasi data aman untuk upgrade skema
// [UPGRADE 1] — Memindahkan data users.api_token yang lama ke tabel api_tokens
// sebelum kolom tersebut benar-benar ditinggalkan (name: "Migrated Token").

package seeders

import (
	"log"

	"ai-generate-api/config"
	"ai-generate-api/models"
)

// MigrateLegacyAPITokens menyalin token lama dari kolom users.api_token ke tabel api_tokens.
// Idempotent: token yang sudah dipindahkan tidak akan diduplikasi.
func MigrateLegacyAPITokens() {
	// Cek apakah kolom users.api_token masih ada.
	// GORM AutoMigrate tidak menghapus kolom, jadi kolom lama biasanya masih tersedia.
	var columnCount int64
	err := config.DB.Raw(`
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name = 'users' AND column_name = 'api_token'
	`).Scan(&columnCount).Error
	if err != nil {
		log.Printf("[WARN] Gagal mengecek kolom users.api_token: %v", err)
		return
	}
	if columnCount == 0 {
		// Kolom sudah tidak ada, tidak ada yang perlu dimigrasikan.
		return
	}

	// Ambil data token lama yang masih terisi.
	type legacyTokenRow struct {
		ID       uint
		APIToken string
	}

	var legacyRows []legacyTokenRow
	if err := config.DB.Raw(`
		SELECT id, api_token
		FROM users
		WHERE api_token IS NOT NULL AND api_token <> ''
	`).Scan(&legacyRows).Error; err != nil {
		log.Printf("[WARN] Gagal membaca token lama dari users: %v", err)
		return
	}

	migrated := 0
	for _, row := range legacyRows {
		// Lewati jika token sudah ada di tabel api_tokens (idempotent).
		var existing models.APIToken
		if err := config.DB.Where("token = ?", row.APIToken).First(&existing).Error; err == nil {
			continue
		}

		migratedToken := models.APIToken{
			UserID:   row.ID,
			Token:    row.APIToken,
			Name:     "Migrated Token",
			IsActive: true,
		}

		if err := config.DB.Create(&migratedToken).Error; err != nil {
			log.Printf("[WARN] Gagal memigrasikan token user %d: %v", row.ID, err)
			continue
		}
		migrated++
	}

	if migrated > 0 {
		log.Printf("[INFO] Berhasil memigrasikan %d token lama ke tabel api_tokens", migrated)
	}
}
