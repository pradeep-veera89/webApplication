package dbrepo

import (
	"context"
	"time"

	"github.com/pradeep-veera89/webApplication/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts reservations into DB.
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	var newId int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservations(first_name, last_name, email, phone, 
			start_date, end_date, room_id,created_at, updated_at)
			values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StarDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

// InsertRoomRestriction inserts room_restrictions into database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions(start_date, end_date, room_id, reservation_id, 
		restriction_id, created_at, updated_at)
		values($1,$2,$3,$4,$5,$6,$7)`
	_, err := m.DB.ExecContext(ctx, stmt,
		res.StarDate,
		res.EndDate,
		res.RoomId,
		res.ReservationId,
		res.RestrictionId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
