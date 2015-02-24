package hackerschool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/delba/community/model"
	"golang.org/x/oauth2"
)

const (
	ClientID     = "2bb3c608d59c16f3749fb75bfe2f54ad5bc268efcbd6d940a2174fbae078f77f"
	ClientSecret = "bd955c3d40dbf1763e4fc0d8ee38c891c218029bc8d157386a7196d219b21701"

	AuthURL     = "https://www.hackerschool.com/oauth/authorize"
	TokenURL    = "https://www.hackerschool.com/oauth/token"
	RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	BaseURL = "https://www.hackerschool.com/api/v1"
)

var accessToken string

var config = oauth2.Config{
	ClientID:     ClientID,
	ClientSecret: ClientSecret,
	Endpoint: oauth2.Endpoint{
		AuthURL:  AuthURL,
		TokenURL: TokenURL,
	},
	RedirectURL: RedirectURL,
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func GetPeople(batch *model.Batch) ([]model.Person, error) {
	var people []model.Person

	b := bytes.NewBuffer([]byte{})

	req, err := http.NewRequest("GET", fmt.Sprintf(BaseURL+"/batches/%d/people", batch.ID), b)
	if err != nil {
		return people, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
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

	b := bytes.NewBuffer([]byte{})
	req, err := http.NewRequest("GET", BaseURL+"/batches", b)
	if err != nil {
		return batches, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
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
