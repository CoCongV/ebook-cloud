package main

import (
	"bytes"
	"ebook-cloud/api/apiv1"
	"ebook-cloud/config"
	"ebook-cloud/models"
	"ebook-cloud/search"
	"ebook-cloud/server"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthorSuit struct {
	suite.Suite
	server    *gin.Engine
	countryID uint
	country   *models.Country
	author    *models.Author
	book      *models.Book
}

func (suit *AuthorSuit) SetupSuite() {
	suit.server = server.CreateServ()
	apiv1.SetRouter(suit.server)
	suit.createData()
	mock()
}

func (suit *AuthorSuit) createData() {
	var (
		china  models.Country
		author models.Author
		book   models.Book
	)
	models.DB.FirstOrCreate(&china, models.Country{
		Name: "China",
	})

	models.DB.FirstOrCreate(&book, models.Book{
		Name: "test",
		File: path.Join(config.Conf.DestPath, "test.mobi"),
	})
	dst, _ := os.Create(book.File)
	src, err := os.Open("./test_file/test.mobi")
	if err != nil {
		assert.Error(suit.T(), err)
	}
	io.Copy(dst, src)
	search.BookIndex.Index(fmt.Sprint(book.ID), search.BookIndexData{
		Name: book.Name,
	})

	models.DB.FirstOrCreate(&author, models.Author{
		Name:      "test",
		CountryID: china.ID,
	})
	models.DB.Model(&author).Association("Books").Append(book)
	suit.country = &china
	suit.author = &author
	suit.book = &book
}

func (suit *AuthorSuit) delData() {
	models.DB.Unscoped().Delete(&models.Book{})
	models.DB.Unscoped().Delete(&models.Author{})
	models.DB.Unscoped().Delete(&models.Country{})
	os.RemoveAll(config.Conf.BookSearchIndexFile)
}

func (suit *AuthorSuit) TearDownSuite() {
	suit.delData()
	httpmock.DeactivateAndReset()
}

func (suit *AuthorSuit) TestAuthors() {
	var (
		authorsResp struct {
			Authors []models.Author `json:"authors"`
		}
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/authors", nil)
	suit.server.ServeHTTP(w, req)
	CustomUnmarshal(w, &authorsResp, suit.T())
	assert.Equal(suit.T(), 200, w.Code)
	assert.Equal(suit.T(), 1, len(authorsResp.Authors))
}

func (suit *AuthorSuit) TestPostAuthor() {
	w := httptest.NewRecorder()
	params := apiv1.AuthorsReqParams{
		Name:      "test1",
		CountryID: suit.country.ID,
	}
	paramsByte, err := json.Marshal(params)
	if err != nil {
		assert.Error(suit.T(), err)
	}
	req, _ := http.NewRequest("POST", "/api/v1/authors", bytes.NewBuffer(paramsByte))
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 201, w.Code)
}

func (suit *AuthorSuit) TestAuthors400() {
	w := httptest.NewRecorder()
	url := CreateQuery("/api/v1/authors", map[string]string{"page": "s"})
	req, _ := http.NewRequest("GET", url, nil)
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 400, w.Code)

	w = httptest.NewRecorder()
	params := apiv1.AuthorsReqParams{
		Name: "test1",
	}
	paramsByte, err := json.Marshal(params)
	if err != nil {
		assert.Error(suit.T(), err)
	}
	req, _ = http.NewRequest("POST", "/api/v1/authors", bytes.NewBuffer(paramsByte))
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 400, w.Code)

	w = httptest.NewRecorder()
	params = apiv1.AuthorsReqParams{
		Name:      "test1",
		CountryID: 1000000,
	}
	paramsByte, err = json.Marshal(params)
	if err != nil {
		assert.Error(suit.T(), err)
	}
	req, _ = http.NewRequest("POST", "/api/v1/authors", bytes.NewBuffer(paramsByte))
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 400, w.Code)
}

func TestAuthorSuit(t *testing.T) {
	suite.Run(t, new(AuthorSuit))
}
