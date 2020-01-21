package apiv1

import (
	"ebook-cloud/models"

	"github.com/gin-gonic/gin"
)

//GetCountries get all authors by pagination
func GetCountries(c *gin.Context) {

	var countries []models.Country
	models.DB.Find(&countries)
	h := gin.H{
		"countries": countries,
	}
	c.JSON(200, h)
}
