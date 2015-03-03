package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/delba/community/hackerschool"
	"github.com/delba/community/model"
	"github.com/delba/community/vcard"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := hackerschool.Authenticate()
	handle(err)

	fmt.Println("Directory name for the vCards please:")
	var vcardsPath string
	_, err = fmt.Scan(&vcardsPath)

	if err := os.Mkdir(vcardsPath, 0755); os.IsExist(err) {
		handle(err)
	}

	var peopleCount int
	peopleChan := make(chan []model.Person)
	done := make(chan bool)

	batches, err := hackerschool.GetBatches()
	handle(err)

	for _, batch := range batches {
		go func(batch model.Batch) {
			batchPath := path.Join(vcardsPath, batch.Name)

			if err := os.Mkdir(batchPath, 0755); os.IsExist(err) {
				handle(err)
			}

			people, err := hackerschool.GetPeople(&batch)
			handle(err)

			peopleChan <- people
		}(batch)
	}

	for range batches {
		people := <-peopleChan

		for _, person := range people {
			peopleCount++

			go func(person model.Person) {
				batchPath := path.Join(vcardsPath, person.Batch.Name)
				personPath := path.Join(batchPath, person.FormattedName()+".vcard")

				vcard.Generate(personPath, &person)

				done <- true
			}(person)
		}
	}

	for i := 0; i < peopleCount; i++ {
		<-done
	}

	fmt.Println("Your vCards have been generated!")
	err = exec.Command("open", vcardsPath).Run()
}
