package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/delba/community/recurse"
)

type Batch struct {
	EndDate   string `json:"end_date"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	People    []Person
}

func (b *Batch) FetchPeople() error {
	url := fmt.Sprintf("%s/batches/%d/people", recurse.BaseURL, b.ID)
	request, err := recurse.GetRequest(url)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &b.People)

	for _, person := range b.People {
		person.BatchName = b.Name
	}

	return err
}

type Batches []Batch

func (b *Batches) Fetch() error {
	url := fmt.Sprintf("%s/batches", recurse.BaseURL)

	request, err := recurse.GetRequest(url)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &b)

	return err
}
