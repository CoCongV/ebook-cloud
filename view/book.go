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
	var (
		books []models.Book
		count int
		prev  bool
		next  bool
	)
	models.DB.Model(&models.Book{}).Count(&count)
	models.DB.Offset(offsetCount).Limit(itemCount).Preload("Authors").Find(&books)
	if page == 1 {
		prev = false
	} else if page > 1 && len(books) > 1 {
		prev = true
	} else {
		prev = false
	}

	if itemCount*page > count {
		next = false
	} else {
		next = true
	}

	c.HTML(http.StatusOK, "index.html", pongo2.Context{
		"books": books,
		"next":  next,
		"prev":  prev,
	})
}
