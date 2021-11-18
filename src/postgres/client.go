package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang_bank_demo/src/config"
)

func CreateClient(config *config.AppConfig) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", config.Postgres.Uri)
}
