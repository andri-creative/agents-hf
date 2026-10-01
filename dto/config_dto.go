// dto/config_dto.go — Data Transfer Objects untuk pengelolaan config aplikasi
// [UPGRADE 3] — DTO CRUD config pada tabel configs.

package dto

// CreateConfigRequest adalah body request untuk endpoint POST /api/admin/configs
type CreateConfigRequest struct {
	Key         string `json:"key" validate:"required,min=2,max=100"`
	Value       string `json:"value" validate:"required"`
	Description string `json:"description" validate:"max=255"`
}

// UpdateConfigRequest adalah body request untuk endpoint PUT /api/admin/configs/:id
type UpdateConfigRequest struct {
	Value       string `json:"value" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty,max=255"`
}
