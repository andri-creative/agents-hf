// models/user.go — Definisi model User untuk database (tabel: users)
// Memetakan struct Go ke schema database menggunakan GORM tags
// [UPGRADE v3] — Menambahkan role untuk otorisasi admin.
// [UPGRADE 1] — Menghapus kolom APIToken (dipindah ke tabel api_tokens) dan menambahkan relasi has-many.

package models

import "time"

// User merepresentasikan entitas pengguna dalam sistem
// Password disembunyikan dari JSON response menggunakan json:"-"
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName  string    `gorm:"type:varchar(120);not null" json:"full_name"`
	Username  string    `gorm:"type:varchar(60);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // TIDAK pernah dikirim ke client
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// APITokens adalah relasi has-many ke tabel api_tokens (1 user → banyak token)
	APITokens []APIToken `gorm:"foreignKey:UserID" json:"api_tokens,omitempty"`
}
