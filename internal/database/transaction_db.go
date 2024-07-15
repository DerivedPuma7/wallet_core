package database

import (
	"database/sql"

	"github.com.br/derivedpuma7/wallet-core/internal/entity"
)

type TransactionDb struct {
	DB *sql.DB
}

func NewTransactionDb(db *sql.DB) *TransactionDb {
	return &TransactionDb{
		DB: db,
	}
}

func (t *TransactionDb) Create(transaction *entity.Transaction) error {
	stmt, err := t.DB.Prepare(`
		INSERT INTO transactions(id, accountIdFrom, accountIdTo, amount, createdAt)
		VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(transaction.ID, transaction.AccountFrom.ID, transaction.AccountTo.ID, transaction.Amount, transaction.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
