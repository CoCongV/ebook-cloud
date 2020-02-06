package client

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/publicsuffix"

	"ebook-cloud/config"
)

//UserClientStruct is ...
type UserClientStruct struct {
	url    string
	client *resty.Client
}

//UserClient is UserClientStruct point
var UserClient *UserClientStruct

//Setup is init UserClient
func Setup() {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	httpClient := http.Client{
		Jar: cookieJar,
	}
	UserClient = &UserClientStruct{
		url:    config.Conf.UserServerURL,
		client: resty.NewWithClient(&httpClient),
	}
}
