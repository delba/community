package model

import "testing"

func TestFormattedName(t *testing.T) {
	person := &Person{
		FirstName:  "John",
		MiddleName: "Middle",
		LastName:   "Doe",
	}

	actual := person.FormattedName()
	expectation := "John Middle Doe"

	if actual != expectation {
		t.Error("Expected '%s'. Got '%s'.", expectation, actual)
	}
}
