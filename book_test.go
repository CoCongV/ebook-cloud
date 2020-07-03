package main

import (
	"bytes"
	"ebook-cloud/api/apiv1"
	"ebook-cloud/client"
	"ebook-cloud/config"
	"ebook-cloud/models"
	"ebook-cloud/search"
	"ebook-cloud/server"
	"ebook-cloud/view"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuit struct {
	suite.Suite
	server    *gin.Engine
	countryID uint
	country   *models.Country
	author    *models.Author
	author1   *models.Author
	book      *models.Book
	tag       *models.Tag
}

func (suit *TestSuit) SetupSuite() {
	suit.server = server.CreateServ()
	apiv1.SetRouter(suit.server)
	view.SetRouter(suit.server)
	suit.createData()
	mock()
}

func (suit *TestSuit) TearDownSuite() {
	suit.delData()
	httpmock.DeactivateAndReset()
}

func mock() {
	responder, _ := httpmock.NewJsonResponder(200, map[string]uint{
		"id": 1,
	})
	httpmock.ActivateNonDefault(client.UserClient.Client.GetClient())
	httpmock.RegisterResponder("GET", client.UserClient.VerifyURL, responder)
}

func (suit *TestSuit) createData() {
	var (
		china   models.Country
		author  models.Author
		author1 models.Author
		book    models.Book
		tag     models.Tag
	)
	models.NewRoles(1)
	models.DB.FirstOrCreate(&china, models.Country{
		Name: "China",
	})
	models.DB.FirstOrCreate(&tag, models.Tag{
		Name: "test",
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
	search.BookIndex.Index(fmt.Sprint(book.ID), search.IndexData{book.Name})

	models.DB.FirstOrCreate(&author, models.Author{
		Name:      "test",
		CountryID: china.ID,
	})
	models.DB.FirstOrCreate(&author1, models.Author{
		Name:      "test1",
		CountryID: china.ID,
	})

	models.DB.Model(&author).Association("Books").Append(book)
	suit.country = &china
	suit.author = &author
	suit.author1 = &author1
	suit.book = &book
	suit.tag = &tag
}

func (suit *TestSuit) delData() {
	models.DB.Unscoped().Delete(&models.Book{})
	models.DB.Unscoped().Delete(&models.Author{})
	models.DB.Unscoped().Delete(&models.Country{})
	models.DB.Unscoped().Delete(&models.User{})
	models.DB.Unscoped().Delete(&models.Role{})
	os.RemoveAll(config.Conf.BookSearchIndexFile)
}

func (suit *TestSuit) TestGetBooks() {
	w := httptest.NewRecorder()
	oneURL := CreateQuery("/api/v1/books", map[string]string{"page": "1"})
	var booksResp struct {
		Books []models.Book `json:"books"`
	}

	req, _ := http.NewRequest("GET", oneURL, nil)
	suit.server.ServeHTTP(w, req)
	CustomUnmarshal(w, &booksResp, suit.T())
	assert.Equal(suit.T(), 200, w.Code)
	assert.Equal(suit.T(), len(booksResp.Books), 1)

	w = httptest.NewRecorder()
	twoURL := CreateQuery("/api/v1/books", map[string]string{"page": "2"})

	req, _ = http.NewRequest("GET", twoURL, nil)
	suit.server.ServeHTTP(w, req)
	CustomUnmarshal(w, &booksResp, suit.T())
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
	writer.WriteField("name", "book")
	writer.WriteField("format", "mobi")
	writer.WriteField("author", fmt.Sprint(suit.author.ID))
	if err = writer.Close(); err != nil {
		assert.Error(suit.T(), err)
	}
	req, _ := http.NewRequest("POST", "/api/v1/books", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 201, w.Code)
}

func (suit *TestSuit) TestPutBook() {
	w := httptest.NewRecorder()
	file, err := os.Open("./test_file/test1.mobi")
	if err != nil {
		assert.Error(suit.T(), err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test1.mobi")
	if err != nil {
		assert.Error(suit.T(), err)
	}
	_, err = io.Copy(part, file)
	if err = writer.Close(); err != nil {
		assert.Error(suit.T(), err)
	}
	url := "/api/v1/books/" + fmt.Sprint(suit.book.ID)
	req, _ := http.NewRequest("PUT", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 200, w.Code)
}

func (suit *TestSuit) TestGetBookByID() {
	w := httptest.NewRecorder()
	url := "/api/v1/books/" + fmt.Sprint(suit.book.ID)
	req, _ := http.NewRequest("GET", url, nil)
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 200, w.Code)
}

func (suit *TestSuit) TestBooks400() {
	w := httptest.NewRecorder()
	url := CreateQuery("/api/v1/books", map[string]string{"page": "s"})
	req, _ := http.NewRequest("GET", url, nil)
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 400, w.Code)
}

func (suit *TestSuit) TestQueryBook() {
	var booksResp struct {
		Books []models.Book `json:"books"`
	}
	w := httptest.NewRecorder()
	url := CreateQuery("/api/v1/books", map[string]string{"page": "1", "name": "test"})
	req, _ := http.NewRequest("GET", url, nil)
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 200, w.Code)
	CustomUnmarshal(w, &booksResp, suit.T())
	assert.Equal(suit.T(), 1, len(booksResp.Books))
}

func (suit *TestSuit) TestPatchBook() {
	var patchResp struct {
		Book models.Book `json:"book"`
	}
	reqJSON := struct {
		Name    string `json:"name"`
		Authors []uint `json:"authors"`
		Tags    []uint `json:"tags"`
	}{
		Name:    "patchtest",
		Authors: []uint{suit.author.ID, suit.author1.ID},
		Tags:    []uint{suit.tag.ID},
	}
	paramsBytes, err := json.Marshal(reqJSON)
	if err != nil {
		assert.Error(suit.T(), err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/books/"+fmt.Sprint(suit.book.ID), bytes.NewBuffer(paramsBytes))
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 200, w.Code)
	CustomUnmarshal(w, &patchResp, suit.T())
	assert.Equal(suit.T(), "patchtest", patchResp.Book.Name)
	assert.Equal(suit.T(), suit.author.ID, patchResp.Book.Authors[0].ID)
	assert.Equal(suit.T(), suit.author1.ID, patchResp.Book.Authors[1].ID)
	assert.Equal(suit.T(), suit.tag.ID, patchResp.Book.Tags[0].ID)
}

func (suit *TestSuit) TestBook400() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/books/s", nil)
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 400, w.Code)
}

func (suit *TestSuit) TestBook404() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/books/"+fmt.Sprint(suit.book.ID+1000), nil)
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 404, w.Code)
}

func (suit *TestSuit) TestCountries() {
	var countriesResp struct {
		Countries []models.Country `json:"countries"`
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/countries", nil)
	suit.server.ServeHTTP(w, req)
	assert.Equal(suit.T(), 200, w.Code)
	CustomUnmarshal(w, &countriesResp, suit.T())
	assert.Equal(suit.T(), 1, len(countriesResp.Countries))
}

// func (suit *TestSuit) TestBookTemplate() {
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	suit.server.ServeHTTP(w, req)
// 	assert.Equal(suit.T(), 200, w.Code)
// }

func TestApiTestSuit(t *testing.T) {
	suite.Run(t, new(TestSuit))
}

func CreateQuery(baseURL string, params map[string]string) string {
	r := url.Values{}
	for k, v := range params {
		r.Set(k, v)
	}
	return baseURL + "?" + r.Encode()
}

func CustomUnmarshal(w *httptest.ResponseRecorder, r interface{}, t assert.TestingT) {
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, r); err != nil {
		assert.Error(t, err)
	}
}
