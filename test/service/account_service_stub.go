package service

import (
	"github.com/stretchr/testify/mock"
	"golang_bank_demo/src/dto"
	"golang_bank_demo/src/model"
)

type StubAccountService struct {
	mock.Mock
}

func (service *StubAccountService) Create(user model.UserId) (*model.Account, error) {
	args := service.Called(user)
	if account, ok := args.Get(0).(*model.Account); ok {
		return account, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (service *StubAccountService) Get(accountId model.AccountId, user model.UserId) (*model.Account, error) {
	args := service.Called(accountId)
	if account, ok := args.Get(0).(*model.Account); ok {
		return account, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (service *StubAccountService) TopUp(request *dto.TopUpRequest, user model.UserId) error {
	args := service.Called(request, user)
	return args.Error(0)
}

func (service *StubAccountService) Transfer(request *dto.TransferRequest, user model.UserId) error {
	args := service.Called(request, user)
	return args.Error(0)
}
