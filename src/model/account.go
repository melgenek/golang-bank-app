package model

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	Id      AccountId       `db:"id"`
	Owner   UserId          `db:"owner_id"`
	Balance decimal.Decimal `db:"balance"`
}
