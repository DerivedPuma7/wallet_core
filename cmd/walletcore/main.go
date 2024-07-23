package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com.br/derivedpuma7/wallet-core/internal/database"
	"github.com.br/derivedpuma7/wallet-core/internal/event"
	"github.com.br/derivedpuma7/wallet-core/internal/event/handler"
	createaccount "github.com.br/derivedpuma7/wallet-core/internal/usecase/create_account"
	createclient "github.com.br/derivedpuma7/wallet-core/internal/usecase/create_client"
	createtransaction "github.com.br/derivedpuma7/wallet-core/internal/usecase/create_transaction"
	"github.com.br/derivedpuma7/wallet-core/internal/web"
	"github.com.br/derivedpuma7/wallet-core/internal/web/webserver"
	"github.com.br/derivedpuma7/wallet-core/pkg/events"
	"github.com.br/derivedpuma7/wallet-core/pkg/kafka"
	"github.com.br/derivedpuma7/wallet-core/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
  if err != nil {
    panic(err)
  }
  defer db.Close()
  startDatabase(db)

  configMap := ckafka.ConfigMap{
    "bootstrap.servers": "kafka:29092",
    "group.id": "wallet",
  }
  kafkaProducer := kafka.NewKafkaProducer(&configMap)

  eventDispatcher := events.NewEventDispatcher()
  transactionCreatedEvent := event.NewTransactionCreated(time.Now())
  balanceUpdatedEvent := event.NewBalanceUpdated()
  eventDispatcher.Register(transactionCreatedEvent.GetName(), handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
  eventDispatcher.Register(balanceUpdatedEvent.GetName(), handler.NewBalanceUpdatedKafkaHandler(kafkaProducer))

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
  createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

  httpPort := ":8080"

  webserver := webserver.NewWebServer(httpPort)
  clientHandler := web.NewWebClientHandler(*createClientUseCase)
  accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
  transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

  webserver.AddHandler("/clients", clientHandler.CreateClient)
  webserver.AddHandler("/accounts", accountHandler.CreateAccount)
  webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

  fmt.Println("server is running on: http://localhost:", httpPort)
  webserver.Start()
}

func startDatabase(db *sql.DB) {
  sqls := []string{
    "CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), createdAt date, updatedAt date);",
    "CREATE TABLE IF NOT EXISTS accounts(id varchar(255), clientId varchar(255), balance float, createdAt date, updatedAt date);",
    "CREATE TABLE IF NOT EXISTS transactions(id varchar(255), accountIdFrom varchar(255), accountIdTo varchar(255), amount float, createdAt date);",
    "TRUNCATE TABLE clients",
    "TRUNCATE TABLE accounts",
    "TRUNCATE TABLE transactions",
    "INSERT INTO clients (id, name, email, createdAt, updatedAt) VALUES ('005c9780-ddd8-4834-87f2-2fc98620892e', 'client 1', 'client1@test.com', current_timestamp, current_timestamp)",
    "INSERT INTO clients (id, name, email, createdAt, updatedAt) VALUES ('17374f16-1ac1-4f86-9b6a-1726da0c5327', 'client 2', 'client2@test.com', current_timestamp, current_timestamp)",
    "INSERT INTO accounts (id, clientId, balance, createdAt, updatedAt) VALUES('fab64237-6e11-4152-80ac-2cdd8a5a86dd', '005c9780-ddd8-4834-87f2-2fc98620892e', 10000, current_timestamp, current_timestamp)",
    "INSERT INTO accounts (id, clientId, balance, createdAt, updatedAt) VALUES('e67b0ed1-7b13-45d7-ad20-1c9682842635', '17374f16-1ac1-4f86-9b6a-1726da0c5327', 10000, current_timestamp, current_timestamp)",
  }

  for _, sql := range sqls {
    _, err := db.Exec(sql)
    if err != nil {
      panic(err)
    }
  }
}
