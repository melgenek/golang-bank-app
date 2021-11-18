package service

import (
	"golang_bank_demo/src/dto"
	"golang_bank_demo/src/errors"
	"golang_bank_demo/src/model"
	"golang_bank_demo/src/storage"
)

type AccountService interface {
	Create(user model.UserId) (*model.Account, error)
	Get(accountId model.AccountId, user model.UserId) (*model.Account, error)
	TopUp(request *dto.TopUpRequest, user model.UserId) error
	Transfer(request *dto.TransferRequest, user model.UserId) error
}

type RealAccountService struct {
	storage storage.AccountStorage
}

func NewAccountService(accountStorage storage.AccountStorage) AccountService {
	return &RealAccountService{storage: accountStorage}
}

func (service *RealAccountService) Create(user model.UserId) (*model.Account, error) {
	return service.storage.Create(user)
}

func (service *RealAccountService) Get(accountId model.AccountId, user model.UserId) (*model.Account, error) {
	if account, err := service.storage.Get(accountId); err != nil {
		return nil, err
	} else if account.Owner != user {
		return nil, &errors.ForbiddenAccountAccessError{AccountId: accountId, UserId: user}
	} else {
		return account, nil
	}
}

func (service *RealAccountService) TopUp(request *dto.TopUpRequest, user model.UserId) error {
	if err := request.Validate(); err != nil {
		return err
	} else if account, err := service.storage.Get(request.Id); err != nil {
		return err
	} else if account.Owner != user {
		return &errors.ForbiddenAccountAccessError{AccountId: request.Id, UserId: user}
	} else {
		return service.storage.TopUp(request.Id, request.Amount)
	}
}

func (service *RealAccountService) Transfer(request *dto.TransferRequest, user model.UserId) error {
	if err := request.Validate(); err != nil {
		return err
	} else if fromAccount, err := service.storage.Get(request.From); err != nil {
		return err
	} else if fromAccount.Owner != user {
		return &errors.ForbiddenAccountAccessError{AccountId: request.From, UserId: user}
	} else if err = service.storage.Transfer(request.From, request.To, request.Amount); err != nil {
		return err
	} else {
		return nil
	}
}
