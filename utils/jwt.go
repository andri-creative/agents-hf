// utils/jwt.go — Utility untuk generate dan validasi JWT token
// Menggunakan github.com/golang-jwt/jwt/v5 dengan algoritma HS256
// [UPGRADE v3] — Menambahkan role ke dalam JWT claims untuk otorisasi admin.

package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims mendefinisikan payload yang disimpan di dalam JWT token
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT membuat JWT token baru untuk user yang berhasil login
// secret adalah JWT_SECRET dari ENV, expiredHours adalah durasi token berlaku
func GenerateJWT(userID uint, username, role, secret string, expiredHours int) (string, error) {
	// Tentukan waktu kadaluarsa token
	expirationTime := time.Now().Add(time.Duration(expiredHours) * time.Hour)

	// Buat claims (payload) JWT
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-generate-api",
		},
	}

	// Buat token dengan algoritma HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT memvalidasi JWT token dan mengembalikan claims jika valid
// Memastikan method signing adalah HS256 untuk mencegah algoritma confusion attack
func ValidateJWT(tokenString, secret string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validasi bahwa method signing adalah HS256 (mencegah algoritma confusion attack)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return claims, nil
}
