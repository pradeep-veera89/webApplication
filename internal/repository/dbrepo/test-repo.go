package dbrepo

import (
	"errors"
	"time"

	"github.com/pradeep-veera89/webApplication/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts reservations into DB.
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	if res.RoomId == 1 || res.RoomId == 2 {
		return 1, nil
	}
	if res.RoomId == 3 {
		return 2, nil
	}
	return 0, errors.New("invalid room Id")

}

// InsertRoomRestriction inserts room_restrictions into database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {

	if res.ReservationId == 2 {
		return errors.New("invalid room Id")
	}
	return nil
}

// SearchAvailabilityByDates returns true if Room is available in given date, else return false.
func (m *testDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {

	if roomId == 100 {
		return false, errors.New("invalid room id")
	}
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range.
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms = []models.Room{
		{Id: 2, RoomName: "testName"},
	}

	td := "2020-03-02"
	layout := "2006-01-02"
	testDate, _ := time.Parse(layout, td)
	if start == testDate {
		return nil, errors.New("query failed")
	}

	td = "2020-03-01"
	testDate, _ = time.Parse(layout, td)
	if start == testDate {
		return rooms, nil
	}
	return []models.Room{}, nil
}

// GetRoomById gets a room by Id
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {

	room := models.Room{}
	if id > 2 {
		return room, errors.New("Some error")
	}

	return room, nil
}

func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	return models.User{}, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 2, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation
	return reservations, nil
}
