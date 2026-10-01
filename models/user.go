// models/user.go — Definisi model User untuk database (tabel: users)
// Memetakan struct Go ke schema database menggunakan GORM tags

package models

import "time"

// User merepresentasikan entitas pengguna dalam sistem
// Password disembunyikan dari JSON response menggunakan json:"-"
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"           json:"id"`
	FullName  string    `gorm:"type:varchar(120);not null"         json:"full_name"`
	Username  string    `gorm:"type:varchar(60);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null"         json:"-"` // TIDAK pernah dikirim ke client
	APIToken  string    `gorm:"type:varchar(130);uniqueIndex"      json:"api_token,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime"                     json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"                     json:"updated_at"`
}
