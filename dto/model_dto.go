// [UPGRADE v3] — Menambahkan DTO untuk CRUD ai_models oleh admin.
package dto

type CreateModelRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=120"`
	Slug        string `json:"slug" validate:"required,min=2,max=200"`
	Description string `json:"description" validate:"max=1000"`
	IsActive    bool   `json:"is_active"`
}

type UpdateModelRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=120"`
	Slug        string `json:"slug" validate:"omitempty,min=2,max=200"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	IsActive    *bool  `json:"is_active"`
}
