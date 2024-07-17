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
