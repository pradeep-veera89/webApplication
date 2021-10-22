package repository

import (
	"time"

	"github.com/pradeep-veera89/webApplication/internal/models"
)

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDates(start, end time.Time, roomId int) (bool, error)
	AllUsers() bool
}
