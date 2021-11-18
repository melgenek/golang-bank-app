package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"golang_bank_demo/src/config"
	"golang_bank_demo/src/postgres"
	"testing"
)

type PostgresTestSuite struct {
	Db *sqlx.DB
}

func (suite *PostgresTestSuite) SetupSuite(t *testing.T) {
	appConfig := config.AppConfig{Postgres: config.Postgres{"postgres://demo:demo@localhost/bank?sslmode=disable"}}
	db, err := postgres.CreateClient(&appConfig)
	assert.NoError(t, err)
	suite.Db = db
}

func (suite *PostgresTestSuite) SetupTest(t *testing.T) {
	err := postgres.SetUp(suite.Db)
	assert.NoError(t, err)
}

func (suite *PostgresTestSuite) TearDownTest(t *testing.T) {
	err := postgres.TearDown(suite.Db)
	assert.NoError(t, err)
}
