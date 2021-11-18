package storage

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang_bank_demo/src/errors"
	"golang_bank_demo/src/storage"
	"golang_bank_demo/test/postgres"
	"testing"
)

type AccountStorageSuite struct {
	suite.Suite
	postgres.PostgresTestSuite
	storage storage.AccountStorage
}

func TestAccountStorageSuite(t *testing.T) {
	suite.Run(t, new(AccountStorageSuite))
}

func (suite *AccountStorageSuite) SetupSuite() {
	suite.PostgresTestSuite.SetupSuite(suite.T())
	suite.storage = storage.NewPostgresAccountStorage(suite.Db)
}

func (suite *AccountStorageSuite) SetupTest() {
	suite.PostgresTestSuite.SetupTest(suite.T())
}

func (suite *AccountStorageSuite) TearDownTest() {
	suite.PostgresTestSuite.TearDownTest(suite.T())
}

func (suite *AccountStorageSuite) TestShouldCreateAndGetAnAccount() {
	createdAccount, err := suite.storage.Create(1)
	assert.NoError(suite.T(), err)

	foundAccount, err := suite.storage.Get(createdAccount.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdAccount, foundAccount)
}

func (suite *AccountStorageSuite) TestShouldNotCreateADuplicateAccount() {
	createdAccount, err := suite.storage.Create(1)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), createdAccount)

	duplicateAccount, err := suite.storage.Create(1)
	assert.ErrorIs(suite.T(), err, &errors.DuplicateAccountError{UserId: 1})
	assert.Nil(suite.T(), duplicateAccount)
}

func (suite *AccountStorageSuite) TestShouldNotGetAccountThatDoesNotExist() {
	foundAccount, err := suite.storage.Get(123)

	assert.ErrorIs(suite.T(), err, &errors.AccountDoesNotExistError{AccountId: 123})
	assert.Nil(suite.T(), foundAccount)
}

func (suite *AccountStorageSuite) TestShouldTopUpTheAccount() {
	createdAccount, err := suite.storage.Create(1)
	assert.NoError(suite.T(), err)

	err = suite.storage.TopUp(createdAccount.Id, decimal.NewFromInt(100))
	assert.NoError(suite.T(), err)
	err = suite.storage.TopUp(createdAccount.Id, decimal.NewFromInt(200))
	assert.NoError(suite.T(), err)

	foundAccount, err := suite.storage.Get(createdAccount.Id)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), foundAccount.Balance, decimal.NewFromInt(300))
}

func (suite *AccountStorageSuite) TestShouldNotTopUpTheAccountThatDoesNotExist() {
	err := suite.storage.TopUp(123, decimal.NewFromInt(100))

	assert.ErrorIs(suite.T(), err, &errors.AccountDoesNotExistError{AccountId: 123})
}

func (suite *AccountStorageSuite) TestShouldTransfer() {
	createdAccount1, err := suite.storage.Create(1)
	assert.NoError(suite.T(), err)

	createdAccount2, err := suite.storage.Create(2)
	assert.NoError(suite.T(), err)

	err = suite.storage.TopUp(createdAccount1.Id, decimal.NewFromInt(200))
	assert.NoError(suite.T(), err)

	err = suite.storage.Transfer(createdAccount1.Id, createdAccount2.Id, decimal.NewFromInt(200))
	assert.NoError(suite.T(), err)

	foundAccount1, err := suite.storage.Get(createdAccount1.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), foundAccount1.Balance, decimal.NewFromInt(0))

	foundAccount2, err := suite.storage.Get(createdAccount2.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), foundAccount2.Balance, decimal.NewFromInt(200))
}

func (suite *AccountStorageSuite) TestShouldNotTransferFromAccountThatDoesNotExist() {
	err := suite.storage.Transfer(1, 2, decimal.NewFromInt(100))

	assert.ErrorIs(suite.T(), err, &errors.AccountDoesNotExistError{AccountId: 1})
}

func (suite *AccountStorageSuite) TestShouldNotTransferToAccountThatDoesNotExist() {
	createdAccount1, err := suite.storage.Create(1)
	assert.NoError(suite.T(), err)
	err = suite.storage.TopUp(createdAccount1.Id, decimal.NewFromInt(300))
	assert.NoError(suite.T(), err)

	err = suite.storage.Transfer(createdAccount1.Id, 2, decimal.NewFromInt(100))

	assert.ErrorIs(suite.T(), err, &errors.AccountDoesNotExistError{AccountId: 2})
}

func (suite *AccountStorageSuite) TestShouldNotTransferWhenNtEnoughMoney() {
	createdAccount1, err := suite.storage.Create(1)
	assert.NoError(suite.T(), err)

	err = suite.storage.Transfer(createdAccount1.Id, 2, decimal.NewFromInt(100))

	assert.ErrorIs(suite.T(), err, &errors.BalanceTooLowError{AccountId: createdAccount1.Id})
}
