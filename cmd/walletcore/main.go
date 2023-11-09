package main

import (
	"context"
	"database/sql"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/fontinelle/fc-ms-wallet/internal/database"
	"github.com/fontinelle/fc-ms-wallet/internal/event"
	"github.com/fontinelle/fc-ms-wallet/internal/event/handler"
	"github.com/fontinelle/fc-ms-wallet/internal/usecase/create_account"
	"github.com/fontinelle/fc-ms-wallet/internal/usecase/create_client"
	"github.com/fontinelle/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com/fontinelle/fc-ms-wallet/internal/web"
	"github.com/fontinelle/fc-ms-wallet/internal/web/web_server"
	"github.com/fontinelle/fc-ms-wallet/pkg/events"
	"github.com/fontinelle/fc-ms-wallet/pkg/kafka"
	"github.com/fontinelle/fc-ms-wallet/pkg/uow"
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
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)
	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)

	web_server := web_server.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	web_server.AddHandler("/clients", clientHandler.CreateClient)
	web_server.AddHandler("/accounts", accountHandler.CreateAccount)
	web_server.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	web_server.Start()
}
