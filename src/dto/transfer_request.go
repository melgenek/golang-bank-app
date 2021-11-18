package dto

import (
	"github.com/shopspring/decimal"
	"golang_bank_demo/src/errors"
	"golang_bank_demo/src/model"
)

type TransferRequest struct {
	From   model.AccountId `json:"from"`
	To     model.AccountId `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}

func (request *TransferRequest) Validate() error {
	if request.From <= 0 {
		return errors.NewValidationError("from", "The id has to be positive")
	} else if request.To <= 0 {
		return errors.NewValidationError("to", "The id has to be positive")
	} else if request.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return errors.NewValidationError("amount", "The amount has to be positive")
	} else {
		return nil
	}
}
