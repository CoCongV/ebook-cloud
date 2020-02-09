package apiv1

import (
	"ebook-cloud/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//GetAuthors get all authors by pagination
func GetAuthors(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
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
	c.JSON(http.StatusCreated, gin.H{
		"id": author.ID,
	})
}
