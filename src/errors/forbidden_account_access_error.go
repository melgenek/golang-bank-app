package errors

import (
	"fmt"
	"golang_bank_demo/src/model"
)

type ForbiddenAccountAccessError struct {
	AccountId model.AccountId
	UserId    model.UserId
}

func (err *ForbiddenAccountAccessError) Error() string {
	return fmt.Sprintf("The user %d cannot access the account %d", err.UserId, err.AccountId)
}

func (err *ForbiddenAccountAccessError) Is(target error) bool {
	t, ok := target.(*ForbiddenAccountAccessError)
	if ok {
		return t.AccountId == err.AccountId && t.UserId == err.UserId
	} else {
		return false
	}
}
