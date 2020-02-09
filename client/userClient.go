package client

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/publicsuffix"

	"ebook-cloud/config"
)

//userClient is ...
type userClient struct {
	VerifyURL string
	Client    *resty.Client
}

type VerifyUserResp struct {
	ID uint
}

//UserClient is UserClientStruct point
var UserClient *userClient

//Setup is init UserClient
func Setup() {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	httpClient := http.Client{
		Jar: cookieJar,
	}

	if config.Conf.UserServerTimeout != 0 {
		httpClient.Timeout = time.Duration(config.Conf.UserServerTimeout) * time.Second
	}
	UserClient = &userClient{
		VerifyURL: config.Conf.VerifyUserURL,
		Client:    resty.NewWithClient(&httpClient),
	}
}

func (u *userClient) VerifyUser(token string) (uint, error) {
	var result VerifyUserResp
	// resp, err := u.Client.R().SetHeader("Authorization", token).Get(u.verifyURL)
	resp, err := u.Client.R().SetAuthToken(token).SetResult(&result).Get(u.VerifyURL)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode() == http.StatusOK {
		return result.ID, nil
	}
	return 0, errors.New("Verify False")
}
