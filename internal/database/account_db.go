package database

import (
	"database/sql"

	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	"github.com.br/derivedpuma7/wallet-core/internal/gateway"
)

type AccountDb struct {
	DB *sql.DB
}

var _ gateway.AccountGateway = (*AccountDb)(nil)

func NewAccountDb(db *sql.DB) *AccountDb {
	return &AccountDb{
		DB: db,
	}
}

func (a *AccountDb) FindById(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client
	stmt, err := a.DB.Prepare(`
		SELECT 
			a.id, a.clientId, a.balance, a.createdAt, a.updatedAt,
			c.id, c.name, c.email, c.createdAt, c.updatedAt
		FROM accounts as a 
		INNER JOIN clients as c 
			ON a.clientId = c.id 
		WHERE a.id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(
		&account.ID,
		&account.Client.ID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
		&client.ID,
		&client.Name,
		&client.Email,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *AccountDb) Save(account *entity.Account) error {
	stmt, err := a.DB.Prepare(`
		INSERT INTO accounts(id, clientId, balance, createdAt, updatedAt)
		VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(account.ID, account.Client.ID, account.Balance, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountDb) UpdateBalance(account *entity.Account) error {
	stmt, err := a.DB.Prepare("UPDATE accounts SET balance = ? WHERE id = ?")
  if err != nil {
    return err
  }
  defer stmt.Close()
  _, err = stmt.Exec(account.Balance, account.ID)
  if err != nil {
    return err
  }
  return nil
}
