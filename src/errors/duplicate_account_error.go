package errors

import (
	"fmt"
	"golang_bank_demo/src/model"
)

type DuplicateAccountError struct {
	UserId model.UserId
}

func (err *DuplicateAccountError) Error() string {
	return fmt.Sprintf("The user %d already has an account", err.UserId)
}

func (err *DuplicateAccountError) Is(target error) bool {
	t, ok := target.(*DuplicateAccountError)
	if ok {
		return t.UserId == err.UserId
	} else {
		return false
	}
}
