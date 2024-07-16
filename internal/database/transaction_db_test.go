package database

import (
	"database/sql"
	"testing"

	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type TransactionDbTestSuite struct {
	suite.Suite
	db * sql.DB
	client1 *entity.Client
	client2 *entity.Client
	accountFrom *entity.Account
	accountTo *entity.Account
	transactionDb *TransactionDb
}

func (s *TransactionDbTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), createdAt date, updatedAt date)")
	db.Exec("CREATE TABLE accounts(id varchar(255), clientId varchar(255), balance float, createdAt date, updatedAt date)")
	db.Exec("CREATE TABLE transactions(id varchar(255), accountIdFrom varchar(255), accountIdTo varchar(255), amount float, createdAt date)")
	s.SetupAccount()
	s.transactionDb = NewTransactionDb(db)
}

func (s *TransactionDbTestSuite) SetupAccount() {
	client1, err := entity.NewClient("any name 1", "any email 1")
	s.Nil(err)
	s.client1 = client1

	client2, err := entity.NewClient("any name 2", "any email 2")
	s.Nil(err)
	s.client2 = client2

	accountFrom, _ := entity.NewAccount(s.client1)
	accountFrom.Balance = 1000
	s.accountFrom = accountFrom

	accountTo, _ := entity.NewAccount(s.client2)
	accountFrom.Balance = 1000
	s.accountTo = accountTo
}

func(s *TransactionDbTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE transactions")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func TestTransactionDbTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDbTestSuite))
}

func (s *TransactionDbTestSuite) TestCreateTransaction() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100)
	s.Nil(err)

	err = s.transactionDb.Create(transaction)

	s.Nil(err)
}
