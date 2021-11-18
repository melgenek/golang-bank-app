package errors

import (
	"fmt"
	"golang_bank_demo/src/model"
)

type AccountDoesNotExistError struct {
	AccountId model.AccountId
}

func (err *AccountDoesNotExistError) Error() string {
	return fmt.Sprintf("The account %d does not exist", err.AccountId)
}

func (err *AccountDoesNotExistError) Is(target error) bool {
	t, ok := target.(*AccountDoesNotExistError)
	if ok {
		return t.AccountId == err.AccountId
	} else {
		return false
	}
}
