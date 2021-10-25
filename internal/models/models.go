package models

import (
	"time"
)

// User is the user model
type User struct {
	Id         int
	FirstName  string
	LastName   string
	Email      string
	Password   string
	AssesLevel int
	CreatedAt  time.Time
	UpdateAt   time.Time
}

// Room is the room model
type Room struct {
	Id        int
	RoomName  string
	CreatedAt time.Time
	UpdateAt  time.Time
}

// Restriction is the resctrictions model
type Restriction struct {
	Id              int
	RestrictionName string
	CreatedAt       time.Time
	UpdateAt        time.Time
}

// Reservation is the reservation model
type Reservation struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StarDate  time.Time
	EndDate   time.Time
	RoomId    int
	CreatedAt time.Time
	UpdateAt  time.Time
	Room      Room
}

type RoomRestriction struct {
	Id            int
	StarDate      time.Time
	EndDate       time.Time
	RoomId        int
	ReservationId int
	RestrictionId int
	CreatedAt     time.Time
	UpdateAt      time.Time

	Room        Room
	Reservation Reservation
	Restriction Restriction
}

// MailData holds email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
