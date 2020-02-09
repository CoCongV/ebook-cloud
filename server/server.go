package server

import (
	"ebook-cloud/config"

	"github.com/gin-gonic/gin"
)

//CreateServ return gin engine
func CreateServ() *gin.Engine {
	r := gin.Default()
	switch config.Conf.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	return r
}
