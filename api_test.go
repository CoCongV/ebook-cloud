package main

import (
	"bytes"
	"ebook-cloud/api/apiv1"
	"ebook-cloud/models"
	"ebook-cloud/server"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuit struct {
	suite.Suite
	server    *gin.Engine
	countryID uint
	country   *models.Country
	author    *models.Author
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
	})
	models.DB.Model(&author).Association("Books").Append(book)
	suit.country = &china
	suit.author = &author
}

func (suit *TestSuit) delData() {
	models.DB.Unscoped().Delete(&models.Book{})
	models.DB.Unscoped().Delete(&models.Author{})
	models.DB.Unscoped().Delete(&models.Country{})

}

func (suit *TestSuit) TestGetBooks() {
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
		assert.Error(suit.T(), err)
	}
	assert.Equal(suit.T(), 200, w.Code)
	assert.Equal(suit.T(), len(booksResp.Books), 1)

	w = httptest.NewRecorder()
	params.Set("page", "2")
	query = params.Encode()
	towURL := url + query

	req, _ = http.NewRequest("GET", towURL, nil)
	suit.server.ServeHTTP(w, req)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &booksResp); err != nil {
		assert.Error(suit.T(), err)
	}
	assert.Equal(suit.T(), 200, w.Code)
	assert.Equal(suit.T(), 0, len(booksResp.Books))
}

func (suit *TestSuit) TestPostBook() {
	w := httptest.NewRecorder()
	file, err := os.Open("./test_file/test.mobi")
	if err != nil {
		assert.Error(suit.T(), err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.mobi")
	if err != nil {
		assert.Error(suit.T(), err)
	}
	_, err = io.Copy(part, file)
	_ = writer.WriteField("name", "book")
	_ = writer.WriteField("author", fmt.Sprint(suit.author.ID))
	if err = writer.Close(); err != nil {
		assert.Error(suit.T(), err)
	}
	req, _ := http.NewRequest("POST", "/api/v1/books", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 201, w.Code)
}

func (suit *TestSuit) TestAuthors() {
	var (
		authorsResp struct {
			Authors []models.Author `json:"authors"`
		}
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/authors", nil)
	suit.server.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &authorsResp); err != nil {
		assert.Error(suit.T(), err)
	}
	assert.Equal(suit.T(), 200, w.Code)
	assert.Equal(suit.T(), 1, len(authorsResp.Authors))

	w = httptest.NewRecorder()
	params := apiv1.AuthorsReqParams{
		Name:      "test1",
		CountryID: suit.country.ID,
	}
	paramsByte, err := json.Marshal(params)
	if err != nil {
		assert.Error(suit.T(), err)
	}
	req, _ = http.NewRequest("POST", "/api/v1/authors", bytes.NewBuffer(paramsByte))
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 201, w.Code)
}

func (suit *TestSuit) TestCountries() {
	var countriesResp struct {
		Countries []models.Country `json:"countries"`
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/countries", nil)
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 200, w.Code)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &countriesResp); err != nil {
		assert.Error(suit.T(), err)
	}
	assert.Equal(suit.T(), 1, len(countriesResp.Countries))
}

func TestUserTestSuit(t *testing.T) {
	suite.Run(t, new(TestSuit))
}
