package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pradeep-veera89/webApplication/internal/config"
)

func TestRoutes(t *testing.T) {
	var app *config.AppConfig

	mux := routes(app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing; test passed
	default:
		t.Error(fmt.Sprintf("%v type is not *chi.Mux", v))
	}
}
