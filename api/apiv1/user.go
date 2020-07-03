package apiv1

import (
	"ebook-cloud/client"
	"ebook-cloud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Login for get user info from user server
func Login(c *gin.Context) {
	var user models.User
	token := c.GetHeader("Authorization")
	uid, err := client.UserClient.VerifyUser(token)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
	}

	models.DB.Where("uid = ?", uid).First(&user)
	if user == (models.User{}) {
		models.DB.Create(&models.User{
			UID: uid,
		})
	}
	c.JSON(http.StatusOK, "")
}
