package model

import "strings"

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
}

func (p *Person) FormattedName() string {
	return strings.Join([]string{p.FirstName, p.MiddleName, p.LastName}, " ")
}
