package vcard

import (
	"os"
	"path"
	"runtime"
	"text/template"

	"github.com/delba/community/models"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func Generate(destinationPath string, person *models.Person) {
	dst, err := os.Create(destinationPath)
	handle(err)
	defer dst.Close()

	_, filename, _, _ := runtime.Caller(1)
	templatePath := path.Join(path.Dir(filename), "templates", "template.vcard")

	t, err := template.ParseFiles(templatePath)
	handle(err)

	err = t.Execute(dst, person)
	handle(err)
}
