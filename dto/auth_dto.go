// dto/auth_dto.go — Data Transfer Objects untuk endpoint autentikasi (register & login)
// Mendefinisikan struktur request/response dengan validasi menggunakan go-playground/validator
// [UPGRADE v3] — Menambahkan role ke response user agar status admin terlihat.

package dto

// RegisterRequest adalah body request untuk endpoint POST /api/register
type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=120"`
	Username string `json:"username"  validate:"required,alphanum,min=3,max=60"`
	Password string `json:"password"  validate:"required,min=6,max=255"`
}

// LoginRequest adalah body request untuk endpoint POST /api/login
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=60"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

// UserResponse adalah data user yang aman dikirimkan ke client (tanpa password)
type UserResponse struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	Username  string `json:"username"`
	Role      string `json:"role,omitempty"`
	APIToken  string `json:"api_token,omitempty"`
	CreatedAt string `json:"created_at"`
}

// RegisterResponse adalah response untuk endpoint register
type RegisterResponse struct {
	APIToken string       `json:"api_token"`
	User     UserResponse `json:"user"`
}

// LoginResponse adalah response untuk endpoint login
type LoginResponse struct {
	JWTToken string       `json:"jwt_token"`
	APIToken string       `json:"api_token"`
	User     UserResponse `json:"user"`
}
