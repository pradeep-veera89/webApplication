package forms

import (
	"testing"
)

func TestAdd(t *testing.T) {

	e := errors{}
	field := "field-name"
	message := "field error message"
	e.Add(field, message)
	if e.Get(field) != message {
		t.Error("failed to add error message")
	}

}

func TestGet(t *testing.T) {
	e := errors{}
	field := "field-name"
	message := "field error message"
	e.Add(field, message)
	if e.Get(field) != message {
		t.Error("failed to add error message")
	}

	if e.Get("not-present-field") != "" {
		t.Error("shows not empty message for the field not present inside the error struct")
	}
}
