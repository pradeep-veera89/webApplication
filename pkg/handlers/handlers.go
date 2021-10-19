package handlers

import (
	"fmt"
	"net/http"

	"github.com/pradeep-veera89/webApplication/pkg/config"
	"github.com/pradeep-veera89/webApplication/pkg/models"
	"github.com/pradeep-veera89/webApplication/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

//Repository is a repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote-ip", remoteIp)
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "hellow"

	remoteIp := m.App.Session.Get(r.Context(), "remote-ip")
	stringMap["remote_ip"] = fmt.Sprintf("%v", remoteIp)

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
