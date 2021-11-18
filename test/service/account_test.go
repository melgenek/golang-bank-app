package service

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang_bank_demo/src/dto"
	"golang_bank_demo/src/errors"
	"golang_bank_demo/src/model"
	"golang_bank_demo/src/service"
	"golang_bank_demo/test/storage"
	"testing"
)

type AccountServiceSuite struct {
	suite.Suite
	storage *storage.StubAccountStorage
	service service.AccountService
}

func TestAccountServiceSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceSuite))
}

func (suite *AccountServiceSuite) SetupTest() {
	suite.storage = new(storage.StubAccountStorage)
	suite.service = service.NewAccountService(suite.storage)
}

func (suite *AccountServiceSuite) TestShouldCreateAnAccount() {
	userId := model.UserId(1)
	account := &model.Account{Id: 1, Owner: userId, Balance: decimal.NewFromInt(20)}
	suite.storage.On("Create", userId).Return(account, nil)

	createdAccount, err := suite.service.Create(userId)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdAccount, account)
	suite.storage.AssertExpectations(suite.T())
}

func (suite *AccountServiceSuite) TestShouldGetAnAccount() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(20)}
	suite.storage.On("Get", accountId).Return(account, nil)

	foundAccount, err := suite.service.Get(accountId, userId)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), foundAccount, account)
	suite.storage.AssertExpectations(suite.T())
}

func (suite *AccountServiceSuite) TestShouldNotGetAnAccountWhenDifferentUser() {
	accountId := model.AccountId(1)
	anotherUserId := model.UserId(2)
	account := &model.Account{Id: accountId, Owner: model.UserId(1), Balance: decimal.NewFromInt(20)}
	suite.storage.On("Get", accountId).Return(account, nil)

	foundAccount, err := suite.service.Get(accountId, anotherUserId)

	assert.ErrorIs(suite.T(), err, &errors.ForbiddenAccountAccessError{AccountId: accountId, UserId: anotherUserId})
	assert.Nil(suite.T(), foundAccount)
	suite.storage.AssertExpectations(suite.T())
}

func (suite *AccountServiceSuite) TestShouldTopUpAnAccount() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	amount := decimal.NewFromInt(20)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", accountId).Return(account, nil)
	suite.storage.On("TopUp", accountId, amount).Return(nil)

	err := suite.service.TopUp(&dto.TopUpRequest{Id: accountId, Amount: amount}, userId)

	assert.NoError(suite.T(), err)
	suite.storage.AssertExpectations(suite.T())
}

func (suite *AccountServiceSuite) TestShouldNotTopUpAnAccountWhenDifferentUser() {
	userId := model.UserId(1)
	anotherUserId := model.UserId(2)
	accountId := model.AccountId(1)
	amount := decimal.NewFromInt(20)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", accountId).Return(account, nil)
	suite.storage.On("TopUp", accountId, amount).Return(nil)

	err := suite.service.TopUp(&dto.TopUpRequest{Id: accountId, Amount: amount}, anotherUserId)

	assert.ErrorIs(suite.T(), err, &errors.ForbiddenAccountAccessError{AccountId: accountId, UserId: anotherUserId})
	suite.storage.AssertNotCalled(suite.T(), "TopUp", accountId, amount)
}

func (suite *AccountServiceSuite) TestShouldTopUpAnAccountWhenTheBalanceIsNotPositive() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	amount := decimal.NewFromInt(-20)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", accountId).Return(account, nil)
	suite.storage.On("TopUp", accountId, amount).Return(nil)

	err := suite.service.TopUp(&dto.TopUpRequest{Id: accountId, Amount: amount}, userId)

	assert.ErrorIs(suite.T(), err, &errors.ValidationError{Field: "amount", Message: "The amount has to be positive"})
	suite.storage.AssertNotCalled(suite.T(), "TopUp", accountId, amount)
}

func (suite *AccountServiceSuite) TestShouldTopUpAnAccountWhenTheIdIsNotPositive() {
	userId := model.UserId(1)
	accountId := model.AccountId(-1)
	amount := decimal.NewFromInt(20)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", accountId).Return(account, nil)
	suite.storage.On("TopUp", accountId, amount).Return(nil)

	err := suite.service.TopUp(&dto.TopUpRequest{Id: accountId, Amount: amount}, userId)

	assert.ErrorIs(suite.T(), err, &errors.ValidationError{Field: "id", Message: "The id has to be positive"})
	suite.storage.AssertNotCalled(suite.T(), "TopUp", accountId, amount)
}

