package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MatheusNP/fc-ms-wallet/internal/database"
	"github.com/MatheusNP/fc-ms-wallet/internal/event"
	"github.com/MatheusNP/fc-ms-wallet/internal/event/handler"
	createaccount "github.com/MatheusNP/fc-ms-wallet/internal/usecase/create_account"
	createclient "github.com/MatheusNP/fc-ms-wallet/internal/usecase/create_client"
	createtransaction "github.com/MatheusNP/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com/MatheusNP/fc-ms-wallet/internal/web"
	"github.com/MatheusNP/fc-ms-wallet/internal/web/webserver"
	"github.com/MatheusNP/fc-ms-wallet/pkg/events"
	"github.com/MatheusNP/fc-ms-wallet/pkg/kafka"
	"github.com/MatheusNP/fc-ms-wallet/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		"user",
		"pass",
		"mysql",
		"3306",
		"wallet",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		panic(err)
	}
	m.Up()

	if err := db.Ping(); err != nil {
		fmt.Printf("Erro ao conectar ao banco de dados: %v \n", err)
	}

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
	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(
		uow,
		eventDispatcher,
		transactionCreatedEvent,
		balanceUpdatedEvent,
	)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Printf("Rodando na porta %s \n", webserver.WebServerPort)

	webserver.Start()
}
