package apiv1

import (
	"github.com/gin-gonic/gin"
)

// GetBooks get all books
func GetBooks(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "books",
	})
}
