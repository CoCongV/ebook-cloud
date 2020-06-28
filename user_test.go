package main

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
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

func (suit *UserSuit) TestCreateRoles() {
	models.NewRoles(0)
}

func TestUserSuit(t *testing.T) {
	suite.Run(t, new(UserSuit))
}
