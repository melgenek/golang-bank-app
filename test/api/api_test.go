package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang_bank_demo/src/api"
	"golang_bank_demo/src/dto"
	"golang_bank_demo/src/errors"
	"golang_bank_demo/src/model"
	"golang_bank_demo/src/service"
	test_service "golang_bank_demo/test/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AccountApiSuite struct {
	suite.Suite
	service *test_service.StubAccountService
	api     *mux.Router
}

func TestAccountServiceSuite(t *testing.T) {
	suite.Run(t, new(AccountApiSuite))
}

func (suite *AccountApiSuite) SetupTest() {
	suite.service = new(test_service.StubAccountService)
	authApi := api.NewAuthenticatedApi(service.NewStubAuthenticationService())
	suite.api = api.NewAccountApi(suite.service, authApi).Router()
}

func (suite *AccountApiSuite) TestShouldGetAccount() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(20)}
	suite.service.On("Get", accountId).Return(account, nil)
	req, _ := http.NewRequest("GET", "/accounts/1", nil)
	req.Header.Set("Authorization", "Bearer token_user_1")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusOK)
	assert.Equal(suite.T(), resp.Body.String(), "{\"id\":1,\"balance\":\"20\"}\n")
	suite.service.AssertExpectations(suite.T())
}

func (suite *AccountApiSuite) TestShouldNotGetAccountWhenNoToken() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(20)}
	suite.service.On("Get", accountId).Return(account, nil)
	req, _ := http.NewRequest("GET", "/accounts/1", nil)
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusUnauthorized)
	assert.Equal(suite.T(), resp.Body.String(), "{\"message\":\"Unauthorized\"}\n")
	suite.service.AssertNotCalled(suite.T(), "Get", accountId)
}

func (suite *AccountApiSuite) TestShouldNotGetAccountWhenUnknownUser() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(20)}
	suite.service.On("Get", accountId).Return(account, nil)
	req, _ := http.NewRequest("GET", "/accounts/1", nil)
	req.Header.Set("Authorization", "Bearer unknown_user")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusForbidden)
	assert.Equal(suite.T(), resp.Body.String(), "{\"message\":\"The user cannot access the api\"}\n")
	suite.service.AssertNotCalled(suite.T(), "Get", accountId)
}

func (suite *AccountApiSuite) TestShouldNotGetAccountWhenAccountAccessForbidden() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	suite.service.On("Get", accountId).Return(nil, &errors.ForbiddenAccountAccessError{AccountId: accountId, UserId: userId})
	req, _ := http.NewRequest("GET", "/accounts/1", nil)
	req.Header.Set("Authorization", "Bearer token_user_1")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusForbidden)
	assert.Equal(suite.T(), resp.Body.String(), "{\"message\":\"The user 1 cannot access the account 1\"}\n")
	suite.service.AssertExpectations(suite.T())
}

func (suite *AccountApiSuite) TestShouldCreateAccount() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	account := &model.Account{Id: accountId, Owner: userId, Balance: decimal.NewFromInt(20)}
	suite.service.On("Create", userId).Return(account, nil)
	req, _ := http.NewRequest("POST", "/accounts", nil)
	req.Header.Set("Authorization", "Bearer token_user_1")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusCreated)
	assert.Equal(suite.T(), resp.Body.String(), "{\"id\":1,\"balance\":\"20\"}\n")
	suite.service.AssertExpectations(suite.T())
}

func (suite *AccountApiSuite) TestShouldNotCreateWhenDuplicateAccount() {
	userId := model.UserId(1)
	suite.service.On("Create", userId).Return(nil, &errors.DuplicateAccountError{UserId: userId})
	req, _ := http.NewRequest("POST", "/accounts", nil)
	req.Header.Set("Authorization", "Bearer token_user_1")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusConflict)
	assert.Equal(suite.T(), resp.Body.String(), "{\"message\":\"The user 1 already has an account\"}\n")
	suite.service.AssertExpectations(suite.T())
}

func (suite *AccountApiSuite) TestShouldTopUp() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	request := &dto.TopUpRequest{Id: accountId, Amount: decimal.NewFromInt(100)}
	suite.service.On("TopUp", request, userId).Return(nil)
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/top-up", bytes.NewReader(body))

	req.Header.Set("Authorization", "Bearer token_user_1")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusOK)
	assert.Equal(suite.T(), resp.Body.String(), "\"{}\"\n")
	suite.service.AssertExpectations(suite.T())
}

func (suite *AccountApiSuite) TestShouldNotTopUpWhenValidationError() {
	userId := model.UserId(1)
	accountId := model.AccountId(1)
	request := &dto.TopUpRequest{Id: accountId, Amount: decimal.NewFromInt(100)}
	suite.service.On("TopUp", request, userId).Return(&errors.ValidationError{Field: "id", Message: "The id has to be positive"})
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/top-up", bytes.NewReader(body))

	req.Header.Set("Authorization", "Bearer token_user_1")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusBadRequest)
	assert.Equal(suite.T(), resp.Body.String(), "{\"message\":\"Invalid field 'id': The id has to be positive\"}\n")
	suite.service.AssertExpectations(suite.T())
}

func (suite *AccountApiSuite) TestShouldTransfer() {
	userId := model.UserId(1)
	fromAccountId := model.AccountId(1)
	toAccountId := model.AccountId(1)
	request := &dto.TransferRequest{From: fromAccountId, To: toAccountId, Amount: decimal.NewFromInt(100)}
	suite.service.On("Transfer", request, userId).Return(nil)
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/transfer", bytes.NewReader(body))

	req.Header.Set("Authorization", "Bearer token_user_1")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusOK)
	assert.Equal(suite.T(), resp.Body.String(), "\"{}\"\n")
	suite.service.AssertExpectations(suite.T())
}

func (suite *AccountApiSuite) TestShouldTransferWhenBalanceTooLow() {
	userId := model.UserId(1)
	fromAccountId := model.AccountId(1)
	toAccountId := model.AccountId(1)
	request := &dto.TransferRequest{From: fromAccountId, To: toAccountId, Amount: decimal.NewFromInt(100)}
	suite.service.On("Transfer", request, userId).Return(&errors.BalanceTooLowError{AccountId: fromAccountId})
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/transfer", bytes.NewReader(body))

	req.Header.Set("Authorization", "Bearer token_user_1")
	resp := httptest.NewRecorder()

	suite.api.ServeHTTP(resp, req)

	assert.Equal(suite.T(), resp.Code, http.StatusBadRequest)
	assert.Equal(suite.T(), resp.Body.String(), "{\"message\":\"The account 1 does not have enough money\"}\n")
	suite.service.AssertExpectations(suite.T())
}
