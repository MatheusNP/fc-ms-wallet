package database

import (
	"database/sql"
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	transactionDB *TransactionDB
	client1       *entity.Client
	client2       *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
}

func (s *TransactionDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)

	s.db = db

	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date default current_timestamp)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date default current_timestamp)")
	db.Exec("CREATE TABLE transactions (id varchar(255), account_from_id varchar(255), account_to_id varchar(255), amount float, created_at date, updated_at date default current_timestamp)")

	s.transactionDB = NewTransactionDB(db)

	s.client1, _ = entity.NewClient("John", "john@mail.com")
	s.client2, _ = entity.NewClient("Joe", "joe@mail.com")

	accountFrom := entity.NewAccount(s.client1)
	accountFrom.Balance = 100.0
	s.accountFrom = accountFrom
	accountTo := entity.NewAccount(s.client2)
	accountTo.Balance = 100.0
	s.accountTo = accountTo

}

func (s *TransactionDBTestSuite) TearDownTest() {
	defer s.db.Close()

	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 10.0)

	s.Nil(err)
	s.Nil(s.transactionDB.Create(transaction))
}
