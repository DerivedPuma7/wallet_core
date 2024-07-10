package entity

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewTransaction(t *testing.T) {
	client1, _ := NewClient("any name 1", "any email 1")
	account1, _ := NewAccount(client1)
	account1.Credit(1000)
	client2, _ := NewClient("any name 2", "any email 2")
	account2, _ := NewAccount(client2)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 100)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, 900.0, account1.Balance)
	assert.Equal(t, 1100.0, account2.Balance)
}

func TestCreateNewTransactionWithInsufficientBalance(t *testing.T) {
	client1, _ := NewClient("any name 1", "any email 1")
	account1, _ := NewAccount(client1)
	account1.Credit(1000)
	client2, _ := NewClient("any name 2", "any email 2")
	account2, _ := NewAccount(client2)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 2000)

	assert.NotNil(t, err)
	assert.Error(t, err, "insufficient funds")
	assert.Nil(t, transaction)
	assert.Equal(t, 1000.0, account1.Balance)
	assert.Equal(t, 1000.0, account2.Balance)
}