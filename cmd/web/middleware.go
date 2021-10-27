package main

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/pradeep-veera89/webApplication/internal/helpers"
)

// NoSurf adds CSRF Protection to ll POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and savess the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Inside Auth")
		if !helpers.IsAuthenticated(r) {
			log.Println("Check if User is Authenticated")
			session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			//return
		}
		next.ServeHTTP(w, r)
	})
}
