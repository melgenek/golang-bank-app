package service

import (
	"errors"
	"golang_bank_demo/src/model"
)

type AuthenticationService interface {
	GetUser(token string) (*model.UserId, error)
}

type StubAuthenticationService struct {
}

func NewStubAuthenticationService() AuthenticationService {
	return &StubAuthenticationService{}
}

var authError = errors.New("The user cannot access the api")

func (service *StubAuthenticationService) GetUser(token string) (*model.UserId, error) {
	if token == "token_user_1" {
		userId := model.UserId(1)
		return &userId, nil
	} else if token == "token_user_2" {
		userId := model.UserId(2)
		return &userId, nil
	} else {
		return nil, authError
	}
}