func (suite *AccountServiceSuite) TestShouldTransfer() {
	userId := model.UserId(1)
	fromAccountId := model.AccountId(1)
	toAccountId := model.AccountId(2)
	amount := decimal.NewFromInt(20)
	account := &model.Account{Id: fromAccountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", fromAccountId).Return(account, nil)
	suite.storage.On("Transfer", fromAccountId, toAccountId, amount).Return(nil)

	err := suite.service.Transfer(&dto.TransferRequest{From: fromAccountId, To: toAccountId, Amount: amount}, userId)

	assert.NoError(suite.T(), err)
	suite.storage.AssertExpectations(suite.T())
}

func (suite *AccountServiceSuite) TestShouldNotTransferWhenAccountOfDifferentUser() {
	userId := model.UserId(1)
	anotherUserId := model.UserId(2)
	fromAccountId := model.AccountId(1)
	toAccountId := model.AccountId(2)
	amount := decimal.NewFromInt(20)
	account := &model.Account{Id: fromAccountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", fromAccountId).Return(account, nil)
	suite.storage.On("Transfer", fromAccountId, toAccountId, amount).Return(nil)

	err := suite.service.Transfer(&dto.TransferRequest{From: fromAccountId, To: toAccountId, Amount: amount}, anotherUserId)

	assert.ErrorIs(suite.T(), err, &errors.ForbiddenAccountAccessError{AccountId: fromAccountId, UserId: anotherUserId})
	suite.storage.AssertNotCalled(suite.T(), "Transfer", fromAccountId, toAccountId, amount)
}

func (suite *AccountServiceSuite) TestShouldNotTransferWhenFromIdIsNotPositive() {
	userId := model.UserId(1)
	fromAccountId := model.AccountId(-1)
	toAccountId := model.AccountId(2)
	amount := decimal.NewFromInt(20)
	account := &model.Account{Id: fromAccountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", fromAccountId).Return(account, nil)
	suite.storage.On("Transfer", fromAccountId, toAccountId, amount).Return(nil)

	err := suite.service.Transfer(&dto.TransferRequest{From: fromAccountId, To: toAccountId, Amount: amount}, userId)

	assert.ErrorIs(suite.T(), err, &errors.ValidationError{Field: "from", Message: "The id has to be positive"})
	suite.storage.AssertNotCalled(suite.T(), "Transfer", fromAccountId, toAccountId, amount)
}

func (suite *AccountServiceSuite) TestShouldNotTransferWhenToIdIsNotPositive() {
	userId := model.UserId(1)
	fromAccountId := model.AccountId(1)
	toAccountId := model.AccountId(-2)
	amount := decimal.NewFromInt(20)
	account := &model.Account{Id: fromAccountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", fromAccountId).Return(account, nil)
	suite.storage.On("Transfer", fromAccountId, toAccountId, amount).Return(nil)

	err := suite.service.Transfer(&dto.TransferRequest{From: fromAccountId, To: toAccountId, Amount: amount}, userId)

	assert.ErrorIs(suite.T(), err, &errors.ValidationError{Field: "to", Message: "The id has to be positive"})
	suite.storage.AssertNotCalled(suite.T(), "Transfer", fromAccountId, toAccountId, amount)
}

func (suite *AccountServiceSuite) TestShouldNotTransferWhenAmountIsNotPositive() {
	userId := model.UserId(1)
	fromAccountId := model.AccountId(1)
	toAccountId := model.AccountId(2)
	amount := decimal.NewFromInt(-20)
	account := &model.Account{Id: fromAccountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", fromAccountId).Return(account, nil)
	suite.storage.On("Transfer", fromAccountId, toAccountId, amount).Return(nil)

	err := suite.service.Transfer(&dto.TransferRequest{From: fromAccountId, To: toAccountId, Amount: amount}, userId)

	assert.ErrorIs(suite.T(), err, &errors.ValidationError{Field: "amount", Message: "The amount has to be positive"})
	suite.storage.AssertNotCalled(suite.T(), "Transfer", fromAccountId, toAccountId, amount)
}

func (suite *AccountServiceSuite) TestShouldNotTransferWhenStorageErrors() {
	userId := model.UserId(1)
	fromAccountId := model.AccountId(1)
	toAccountId := model.AccountId(2)
	amount := decimal.NewFromInt(20)
	account := &model.Account{Id: fromAccountId, Owner: userId, Balance: decimal.NewFromInt(10)}
	suite.storage.On("Get", fromAccountId).Return(account, nil)
	suite.storage.On("Transfer", fromAccountId, toAccountId, amount).Return(&errors.BalanceTooLowError{AccountId: fromAccountId})

	err := suite.service.Transfer(&dto.TransferRequest{From: fromAccountId, To: toAccountId, Amount: amount}, userId)

	assert.ErrorIs(suite.T(), err, &errors.BalanceTooLowError{AccountId: fromAccountId})
	suite.storage.AssertExpectations(suite.T())
}
