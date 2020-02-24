package server

import (
	"ebook-cloud/config"

	"github.com/gin-gonic/gin"
)

//CreateServ return gin engine
func CreateServ() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	switch config.Conf.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		r.Use(gin.Logger())
		gin.SetMode(gin.DebugMode)
	}
	return r
}
