package apiv1

import (
	"ebook-cloud/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

//GetAuthors get all authors by pagination
func GetAuthors(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithError(401, err)
		return
	}
	offsetCount := (page - 1) * 20
	itemCount := 20

	var authors []models.Author
	models.DB.Offset(offsetCount).Limit(itemCount).Order("name").Find(&authors)
	h := gin.H{
		"authors": authors,
	}
	c.JSON(200, h)
}
