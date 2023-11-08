package main

import (
	"database/sql"
	"fmt"

	"github.com/fontinelle/fc-ms-wallet/internal/database"
	"github.com/fontinelle/fc-ms-wallet/internal/event"
	"github.com/fontinelle/fc-ms-wallet/internal/usecase/create_account"
	"github.com/fontinelle/fc-ms-wallet/internal/usecase/create_client"
	"github.com/fontinelle/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com/fontinelle/fc-ms-wallet/internal/web"
	"github.com/fontinelle/fc-ms-wallet/internal/web/web_server"
	"github.com/fontinelle/fc-ms-wallet/pkg/events"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@/wallet?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	transactionDb := database.NewTransactionDB(db)

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(accountDb, transactionDb, eventDispatcher, transactionCreatedEvent)

	web_server := web_server.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	web_server.AddHandler("/clients", clientHandler.CreateClient)
	web_server.AddHandler("/accounts", accountHandler.CreateAccount)
	web_server.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	web_server.Start()
}
