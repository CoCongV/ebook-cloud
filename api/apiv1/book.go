package apiv1

import (
	"ebook-cloud/config"
	"ebook-cloud/models"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetBooks get all books
func GetBooks(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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

//BookForm is form struct
type BookForm struct {
	Name     string `form:"name" binding:"required"`
	AuthorID int    `form:"author"`
	Format   string `form:"format" binding:"required"`
	// File       *multipart.File
	// Fileheader *multipart.FileHeader
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

	models.DB.Find(&author, bookForm.AuthorID)

	book := models.Book{
		Name:   bookForm.Name,
		File:   dstname,
		Author: []*models.Author{&author},
	}
	models.DB.Create(&book)
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
	filename := book.Name + ss[len(ss)-1]
	c.FileAttachment(book.File, filename)
}
