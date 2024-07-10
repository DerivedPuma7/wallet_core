package entity

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient("any name", "any email")

	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "any name", client.Name)
	assert.Equal(t, "any email", client.Email)
}

func TestCreateNewClientWithInvalidArgs(t *testing.T) {
	client, err := NewClient("", "")

	assert.Nil(t, client)
	assert.NotNil(t, err)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("any name", "any email")

	err := client.Update("new name", "new email")

	assert.Nil(t, err)
	assert.Equal(t, "new name", client.Name)
	assert.Equal(t, "new email", client.Email)
}

func TestUpdateClientWithInvalidArgs(t *testing.T) {
	client, _ := NewClient("any name", "any email")

	err := client.Update("", "new email")

	assert.NotNil(t, err)
}

func TestAddAccountToClient(t *testing.T) {
	client, _ := NewClient("any name", "any email")
	account, _ := NewAccount(client)

	err := client.AddAccount(account)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(client.Accounts))
}

func TestAddAccountToInvalidClient(t *testing.T) {
	client, _ := NewClient("client 1", "email 1")
	accountOwner, _ := NewClient("any name", "any email")
	account, _ := NewAccount(accountOwner)

	err := client.AddAccount(account)

	assert.NotNil(t, err)
	assert.Equal(t, "account does not belong to client", err.Error())
}
