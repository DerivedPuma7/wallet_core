package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com.br/derivedpuma7/wallet-core/internal/database"
	"github.com.br/derivedpuma7/wallet-core/internal/event"
	createaccount "github.com.br/derivedpuma7/wallet-core/internal/usecase/create_account"
	createclient "github.com.br/derivedpuma7/wallet-core/internal/usecase/create_client"
	createtransaction "github.com.br/derivedpuma7/wallet-core/internal/usecase/create_transaction"
	"github.com.br/derivedpuma7/wallet-core/internal/web"
	"github.com.br/derivedpuma7/wallet-core/internal/web/webserver"
	"github.com.br/derivedpuma7/wallet-core/pkg/events"
	"github.com.br/derivedpuma7/wallet-core/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "localhost", "3306", "wallet"))
  if err != nil {
    panic(err)
  }
  defer db.Close()

  eventDispatcher := events.NewEventDispatcher()
  transactionCreatedEvent := event.NewTransactionCreated(time.Now())
  // eventDispatcher.Register("TransactionCreated", handler)

  clientDb := database.NewClientDb(db)
  accountDb := database.NewAccountDb(db)
  
  ctx := context.Background()
  uow := uow.NewUow(ctx, db)
  uow.Register("AccountDb", func(tx *sql.Tx) interface{} {
    return database.NewAccountDb(db)
  })
  uow.Register("TransactionDb", func(tx *sql.Tx) interface{} {
    return database.NewTransactionDb(db)
  })

  createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
  createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
  createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent)

  webserver := webserver.NewWebServer(":3000")
  clientHandler := web.NewWebClientHandler(*createClientUseCase)
  accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
  transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

  webserver.AddHandler("/clients", clientHandler.CreateClient)
  webserver.AddHandler("/accounts", accountHandler.CreateAccount)
  webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

  webserver.Start()
}
