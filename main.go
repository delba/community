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

	var peopleWG sync.WaitGroup
	var vcardsWG sync.WaitGroup

	batches, err := hackerschool.GetBatches()
	handle(err)

	for _, batch := range batches {
		peopleWG.Add(1)
		go func(batch model.Batch) {
			batchPath := path.Join(vcardsPath, batch.Name)

			if err := os.Mkdir(batchPath, 0755); os.IsExist(err) {
				handle(err)
			}

			people, err := hackerschool.GetPeople(&batch)
			handle(err)

			for _, person := range people {
				vcardsWG.Add(1)
				go func(person model.Person) {
					personPath := path.Join(batchPath, person.FormattedName()+".vcard")
					vcard.Generate(personPath, &person)
					vcardsWG.Done()
				}(person)
			}

			peopleWG.Done()
		}(batch)
	}

	peopleWG.Wait()
	vcardsWG.Wait()

	fmt.Println("Your vCards have been generated!")
	err = exec.Command("open", vcardsPath).Run()
}
