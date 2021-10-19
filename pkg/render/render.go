package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/pradeep-veera89/webApplication/pkg/config"
	"github.com/pradeep-veera89/webApplication/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders with HTML Template file.
func RenderTemplate(w http.ResponseWriter, html string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[html]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	//td = AddDefaultData(td)
	log.Println(td)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser ", err)
	}
}

// CreateTemplateCache creates a template cache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return nil, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return nil, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
