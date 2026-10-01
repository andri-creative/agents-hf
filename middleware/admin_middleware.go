// [UPGRADE v3] — Menambahkan pembatasan endpoint khusus admin.
package middleware

import (
	"net/http"

	"ai-generate-api/utils"

	"github.com/gin-gonic/gin"
)

// AdminOnly memastikan endpoint hanya bisa diakses oleh user dengan role admin.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.ErrorResponse(c, http.StatusForbidden, "Akses hanya untuk admin", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
