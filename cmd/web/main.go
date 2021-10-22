package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pradeep-veera89/webApplication/internal/config"
	"github.com/pradeep-veera89/webApplication/internal/driver"
	"github.com/pradeep-veera89/webApplication/internal/handlers"
	"github.com/pradeep-veera89/webApplication/internal/helpers"
	"github.com/pradeep-veera89/webApplication/internal/models"
	"github.com/pradeep-veera89/webApplication/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	//  What am i goin to store in Session
	gob.Register(models.Reservation{})
	// Change this to true when in production
	app.InProduction = false

	// initializing infoLog
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// initializing errorLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// initializing session from scs sessionManager.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Println("Connecting to Database ....")
	db, err := driver.ConnectSQL("host=localhost post=5432 dbname=bookings user= password=")
	if err != nil {
		log.Fatal("Cannot connect to DB ")
	}
	log.Println("Connected to DB ...")

	// initializes the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	// creating new Repositories inside the handler package
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// assign the render package with AppConfig
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
