package models

import "time"

// Reservation holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// Users is the user model
type Users struct {
	Id         int
	FirstName  string
	LastName   string
	Email      string
	Password   string
	AssesLevel int
	CreatedAt  time.Time
	UpdateAt   time.Time
}

// Rooms is the room model
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdateAt  time.Time
}

// Restrictions is the resctrictions model
type Restrictions struct {
	Id              int
	RestrictionName string
	CreatedAt       time.Time
	UpdateAt        time.Time
}

// Reservations is the reservation model
type Reservations struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StarDate  time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdateAt  time.Time
	Room      Rooms
}

type RoomRestrictions struct {
	Id            int
	StarDate      time.Time
	EndDate       time.Time
	RoomId        int
	ReservationId int
	RestrictionId int
	CreatedAt     time.Time
	UpdateAt      time.Time

	Room        Rooms
	Reservation Reservations
	Restriction Restrictions
}
