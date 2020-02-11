package apiv1

import (
	"ebook-cloud/config"
	"ebook-cloud/models"
	"ebook-cloud/search"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/gin-gonic/gin"
)

// GetBooks get all books
func GetBooks(c *gin.Context) {
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

	models.DB.Model(&models.Book{}).Count(&count)
	db := models.DB.Offset(offsetCount).Limit(itemCount)

	if ok == false {
		db.Find(&books)
	} else {
		bleveQuery := bleve.NewMatchQuery(queryName)
		searchReq := bleve.NewSearchRequest(bleveQuery)
		searchResults, err := search.Index.Search(searchReq)
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

	c.JSON(200, gin.H{
		"books": books,
		"prev":  prev,
		"next":  next,
	})
}

//BookForm is form struct
type BookForm struct {
	Name     string `form:"name" binding:"required"`
	AuthorID int    `form:"author"`
	Format   string `form:"format" binding:"required"`
}

//PostBooks is create and save book
func PostBooks(c *gin.Context) {
	var (
		bookForm BookForm
		author   models.Author
	)
	c.Bind(&bookForm)
	file, _ := c.FormFile("file")
	filename := strings.Join([]string{bookForm.Name, bookForm.Format}, ".")
	dstname := path.Join(config.Conf.DestPath, filename)
	if err := c.SaveUploadedFile(file, dstname); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// uidInterfae, _ := c.Get("uid")
	// uid := uidInterfae.(uint)
	book := models.Book{
		Name: bookForm.Name,
		File: dstname,
		// UserID: uid,
	}
	if bookForm.AuthorID != 0 {
		models.DB.Find(&author, bookForm.AuthorID)
		book.Authors = []*models.Author{&author}
	}
	models.DB.Create(&book)
	search.Index.Index(fmt.Sprint(book.ID), search.BookIndex{book.Name})
	c.JSON(http.StatusCreated, gin.H{
		"id": book.ID,
	})
}

type BookURI struct {
	ID int `uri:"id" binding:"required"`
}

func GetBook(c *gin.Context) {
	var (
		bookuri BookURI
		book    models.Book
	)
	if err := c.ShouldBindUri(&bookuri); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	models.DB.First(&book, bookuri.ID)
	if book.ID == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	ss := strings.Split(book.File, ".")
	filename := book.Name + "." + ss[len(ss)-1]
	c.FileAttachment(book.File, filename)
}
