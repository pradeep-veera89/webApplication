package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pradeep-veera89/webApplication/internal/config"
	"github.com/pradeep-veera89/webApplication/internal/driver"
	"github.com/pradeep-veera89/webApplication/internal/forms"

	//"github.com/pradeep-veera89/webApplication/internal/helpers"
	"github.com/pradeep-veera89/webApplication/internal/models"
	"github.com/pradeep-veera89/webApplication/internal/render"
	"github.com/pradeep-veera89/webApplication/internal/repository"
	"github.com/pradeep-veera89/webApplication/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

//Repository is a repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

// Reservation is the make reservations page.
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.InfoLog.Println("RoomId from reservation", reservation.RoomId)
	// get room details by id
	room, err := m.DB.GetRoomById(reservation.RoomId)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find a room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// update the reservatio with room name
	reservation.Room.RoomName = room.RoomName

	// save reservation details into session again
	m.App.Session.Put(r.Context(), "reservation", reservation)
	// reversal converting Date from Go layout to string.
	layout := "2006-01-02"
	sd := reservation.StarDate.Format(layout)
	ed := reservation.EndDate.Format(layout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the post data of the reservation form.
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// check if there are any errors while parsing the form.
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "failed to parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// Get reservation details from session.
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "could not get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// update reservation from user information
	reservation.Email = r.Form.Get("email")
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "invalid form data", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	var newReservationId int
	newReservationId, err = m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation.Id = newReservationId
	m.App.Session.Put(r.Context(), "reservation", reservation)
	// insert into Room data to Room Restrictions.
	restriction := models.RoomRestriction{
		StarDate:      reservation.StarDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: reservation.Id,
		RestrictionId: 1,
	}
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot insert room restriction.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// store the reservation to session.
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// send EMail Notification - first to guest
	layout := "2006-01-02"
	htmlMessageToGuest := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear %s:, <br>
		This is to confirm your reservatioin from %s to %s.
	
	`, reservation.FirstName, reservation.StarDate.Format(layout), reservation.EndDate.Format(layout))
	msg := models.MailData{
		To:       reservation.Email,
		From:     "admin@booking.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessageToGuest,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	htmlMessageToOwner := fmt.Sprintf(`
		<strong>New Reservation</strong><br>
		Dear %s: <br><br>
		New reservation details<br><br>
		Name  : %s&nbsp;%s<br>
		Email : %s<br>
		Phone : %s<br>
		StartDate : %s<br>
		EndDate : %s<br>
		RoomName : %s<br>
	`, "Owner",
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		reservation.StarDate.Format(layout),
		reservation.EndDate.Format(layout),
		reservation.Room.RoomName,
	)

	msg = models.MailData{
		To:   "owner@bookings.com",
		From: "admin@bookings.com",
		Subject: fmt.Sprintf("New Reservation from %s to %s for %s",
			reservation.StarDate.Format(layout),
			reservation.EndDate.Format(layout),
			reservation.Room.RoomName,
		),
		Content: htmlMessageToOwner,
	}

	m.App.MailChan <- msg

	// redirect to other page.
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals is the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

// Majors is the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

// Availability is the search avialability  page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

// Availability is the search avialability  page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	var rooms []models.Room
	var data = make(map[string]interface{})

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "failed to parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err = m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//m.App.InfoLog.Println("Rooms : ", rooms)
	// no availability of rooms
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	data["rooms"] = rooms

	res := models.Reservation{
		StarDate: startDate,
		EndDate:  endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
	//fmt.Fprintf(w, "startDate is %s and end date is %s", start, end)
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    string `json:room_id`
	StartDate string `json:start_date`
	EndDate   string `json:end_date`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}
		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)

	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error 1",
		}
		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error 2",
		}
		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))

	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error 3",
		}
		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomId)
	log.Println(err)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error 4",
		}
		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    strconv.Itoa(roomId),
	}

	out, _ := json.MarshalIndent(resp, "", "   ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Availability is the search avialability  page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// ReservationSummary displays reservation summary
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("cannot get reservation summary from session")
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	// removing data from the session.
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	// converting startdate and enddate to stringMap
	layout := "2006-01-02"
	sd := reservation.StarDate.Format(layout)
	ed := reservation.EndDate.Format(layout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays the user to choose room.
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	expected := strings.Split(r.RequestURI, "/")
	roomId, err := strconv.Atoi(expected[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", errors.New("could not fetch reservation from session "))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.RoomId = roomId

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL parameters, builds a sessional variable and takes url to make-reservation page.
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	// get the params from URL id, e, s
	roomId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomById(roomId)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// insert into reservations model.
	var reservation models.Reservation
	reservation.RoomId = roomId
	reservation.Room.RoomName = room.RoomName
	reservation.StarDate = startDate
	reservation.EndDate = endDate

	// Make a new session reservation.
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		// Todo- take user back to page.
	}
	id, _, err := m.DB.Authenticate(r.Form.Get("email"), r.Form.Get("password"))

	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
