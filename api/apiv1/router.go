package apiv1

import (
	"github.com/gin-gonic/gin"
)

// SetRouter init api v1 blueprint
func SetRouter(e *gin.Engine) {
	v1 := e.Group("/api/v1")
	v1.GET("/books", GetBooks)
	v1.GET("/books/:id", GetBook)
	v1.GET("/authors", GetAuthors)
	v1.POST("/authors", PostAuthors)
	v1.GET("/countries", GetCountries)
}
