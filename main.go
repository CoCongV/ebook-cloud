package main

import (
	"github.com/gin-gonic/gin"
	// "github.com/go-pg/pg/v9/orm"

	"EbookCloud/app/apiv1"
)

func main() {

	r := gin.Default()
	apiv1.SetRouter(r)
	r.Run()
}
