package apiv1

import (
	"ebook-cloud/api"
	"ebook-cloud/models"

	"github.com/gin-gonic/gin"
)

// SetRouter init api v1 blueprint
func SetRouter(e *gin.Engine) {
	v1 := e.Group("/api/v1")
	v1.GET("/books", GetBooks)
	v1.POST("/books", api.AuthHandler(models.Moderator), PostBooks)
	v1.GET("/books/:id", GetBook)
	v1.PATCH("/books/:id", api.AuthHandler(models.Moderator), PatchBook)
	v1.PUT("/books/:id", api.AuthHandler(models.Moderator), PutBook)
	v1.GET("/authors", GetAuthors)
	v1.POST("/authors", api.AuthHandler(models.Moderator), PostAuthors)
	v1.GET("/countries", GetCountries)
}
