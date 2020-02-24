package apiv1

import (
	"ebook-cloud/config"
	"ebook-cloud/models"
	"ebook-cloud/search"
	"log"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/gin-gonic/gin"
)

//GetAuthors get all authors by pagination
func GetAuthors(c *gin.Context) {
	var (
		count   int
		authors []models.Author
		idSlice []int
	)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	queryName, ok := c.GetQuery("name")
	offsetCount := (page - 1) * config.Conf.PerPageItem

	db := models.DB.Offset(offsetCount).Limit(config.Conf.PerPageItem)
	if ok == false {
		models.DB.Model(&models.Book{}).Count(&count)
		db.Find(&authors)
	} else {
		bleveQuery := bleve.NewMatchQuery(queryName)
		searchReq := bleve.NewSearchRequest(bleveQuery)
		searchResults, err := search.AuthorIndex.Search(searchReq)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		for _, s := range searchResults.Hits {
			id, err := strconv.Atoi(s.ID)
			if err != nil {
				log.Println(err)
				continue
			}
			idSlice = append(idSlice, id)
		}
		models.DB.Where("id in (?)", idSlice).Model(&models.Book{}).Count(&count)
		db.Where("id in (?)", idSlice).Find(&authors)
	}

	h := gin.H{
		"authors": authors,
	}
	c.JSON(200, h)
}

//AuthorsReqParams is post authors json struct
type AuthorsReqParams struct {
	Name      string `json:"name" binding:"required"`
	CountryID uint   `json:"country_id" binding:"required"`
}

//PostAuthors is ...
func PostAuthors(c *gin.Context) {
	var (
		authorReq AuthorsReqParams
		country   models.Country
	)
	if err := c.BindJSON(&authorReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	models.DB.Find(&country, authorReq.CountryID)
	if country.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "country not found",
		})
		return
	}
	uidInterface, _ := c.Get("uid")
	uid := uidInterface.(uint)
	author := models.Author{
		Name:      authorReq.Name,
		CountryID: authorReq.CountryID,
		UserID:    uid,
	}
	models.DB.Create(&author)
	search.AuthorIndex.Index(strconv.FormatUint(uint64(author.ID), 10), search.IndexData{
		Name: author.Name,
	})
	c.JSON(http.StatusCreated, gin.H{
		"id": author.ID,
	})
}
