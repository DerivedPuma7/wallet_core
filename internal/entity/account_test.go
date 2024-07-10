package entity

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewAccount(t *testing.T) {
	client, _ := NewClient("any name", "any email")
	account, err := NewAccount(client)

	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, client, account.Client)
	assert.Equal(t, 0.0, account.Balance)
}

func TestCreateNewAccountWithInvalidArgs(t *testing.T) {
	account, err := NewAccount(nil)

	assert.Nil(t, account)
	assert.NotNil(t, err)
	assert.Equal(t, "client must exist", err.Error())
}

func TestCreditAccount(t *testing.T) {
	client, _ := NewClient("any name", "any email")
	account, _ := NewAccount(client)

	err := account.Credit(10)

	assert.Nil(t, err)
	assert.Equal(t, 10.0, account.Balance)

	err = account.Credit(5)

	assert.Nil(t, err)
	assert.Equal(t, 15.0, account.Balance)
}

func TestCreditAccountWithInvalidValue(t *testing.T) {
	client, _ := NewClient("any name", "any email")
	account, _ := NewAccount(client)

	err := account.Credit(0)

	assert.NotNil(t, err)
	assert.Equal(t, "credit must be a value bigger than zero", err.Error())

	err = account.Credit(-1)
	assert.NotNil(t, err)
	assert.Equal(t, "credit must be a value bigger than zero", err.Error())
}

func TestDebitAccount(t *testing.T) {
	client, _ := NewClient("any name", "any email")
	account, _ := NewAccount(client)

	account.Credit(10)
	err := account.Debit(8)
	
	assert.Nil(t, err)
	assert.Equal(t, 2.0, account.Balance)
}

func TestDebitAccountWithInvalidValue(t *testing.T) {
	client, _ := NewClient("any name", "any email")
	account, _ := NewAccount(client)

	err := account.Debit(-8)

	assert.NotNil(t, err)
	assert.Equal(t, "debit must be a value bigger than zero", err.Error())
}

func TestDebitAccountWithInsufficientBalance(t *testing.T) {
	client, _ := NewClient("any name", "any email")
	account, _ := NewAccount(client)

	account.Credit(5)
	err := account.Debit(8)

	assert.NotNil(t, err)
	assert.Equal(t, "insufficient balance", err.Error())
}
