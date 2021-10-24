package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pradeep-veera89/webApplication/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	// GET Methods
	// {"home", "/", "GET", []postData{}, http.StatusOK},
	// {"about", "/about", "GET", []postData{}, http.StatusOK},
	// {"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	// {"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	// {"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	// {"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// {"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},

	// // POST Methods
	// {"post-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "john@smith.com"},
	// 	{key: "phone", value: "555-555-555"},
	// }, http.StatusOK},
	// {"post-search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-02-12"},
	// 	{key: "start", value: "2020-02-12"},
	// }, http.StatusOK},
	// {"post-search-availability", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-02-12"},
	// 	{key: "start", value: "2020-02-12"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, test := range theTests {
		testingUrl := ts.URL + test.url
		if test.method == "GET" {
			resp, err := ts.Client().Get(testingUrl)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		}
		if test.method == "POST" {
			values := url.Values{}
			for _, x := range test.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(testingUrl, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {

	reservation := models.Reservation{
		RoomId: 1
		Room:models.Room{
			Id : 1,
			RoomName: "Generals Quarters",
		},
	}

	req, _  := http.NewRequest("GET", "/make-reservation",nil)
	ctx := getCtx(req)
	req =  req.WithContext(ctx)

	rr := httptest.NewRecorder() 
	session.Put(ctx,"reservation", reservation)
	
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wront response code got %d wanted %d", rr.Code, http.StatusOK)
	}
}

func getCtx(req * http.Request) context.Context{
	ctx, err := session.Load(req.Context(),req.Header.Get("X-Session"))
	if err !=  nil {
		log.Println(err)
	}
	return ctx
}
