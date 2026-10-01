// models/api_token.go — Definisi model APIToken (tabel: api_tokens)
// [UPGRADE 1] — Memisahkan API token dari tabel users agar 1 user bisa punya banyak token (multi-device).

package models

import "time"

// APIToken merepresentasikan satu API token milik user.
// Relasi: satu User bisa memiliki banyak APIToken (has many).
type APIToken struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"index;not null" json:"user_id"`
	Token      string     `gorm:"type:varchar(80);uniqueIndex;not null" json:"token"`
	Name       string     `gorm:"type:varchar(100)" json:"name"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
