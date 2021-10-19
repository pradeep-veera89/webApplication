package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

// RenderTemplate renders with HTML Template file.
func RenderTemplate(w http.ResponseWriter, html string) {

	 
	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	t, ok := tc[html]
	if !ok {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, nil)
	_, err = buf.WriteTo(w)
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
