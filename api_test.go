package main

import (
	"ebook-cloud/api/apiv1"
	"ebook-cloud/models"
	"ebook-cloud/server"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuit struct {
	suite.Suite
	server *gin.Engine
}

func (suit *TestSuit) SetupSuite() {
	suit.server = server.CreateServ()
	apiv1.SetRouter(suit.server)
	suit.createData()
}

func (suit *TestSuit) TearDownSuite() {
	suit.delData()
}

func (suit *TestSuit) createData() {
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
		File: "test",
	})
	models.DB.FirstOrCreate(&author, models.Author{
		Name:      "test",
		CountryID: china.ID,
		BookID:    book.ID,
	})
}

func (suit *TestSuit) delData() {
	models.DB.Unscoped().Where("name = ?", "test").Delete(&models.Book{})
	models.DB.Unscoped().Where("name = ?", "test").Delete(&models.Author{})
	models.DB.Unscoped().Where("name = ?", "China").Delete(&models.Country{})

}

func (suit *TestSuit) TestBooks() {
	w := httptest.NewRecorder()
	params := url.Values{}
	params.Set("page", "1")
	query := params.Encode()
	url := "/api/v1/books?"
	oneURL := url + query
	var booksResp struct {
		Books []models.Book `json:"books"`
	}

	req, _ := http.NewRequest("GET", oneURL, nil)
	suit.server.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &booksResp); err != nil {
		log.Fatal(err)
	}
	assert.Equal(suit.T(), 200, w.Code)
	assert.Equal(suit.T(), len(booksResp.Books), 1)

	params.Set("page", "2")
	query = params.Encode()
	towURL := url + query

	req, _ = http.NewRequest("GET", towURL, nil)
	suit.server.ServeHTTP(w, req)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &booksResp); err != nil {
		log.Fatal(err)
	}
	assert.Equal(suit.T(), 200, w.Code)
	assert.Equal(suit.T(), len(booksResp.Books), 0)
}

func TestUserTestSuit(t *testing.T) {
	suite.Run(t, new(TestSuit))
}
