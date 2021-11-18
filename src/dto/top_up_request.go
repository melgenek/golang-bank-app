package dto

import (
	"github.com/shopspring/decimal"
	"golang_bank_demo/src/errors"
	"golang_bank_demo/src/model"
)

type TopUpRequest struct {
	Id     model.AccountId `json:"id"`
	Amount decimal.Decimal `json:"amount"`
}

func (request *TopUpRequest) Validate() error {
	if request.Id <= 0 {
		return errors.NewValidationError("id", "The id has to be positive")
	} else if request.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return errors.NewValidationError("amount", "The amount has to be positive")
	} else {
		return nil
	}
}
