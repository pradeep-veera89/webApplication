package dbrepo

import (
	"time"

	"github.com/pradeep-veera89/webApplication/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts reservations into DB.
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	var newId int

	return newId, nil
}

// InsertRoomRestriction inserts room_restrictions into database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {

	return nil
}

// SearchAvailabilityByDates returns true if Room is available in given date, else return false.
func (m *testDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range.
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms = []models.Room{}
	return rooms, nil
}

// GetRoomById gets a room by Id
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {

	room := models.Room{}
	return room, nil
}
