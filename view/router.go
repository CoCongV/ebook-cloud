package view

import (
	"github.com/gin-gonic/gin"
)

func SetRouter(e *gin.Engine) {
	view := e.Group("/")
	view.GET("", BookView)
}
