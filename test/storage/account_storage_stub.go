package storage

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"golang_bank_demo/src/model"
)

type StubAccountStorage struct {
	mock.Mock
}

func (storage *StubAccountStorage) Create(owner model.UserId) (*model.Account, error) {
	args := storage.Called(owner)
	if account, ok := args.Get(0).(*model.Account); ok {
		return account, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (storage *StubAccountStorage) Get(accountId model.AccountId) (*model.Account, error) {
	args := storage.Called(accountId)
	if account, ok := args.Get(0).(*model.Account); ok {
		return account, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (storage *StubAccountStorage) TopUp(accountId model.AccountId, amount decimal.Decimal) error {
	args := storage.Called(accountId, amount)
	return args.Error(0)
}

func (storage *StubAccountStorage) Transfer(from, to model.AccountId, amount decimal.Decimal) error {
	args := storage.Called(from, to, amount)
	return args.Error(0)
}
