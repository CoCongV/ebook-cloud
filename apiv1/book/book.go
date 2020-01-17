package book

import (
	"github.com/gin-gonic/gin"

	"EbookCloud/app"
	"EbookCloud/app/models"
)

// GetBooks get all books
func GetBooks(c *gin.Context) {
	var books []models.Book
	app.db.Model(&books).Select()
	c.JSON(200, gin.H{
		"message": "books",
	})
}
