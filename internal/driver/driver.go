package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database conenction pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDB = 5
const maxDBLifeTime = 5 * time.Minute

// ConnectSQL create database pool for postgres
func ConnectSQL(dsn string) (*DB, error) {

	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	// Setting some parameters to DB pool
	d.SetConnMaxIdleTime(maxIdleDB)
	d.SetConnMaxLifetime(maxDBLifeTime)
	d.SetMaxOpenConns(maxOpenDbConn)

	dbConn.SQL = d

	err = testDB(dbConn.SQL)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// testDB tries to ping database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates new database for application.
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
