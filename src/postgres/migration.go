package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
)

var migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		{
			Id: "1",
			Up: []string{"CREATE TABLE accounts (" +
				"id BIGSERIAL PRIMARY KEY," +
				"owner_id BIGINT NOT NULL," +
				"balance DECIMAL NOT NULL DEFAULT 0," +
				"UNIQUE(owner_id)" +
				")"},
			Down: []string{"DROP TABLE accounts"},
		},
	},
}

func SetUp(db *sqlx.DB) error {
	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	return err
}

func TearDown(db *sqlx.DB) error {
	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Down)
	return err
}
