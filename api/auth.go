package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ebook-cloud/client"
)

//AuthHandler is auth decorator
func AuthHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	id, err := client.UserClient.VerifyUser(token)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
	}
	c.Set("uid", id)
}
