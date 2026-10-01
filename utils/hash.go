// utils/hash.go — Utility untuk hashing dan verifikasi password menggunakan bcrypt
// Menggunakan golang.org/x/crypto/bcrypt untuk keamanan maksimal

package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword menghasilkan hash bcrypt dari password plaintext
// Cost default (bcrypt.DefaultCost = 10) digunakan untuk keseimbangan keamanan & performa
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword memverifikasi apakah password plaintext cocok dengan hash yang tersimpan
// Mengembalikan true jika cocok, false jika tidak
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
