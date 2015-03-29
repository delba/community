package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/delba/community/models"
	"github.com/delba/community/recurse"
	"github.com/delba/community/vcard"
)

var vcardsPath string

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := recurse.Authenticate()
	handle(err)

	fmt.Println("Directory name for the vCards please:")
	_, err = fmt.Scan(&vcardsPath)

	if err := os.Mkdir(vcardsPath, 0755); os.IsExist(err) {
		handle(err)
	}

	var batches models.Batches
	err = batches.Fetch()
	handle(err)

	c := make(chan []models.Person)

	for _, batch := range batches {
		go GetPeople(batch, c)
	}

	var wg sync.WaitGroup

	for range batches {
		for _, person := range <-c {
			wg.Add(1)
			go GenerateVCard(person, &wg)
		}
	}

	wg.Wait()

	fmt.Println("Your vCards have been generated!")
	err = exec.Command("open", vcardsPath).Run()
	handle(err)
}

func GetPeople(b models.Batch, c chan []models.Person) {
	batchPath := path.Join(vcardsPath, b.Name)

	if err := os.Mkdir(batchPath, 0755); os.IsExist(err) {
		handle(err)
	}

	err := b.FetchPeople()
	handle(err)

	c <- b.People
}

func GenerateVCard(p models.Person, wg *sync.WaitGroup) {
	personPath := path.Join(vcardsPath, p.Batch.Name, p.FormattedName()+".vcard")
	vcard.Generate(personPath, &p)
	wg.Done()
}
