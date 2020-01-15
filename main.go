package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"EookCloud/models"
)

var db *(pg.DB)

func main() {
	db = pg.Connect(&pg.Options{
		models.User: "postgres",
	})
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
