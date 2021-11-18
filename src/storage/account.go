package storage

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"golang_bank_demo/src/errors"
	"golang_bank_demo/src/model"
)

type AccountStorage interface {
	Create(owner model.UserId) (*model.Account, error)
	Get(accountId model.AccountId) (*model.Account, error)
	TopUp(accountId model.AccountId, amount decimal.Decimal) error
	Transfer(from, to model.AccountId, amount decimal.Decimal) error
}

const uniqueConstraintErrorCode = pq.ErrorCode("23505")

type PostgresAccountStorage struct {
	db *sqlx.DB
}

func NewPostgresAccountStorage(db *sqlx.DB) AccountStorage {
	return &PostgresAccountStorage{db}
}

func (storage *PostgresAccountStorage) Create(owner model.UserId) (*model.Account, error) {
	newId := model.AccountId(-1)
	if err := storage.db.Get(&newId, "INSERT INTO accounts (owner_id) VALUES ($1) RETURNING id", owner); err == nil {
		return &model.Account{Id: newId, Owner: owner, Balance: decimal.NewFromInt(0)}, nil
	} else if pgErr := err.(*pq.Error); pgErr.Code == uniqueConstraintErrorCode {
		return nil, &errors.DuplicateAccountError{UserId: owner}
	} else {
		return nil, &errors.InternalServerError{Err: pgErr}
	}
}

func (storage *PostgresAccountStorage) Get(accountId model.AccountId) (account *model.Account, err error) {
	storage.executeInTransaction(func(tx *sqlx.Tx) error {
		account, err = storage.get(tx, accountId)
		return err
	})
	return
}

func (storage *PostgresAccountStorage) get(tx *sqlx.Tx, accountId model.AccountId) (*model.Account, error) {
	account := &model.Account{}
	if err := tx.Get(account, "SELECT * FROM accounts WHERE id=$1", accountId); err == nil {
		return account, nil
	} else if err == sql.ErrNoRows {
		return nil, &errors.AccountDoesNotExistError{AccountId: accountId}
	} else {
		return nil, &errors.InternalServerError{Err: err}
	}
}

func (storage *PostgresAccountStorage) TopUp(accountId model.AccountId, amount decimal.Decimal) error {
	if result, err := storage.db.Exec("UPDATE accounts SET balance = balance + $2 WHERE id = $1", accountId, amount); err != nil {
		return &errors.InternalServerError{Err: err}
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return &errors.InternalServerError{Err: err}
	} else if rowsAffected == 0 {
		return &errors.AccountDoesNotExistError{AccountId: accountId}
	} else {
		return nil
	}
}

func (storage *PostgresAccountStorage) Transfer(from, to model.AccountId, amount decimal.Decimal) error {
	return storage.executeInTransaction(func(tx *sqlx.Tx) error {
		if _, err := storage.get(tx, from); err != nil {
			return err
		} else if decreaseResult, err := tx.Exec("UPDATE accounts SET balance = balance - $2 WHERE id = $1 AND balance >= $2", from, amount); err != nil {
			return &errors.InternalServerError{Err: err}
		} else if decreaseBalanceRows, err := decreaseResult.RowsAffected(); err != nil {
			return &errors.InternalServerError{Err: err}
		} else if decreaseBalanceRows == 0 {
			return &errors.BalanceTooLowError{AccountId: from}
		} else if increaseResult, err := tx.Exec("UPDATE accounts SET balance = balance + $2 WHERE id = $1", to, amount); err != nil {
			return &errors.InternalServerError{Err: err}
		} else if increaseBalanceRows, err := increaseResult.RowsAffected(); err != nil {
			return &errors.InternalServerError{Err: err}
		} else if increaseBalanceRows == 0 {
			return &errors.AccountDoesNotExistError{AccountId: to}
		} else {
			return nil
		}
	})
}

func (storage *PostgresAccountStorage) executeInTransaction(f func(*sqlx.Tx) error) error {
	if tx, err := storage.db.Beginx(); err != nil {
		return err
	} else if err := f(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return &errors.InternalServerError{Err: rollbackErr}
		} else {
			return err
		}
	} else if commitErr := tx.Commit(); commitErr != nil {
		return &errors.InternalServerError{Err: commitErr}
	} else {
		return nil
	}
}
