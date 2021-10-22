package dbrepo

import (
	"database/sql"

	"github.com/pradeep-veera89/webApplication/internal/config"
	"github.com/pradeep-veera89/webApplication/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo initializes postgressRepo
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
