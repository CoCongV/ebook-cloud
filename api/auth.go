package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ebook-cloud/client"
	"ebook-cloud/models"
)

//AuthHandler is auth decorator
func AuthHandler(permission string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			user           models.User
			userRole       models.Role
			permissionRole models.Role
		)
		token := c.GetHeader("Authorization")
		id, err := client.UserClient.VerifyUser(token)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		models.DB.Where("name = ?", permission).First(&permissionRole)
		models.DB.Where("UID = ?", id).First(&user)
		if user == (models.User{}) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		models.DB.Model(&user).Related(&userRole)
		if permissionRole.HasPermission(permissionRole.Permmission) {
			c.Set("uid", id)
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}
