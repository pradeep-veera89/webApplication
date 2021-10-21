package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom from struct embed as url.Values
type Form struct {
	url.Values
	Errors errors
}

// New initializes the from struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if the form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}
