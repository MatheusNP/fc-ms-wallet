package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MatheusNP/fc-ms-wallet/internal/database"
	"github.com/MatheusNP/fc-ms-wallet/internal/event"
	createaccount "github.com/MatheusNP/fc-ms-wallet/internal/usecase/create_account"
	createclient "github.com/MatheusNP/fc-ms-wallet/internal/usecase/create_client"
	createtransaction "github.com/MatheusNP/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com/MatheusNP/fc-ms-wallet/internal/web"
	"github.com/MatheusNP/fc-ms-wallet/internal/web/webserver"
	"github.com/MatheusNP/fc-ms-wallet/pkg/events"
	"github.com/MatheusNP/fc-ms-wallet/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"root",
		"root",
		"localhost",
		"3306",
		"wallet",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)
	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := createclient.NewCreateClientUseCase(clientDB)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(clientDB, accountDB)
	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent)

	webserver := webserver.NewWebServer(
		":3000",
	)

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()
}
