package errors

import (
	"fmt"
	"golang_bank_demo/src/model"
)

type BalanceTooLowError struct {
	AccountId model.AccountId
}

func (err *BalanceTooLowError) Error() string {
	return fmt.Sprintf("The account %d does not have enough money", err.AccountId)
}

func (err *BalanceTooLowError) Is(target error) bool {
	t, ok := target.(*BalanceTooLowError)
	if ok {
		return t.AccountId == err.AccountId
	} else {
		return false
	}
}
