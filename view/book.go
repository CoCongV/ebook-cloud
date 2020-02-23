package view

import (
	"ebook-cloud/models"
	"ebook-cloud/search"
	"log"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

//BookView return book template
func BookView(c *gin.Context) {
	var (
		books   []models.Book
		count   int
		prev    bool
		next    bool
		idSlice []int
	)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	queryName, ok := c.GetQuery("name")
	offsetCount := (page - 1) * 20
	itemCount := 20

	db := models.DB.Offset(offsetCount).Limit(itemCount).Preload("Authors")

	if ok == false {
		models.DB.Model(&models.Book{}).Count(&count)
		db.Find(&books)
	} else {
		bleveQuery := bleve.NewMatchQuery(queryName)
		searchReq := bleve.NewSearchRequest(bleveQuery)
		searchResults, err := search.BookIndex.Search(searchReq)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		for _, s := range searchResults.Hits {
			id, err := strconv.Atoi(s.ID)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			idSlice = append(idSlice, id)
			log.Println(id)
		}
		models.DB.Where("id in (?)", idSlice).Model(&models.Book{}).Count(&count)
		db.Where("id in (?)", idSlice).Find(&books)
	}

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
		"page":  page,
	})
}
