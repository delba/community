package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/delba/community/hackerschool"
	"github.com/delba/community/model"
	"github.com/delba/community/vcard"
)

var vcardsPath string

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := hackerschool.Authenticate()
	handle(err)

	fmt.Println("Directory name for the vCards please:")
	_, err = fmt.Scan(&vcardsPath)

	if err := os.Mkdir(vcardsPath, 0755); os.IsExist(err) {
		handle(err)
	}

	c := make(chan []model.Person)
	var wg sync.WaitGroup

	batches, err := hackerschool.GetBatches()
	handle(err)

	for _, batch := range batches {
		go GetPeople(batch, c)
	}

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

func GetPeople(b model.Batch, c chan []model.Person) {
	batchPath := path.Join(vcardsPath, b.Name)

	if err := os.Mkdir(batchPath, 0755); os.IsExist(err) {
		handle(err)
	}

	people, err := hackerschool.GetPeople(&b)
	handle(err)

	c <- people
}

func GenerateVCard(p model.Person, wg *sync.WaitGroup) {
	personPath := path.Join(vcardsPath, p.Batch.Name, p.FormattedName()+".vcard")
	vcard.Generate(personPath, &p)
	wg.Done()
}
