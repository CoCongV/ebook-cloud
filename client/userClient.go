package client

import (
	"net/http"
	"net/http/cookiejar"
	"path"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/publicsuffix"

	"ebook-cloud/config"
)

//userClient is ...
type userClient struct {
	baseURL   string
	verifyURL string
	client    *resty.Client
}

type verifyUserResp struct {
	id int
}

//UserClient is UserClientStruct point
var UserClient *userClient

//Setup is init UserClient
func Setup() {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	httpClient := http.Client{
		Jar: cookieJar,
	}
	UserClient = &userClient{
		baseURL:   config.Conf.UserServerURL,
		verifyURL: path.Join(config.Conf.UserServerURL, "/api/v1/verify_auth_token"),
		client:    resty.NewWithClient(&httpClient),
	}
}

func (u *userClient) VerifyUser(token string) (int, bool, error) {
	resp, err := u.client.R().SetHeader("Authorization", token).Get(u.verifyURL)
	if err != nil {
		return 0, false, err
	}
	if resp.StatusCode() == http.StatusOK {
		return resp.Result().(verifyUserResp).id, true, nil
	} else {
		return 0, false, nil
	}
}
