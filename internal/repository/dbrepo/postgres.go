package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/pradeep-veera89/webApplication/internal/models"
	"golang.org/x/crypto/bcrypt"
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

// GetRoomById gets a room by Id
func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {

	room := models.Room{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select id, room_name, created_at, updated_at from rooms where id=$1`
	row := m.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(
		&room.Id,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdateAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil

}

func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
			from users 
			where id=$1`
	row := m.DB.QueryRowContext(ctx, stmt, id)
	var user models.User
	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AssesLevel,
		&user.CreatedAt,
		&user.UpdateAt,
	)

	if err != nil {
		return models.User{}, err
	}

	if row.Err() != nil {
		return models.User{}, err
	}

	return user, nil
}

// UpdateUser updates user to DB.
func (m *postgresDBRepo) UpdateUser(u models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update 
				user 
			set
				first_name =$1
				last_name=$2
				email=$3
				access_level=$4
				updated_at=$5
			where 
				id = $8`
	_, err := m.DB.ExecContext(ctx, stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AssesLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

//Authenticate authenticates the user email with password.
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id , password from users where email =$1", email)
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		select 
			r.id,
			r.first_name,
			r.last_name, 
		 	r.email, 
		 	r.phone, 
		 	r.start_date,
		 	r.end_date,
			r.processed,
		 	rm.id,
		 	rm.room_name
		 from
		 	reservations r
		 left join rooms rm on r.room_id = rm.id
		 order by r.start_date asc`

	var reservations []models.Reservation
	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.Id,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StarDate,
			&i.EndDate,
			&i.Processed,
			&i.Room.Id,
			&i.Room.RoomName,
		)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		select 
			r.id,
			r.first_name,
			r.last_name, 
		 	r.email, 
		 	r.phone, 
		 	r.start_date,
		 	r.end_date,
			r.processed,
		 	rm.id,
		 	rm.room_name
		 from
		 	reservations r
		    left join rooms rm on r.room_id = rm.id
		 where r.processed = 0
		    order by r.start_date asc`

	var reservations []models.Reservation
	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.Id,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StarDate,
			&i.EndDate,
			&i.Processed,
			&i.Room.Id,
			&i.Room.RoomName,
		)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}
