// dto/token_dto.go — Data Transfer Objects untuk pengelolaan API token
// [UPGRADE 1] — Menambah DTO untuk CRUD token pada tabel api_tokens.

package dto

import "time"

// CreateTokenRequest adalah body request untuk endpoint POST /api/tokens
type CreateTokenRequest struct {
	Name      string     `json:"name" validate:"required,min=1,max=100"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// UpdateTokenStatusRequest adalah body request untuk endpoint PUT /api/tokens/:id/status
type UpdateTokenStatusRequest struct {
	IsActive bool `json:"is_active"`
}
