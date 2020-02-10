package view

import (
	"ebook-cloud/models"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

//BookView return book template
func BookView(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	offsetCount := (page - 1) * 20
	itemCount := 20
	var books []models.Book
	models.DB.Offset(offsetCount).Limit(itemCount).Preload("Authors").Find(&books)
	c.HTML(http.StatusOK, "index.html", pongo2.Context{"books": books})
}
