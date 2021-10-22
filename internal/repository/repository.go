package repository

import "github.com/pradeep-veera89/webApplication/internal/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) error
	AllUsers() bool
}
