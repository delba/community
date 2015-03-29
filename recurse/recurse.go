package recurse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/delba/community/model"
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

func GetPeople(batch *model.Batch) ([]model.Person, error) {
	var people []model.Person

	request, err := getRequest(fmt.Sprintf(BaseURL+"/batches/%d/people", batch.ID))
	if err != nil {
		return people, err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return people, err
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return people, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(contents, &people)
	if err != nil {
		return people, err
	}

	return people, nil
}

func GetBatches() ([]model.Batch, error) {
	var batches []model.Batch

	request, err := getRequest(BaseURL + "/batches")
	if err != nil {
		return batches, err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return batches, err
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return batches, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(contents, &batches)
	if err != nil {
		return batches, err
	}

	return batches, nil
}

func getRequest(urlStr string) (*http.Request, error) {
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
