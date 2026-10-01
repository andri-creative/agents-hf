// [UPGRADE v3] — Menambahkan model AI yang dikelola dari database.
package models

import "time"

// AIModel menyimpan daftar model HuggingFace yang boleh digunakan aplikasi.
type AIModel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(120);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
