package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	h := NoSurf(&myHandler{})

	switch v := h.(type) {
	case http.Handler:

	default:
		t.Error(fmt.Sprintf(" %T type is not http.Hanlder", v))
	}
}

func TestSessionLoad(t *testing.T) {

	h := SessionLoad(&myHandler{})

	switch v := h.(type) {
	case http.Handler:

	default:
		t.Error(fmt.Sprintf(" %T type is not http.Hanlder", v))
	}
}
