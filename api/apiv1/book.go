package apiv1

import (
	"ebook-cloud/config"
	"ebook-cloud/models"
	"ebook-cloud/search"
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
	offsetCount := (page - 1) * config.Conf.PerPageItem

	db := models.DB.Offset(offsetCount).Limit(config.Conf.PerPageItem).Preload("Authors")

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

	if config.Conf.PerPageItem*page > count {
		next = false
	} else {
		next = true
	}

	c.JSON(200, gin.H{
		"books": books,
		"prev":  prev,
		"next":  next,
		"count": count,
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
		log.Println(filename)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	uidInterfae, _ := c.Get("uid")
	uid := uidInterfae.(uint)
	book := models.Book{
		Name:   bookForm.Name,
		File:   dstname,
		UserID: uid,
	}
	if bookForm.AuthorID != 0 {
		models.DB.Find(&author, bookForm.AuthorID)
		book.Authors = []*models.Author{&author}
	}
	models.DB.Create(&book)
	search.BookIndex.Index(strconv.FormatUint(uint64(book.ID), 10), search.IndexData{book.Name})
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
	filename := strings.Join([]string{book.Name, ss[len(ss)-1]}, ".")
	c.FileAttachment(book.File, filename)
}

type UpdateBookJson struct {
	Name    string  `json:"name" binding:"required"`
	Authors *[]uint `json:"authors"`
	Tags    *[]uint `json:"tags"`
}

func PatchBook(c *gin.Context) {
	var (
		bookuri  BookURI
		book     models.Book
		bookJSON UpdateBookJson
		authors  []models.Author
		tags     []models.Tag
	)

	if err := c.ShouldBindUri(&bookuri); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.BindJSON(&bookJSON); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	models.DB.First(&book, bookuri.ID)

	if bookJSON.Authors != nil {
		models.DB.Where("id in (?)", *(bookJSON.Authors)).Find(&authors)
		models.DB.Model(&book).Association("Authors").Replace(authors)
	}
	if bookJSON.Tags != nil {
		models.DB.Where("id in (?)", *(bookJSON.Tags)).Find(&tags)
		models.DB.Model(&book).Association("Tags").Append(tags)
	}

	models.DB.First(&book, bookuri.ID).Update("name", bookJSON.Name)
	c.AbortWithStatusJSON(
		http.StatusOK, gin.H{
			"book": book,
		},
	)
}

func PutBook(c *gin.Context) {
	var (
		bookuri BookURI
		book    models.Book
	)
	if err := c.ShouldBindUri(&bookuri); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	models.DB.Find(&book, bookuri.ID)

	file, _ := c.FormFile("file")
	ss := strings.Split(file.Filename, ".")
	format := ss[len(ss)-1]
	filename := strings.Join([]string{book.Name, format}, ".")
	dst := path.Join(config.Conf.DestPath, filename)
	book.File = dst
	if err := models.DB.Save(&book).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "file existed",
		})
		return
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
