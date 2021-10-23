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

// SearchAvailabilityByDates returns true if Room is available in given date, else return false.
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {

	var numRows int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query searches within the range, with start and end date.
	stmt := `
		select count(id) 
		from room_restrictions
		where 
			$1 < end_date and 
			$2 >start_date and 
			room_id = $3;`

	row := m.DB.QueryRowContext(ctx, stmt, start, end, roomId)
	err := row.Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range.
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms = []models.Room{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	
	select 
		r.id, r.room_name
	from 
		rooms r
	where
		r.id not in 
		(select room_id from room_restrictions rr where $1 <rr.end_date and $2 >rr.start_date)
	`
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.Id,
			&room.RoomName,
		)
		if err != nil {
			return rooms, nil
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}
