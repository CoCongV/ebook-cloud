package apiv1

import (
	"ebook-cloud/models"

	"github.com/gin-gonic/gin"
)

//GetAuthors get all authors by pagination
func GetAuthors(c *gin.Context) {
	var authors []models.Author
	models.DB.Take(&authors)
	h := gin.H{
		"message": "json",
	}
	c.JSON(200, h)
}
