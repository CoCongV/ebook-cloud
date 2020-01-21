package server

import (
	"github.com/gin-gonic/gin"
)

//CreateServ return gin engine
func CreateServ() *gin.Engine {
	return gin.Default()
}
