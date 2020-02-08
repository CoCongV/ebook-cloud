package view

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/flosch/pongo2"
)

//BookView return book template
func BookView(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", pongo2.Context{"title": "AMD, Yes"})
}