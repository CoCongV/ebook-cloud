package apiv1

import (
	"ebook-cloud/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBooks get all books
func GetBooks(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithError(401, err)
		return
	}
	offsetCount := (page - 1) * 20
	itemCount := 20

	var books []models.Book
	models.DB.Offset(offsetCount).Limit(itemCount).Find(&books)
	c.JSON(200, gin.H{
		"books": books,
	})
}
