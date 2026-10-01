// utils/response.go — Helper fungsi untuk format response HTTP yang konsisten
// Memastikan semua response menggunakan format: { success, message, data/error }

package utils

import "github.com/gin-gonic/gin"

// SuccessResponse mengirimkan response sukses dengan format standar
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse mengirimkan response error dengan format standar
func ErrorResponse(c *gin.Context, statusCode int, message string, errDetail interface{}) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"error":   errDetail,
	})
}
