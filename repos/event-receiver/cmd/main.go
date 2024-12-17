package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/database"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/event"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/event/handler"
	getaccountbalance "github.com/MatheusNP/fc-ms-wallet/ms-account/internal/usecase/get_account_balance"
	updateaccountbalance "github.com/MatheusNP/fc-ms-wallet/ms-account/internal/usecase/update_account_balance"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/web"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/web/webserver"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/pkg/events"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/pkg/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"user",
		"pass",
		"mysqlaccount",
		"3306",
		"account",
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
	if err := m.Up(); err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("Erro ao conectar ao banco de dados: %v \n", err)
	}

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "account",
	}
	kafkaConsumer := kafka.NewConsumer(
		&configMap,
		[]string{"balances"},
	)

	accountDB := database.NewAccountDB(db)

	getAccountBalanceUseCase := getaccountbalance.NewGetAccountBalanceUseCase(accountDB)
	updateAccountBalanceUseCase := updateaccountbalance.NewUpdateAccountBalanceUseCase(accountDB)

	msgChan := make(chan *ckafka.Message)

	go func() {
		err := kafkaConsumer.Consume(msgChan)
		if err != nil {
			fmt.Printf("Erro ao consumir mensagens: %v \n", err)
		}
	}()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(updateAccountBalanceUseCase))
	balanceUpdatedEvent := event.NewBalanceUpdated()

	go func() {
		for msg := range msgChan {
			var eventMsg event.BalanceUpdated

			err := json.Unmarshal(msg.Value, &eventMsg)
			if err != nil {
				fmt.Printf("Erro ao decodificar JSON: %v \n", err)
				continue
			}

			switch eventMsg.GetName() {
			case "BalanceUpdated":
				balanceUpdatedEvent.SetPayload(eventMsg.Payload)
				eventDispatcher.Dispatch(balanceUpdatedEvent)
			}

			fmt.Printf("Mensagem recebida: %s\n", string(msg.Value))
		}
	}()

	webserver := webserver.NewWebServer(":3003")

	accountHandler := web.NewWebAccountHandler(*getAccountBalanceUseCase)

	webserver.AddHandler("/accounts", accountHandler.GetAccountBalance)

	webserver.AddHandler("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	webserver.Start()
}
