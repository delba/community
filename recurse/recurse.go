package recurse

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"golang.org/x/oauth2"
)

const BaseURL = "https://www.recurse.com/api/v1"

var accessToken string

var config = oauth2.Config{
	ClientID:     os.Getenv("COMMUNITY_CLIENT_ID"),
	ClientSecret: os.Getenv("COMMUNITY_CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.recurse.com/oauth/authorize",
		TokenURL: "https://www.recurse.com/oauth/token",
	},
	RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func GetRequest(urlStr string) (*http.Request, error) {
	var err error

	u, err := url.Parse(urlStr)

	request := &http.Request{
		Method: "GET",
		URL:    u,
		Header: http.Header{"Authorization": []string{"Bearer " + accessToken}},
	}

	return request, err
}

func Authenticate() error {
	var err error

	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	err = exec.Command("open", url).Run()

	fmt.Println("Authorization grant here please:")
	var code string
	_, err = fmt.Scan(&code)

	token, err := config.Exchange(oauth2.NoContext, code)
	accessToken = token.AccessToken

	return err
}
