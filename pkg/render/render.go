package render

import (
	"log"
	"net/http"
	"text/template"
)

// RenderTemplate renders with HTML Template file.
func RenderTemplate(w http.ResponseWriter, html string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + html)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		log.Println("error parsing template ", err)
		return
	}
}
