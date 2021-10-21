package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	r := httptest.NewRequest("GET", "/whatever", nil)

	form := New(r.PostForm)
	if !form.Valid() {
		t.Error("Invalid Form type created")
	}
}

func TestRequired(t *testing.T) {

	// Invalid Required fields
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("invalid-field", "no-field", "not-present-field")
	if form.Valid() {
		t.Error("shows form is valid, even when invalid fields are passed")
	}

	// Valid Required fields
	r2 := httptest.NewRequest("POST", "/valid-form-test", nil)

	postDate := url.Values{}
	postDate.Add("valid-field-1", "value-1")
	postDate.Add("valid-field-2", "value-2")
	postDate.Add("valid-field-3", "value-3")
	r2.PostForm = postDate

	form2 := New(r2.PostForm)
	form2.Required("valid-field-1", "valid-field-2", "valid-field-3")
	if !form2.Valid() {
		t.Error("shows form as invalid , even when valid fields are passed")
	}
}

func TestHas(t *testing.T) {
	// Valid Test
	r := httptest.NewRequest("POST", "/whatever", nil)
	postDate := url.Values{}
	postDate.Add("valid-field-1", "value-1")
	postDate.Add("valid-field-2", "value-2")
	postDate.Add("valid-field-3", "value-3")
	r.PostForm = postDate
	form := New(r.PostForm)

	if !form.Has("valid-field-1", r) {
		t.Error("shows value not present, even when fieldvalue is present for the field ")
	}
	// Invalid Test
	if form.Has("not-present-field", r) {
		t.Error("shows value is present, even when fieldvalue is not present for the field ")
	}
}

func TestMinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	postDate := url.Values{}
	postDate.Add("valid-field-1", "value-1")
	postDate.Add("valid-field-2", "value-2")
	postDate.Add("valid-field-3", "value-3")
	r.PostForm = postDate
	form := New(r.PostForm)
	form.MinLength("valid-field-1", 2, r)
	if !form.Valid() {
		t.Error("failed to check for minlength for valid data")
	}

	form.MinLength("valid-field-2", 20, r)
	if form.Valid() {
		t.Error("failed to check for minlength for valid data")
	}
}

func TestIsEmail(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever", nil)
	postDate := url.Values{}
	postDate.Add("email", "test@email.com")
	postDate.Add("invalid-email", "test@test")
	r.PostForm = postDate
	form := New(r.PostForm)

	// valid check
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("shows invalid email, even valid email is passed")
	}

	// invalid check
	form.IsEmail("invalid-email")
	if form.Valid() {
		t.Error("shows valid email, even invalid email is passed")
	}
}
