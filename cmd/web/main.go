package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pradeep-veera89/webApplication/pkg/config"
	"github.com/pradeep-veera89/webApplication/pkg/handlers"
	"github.com/pradeep-veera89/webApplication/pkg/render"
)

const portNumber = ":8080"

// main is the main application function
func main() {

	var app config.AppConfig

	// initializes the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	// creating new Repositories inside the handler package
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// assign the render package with AppConfig
	render.NewTemplates(&app)

	// Other Handler functions
	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting application on port %s", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
