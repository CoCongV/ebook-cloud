package main

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"ebook-cloud/api/apiv1"
	"ebook-cloud/config"
	"ebook-cloud/models"
	"ebook-cloud/server"
	"ebook-cloud/view"
)

type UserSuit struct {
	suite.Suite
	server *gin.Engine
}

func (suit *UserSuit) SetupSuite() {
	suit.server = server.CreateServ()
	apiv1.SetRouter(suit.server)
	view.SetRouter(suit.server)
	suit.CreateRoles()
}

func (suit *UserSuit) TearDownSuite() {
	suit.delData()
}

func (suit *UserSuit) delData() {
	models.DB.Unscoped().Delete(&models.Book{})
	models.DB.Unscoped().Delete(&models.Author{})
	models.DB.Unscoped().Delete(&models.Country{})
	models.DB.Unscoped().Delete(&models.User{})
	models.DB.Unscoped().Delete(&models.Role{})
	os.RemoveAll(config.Conf.BookSearchIndexFile)
}

func (suit *UserSuit) CreateRoles() {
	models.NewRoles(1)
}

func (suit *UserSuit) TestQueryRole() {
	var (
		user models.User
		role models.Role
	)
	models.DB.Where("UID = ?", 1).First(&user)
	models.DB.Model(&user).Related(&role)
	assert.Equal(suit.T(), models.Administrator, role.Name)
}

func (suit *UserSuit) TestCRUDUser() {
	var (
		role models.Role
	)
	user := models.User{UID: 2}
	assert.Equal(suit.T(), true, models.DB.NewRecord(user))
	models.DB.Create(&user)
	assert.Equal(suit.T(), false, models.DB.NewRecord(user))
	models.DB.Model(&user).Related(&role)
	assert.Equal(suit.T(), models.CommonUser, role.Name)
}

func TestUserSuit(t *testing.T) {
	suite.Run(t, new(UserSuit))
}
