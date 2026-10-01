// models/config.go — Definisi model Config (tabel: configs)
// [UPGRADE 3] — Menyimpan konfigurasi aplikasi (mis. HF_API_KEY/HF_MODEL/HF_BASE_URL) di database.

package models

import "time"

// Config merepresentasikan satu pasangan key-value konfigurasi aplikasi.
type Config struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value       string    `gorm:"type:text;not null" json:"value"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
