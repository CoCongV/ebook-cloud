package book

import "github.com/gin-gonic/gin"

func getBooks(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "books",
	})
}
