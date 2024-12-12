package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/database"
	getaccountbalance "github.com/MatheusNP/fc-ms-wallet/ms-account/internal/usecase/get_account_balance"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/web"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/web/webserver"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"user",
		"pass",
		"mysqlaccount",
		"3307",
		"account",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Erro ao conectar ao banco de dados: %v \n", err)
	}

	accountDB := database.NewAccountDB(db)

	getAccountBalanceUseCase := getaccountbalance.NewGetAccountBalanceUseCase(accountDB)

	webserver := webserver.NewWebServer(":3003")

	accountHandler := web.NewWebAccountHandler(*getAccountBalanceUseCase)

	webserver.AddHandler("/accounts", accountHandler.GetAccountBalance)

	webserver.AddHandler("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	webserver.Start()
}
