package database

import (
	"database/sql"
	"testing"

	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDbTestSuite struct {
	suite.Suite
	db *sql.DB
	accountDb *AccountDb
	client *entity.Client
}

func (s *AccountDbTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id string, name string, email string, createdAt date, updatedAt date)")
	db.Exec("CREATE TABLE accounts(id string, clientId string, balance float, createdAt date, updatedAt date)")
	s.accountDb = NewAccountDb(db)
	s.client, _ = entity.NewClient("any name", "any email")
}

func(s *AccountDbTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func TestAccountDbTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDbTestSuite))
}

func (s *AccountDbTestSuite) TestSave() {
	account, _ := entity.NewAccount(s.client)
	err := s.accountDb.Save(account)

	s.Nil(err)
}

func (s *AccountDbTestSuite) TestFindById() {
	s.db.Exec(`
		INSERT INTO clients(id, name, email, createdAt, updatedAt)
		VALUES(?, ?, ?, ?, ?)
	`, s.client.ID, s.client.Name, s.client.Email, s.client.CreatedAt, s.client.UpdatedAt,
	)
	account, _ := entity.NewAccount(s.client)
	err := s.accountDb.Save(account)
	s.Nil(err)

	accountDb, err := s.accountDb.FindById(account.ID)

	s.Nil(err)
	s.Equal(account.ID, accountDb.ID)
	s.Equal(account.Client.ID, accountDb.Client.ID)
	s.Equal(account.Balance, accountDb.Balance)
}
