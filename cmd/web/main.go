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

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Printf("Starting application on port %s", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
