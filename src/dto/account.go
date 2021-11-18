package dto

import (
	"github.com/shopspring/decimal"
	"golang_bank_demo/src/model"
)

type Account struct {
	Id      model.AccountId `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}

func AccountFromModel(account *model.Account) *Account {
	return &Account{Id: account.Id, Balance: account.Balance}
}
