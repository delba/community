package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/delba/community/recurse"
)

type Project struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

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

type Link struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Person struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	MiddleName       string    `json:"middle_name"`
	LastName         string    `json:"last_name"`
	PhoneNumber      string    `json:"phone_number"`
	Email            string    `json:"email"`
	Github           string    `json:"github"`
	Twitter          string    `json:"twitter"`
	Links            []Link    `json:"links"`
	Bio              string    `json:"bio"`
	HasPhoto         bool      `json:"has_photo"`
	Image            string    `json:"image"`
	IsFaculty        bool      `json:"is_faculty"`
	IsHackerSchooler bool      `json:"is_hacker_schooler"`
	Job              string    `json:"job"`
	Projects         []Project `json:"projects"`
	Skills           []string  `json:"skills"`
	BatchID          int       `json:"batch_id"`
	Batch            *Batch    `json:"batch"`
	BatchName        string
}

func (p *Person) FormattedName() string {
	return strings.Join([]string{p.FirstName, p.MiddleName, p.LastName}, " ")
}
