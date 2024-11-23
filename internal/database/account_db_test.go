package database

import (
	"database/sql"
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
	client    *entity.Client
}

func (s *AccountDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)

	s.db = db

	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date default current_timestamp)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date default current_timestamp)")

	s.accountDB = NewAccountDB(db)

	s.client, _ = entity.NewClient("name", "email")

}

func (s *AccountDBTestSuite) TearDownTest() {
	defer s.db.Close()

	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSave() {
	account := entity.NewAccount(s.client)

	err := s.accountDB.Save(account)
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestFindByID() {
	s.db.Exec(
		"INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)",
		s.client.ID,
		s.client.Name,
		s.client.Email,
		s.client.CreatedAt,
	)

	account := entity.NewAccount(s.client)

	err := s.accountDB.Save(account)
	s.Nil(err)

	accountFound, err := s.accountDB.FindByID(account.ID)
	s.Nil(err)
	s.Equal(account.ID, accountFound.ID)
	s.Equal(account.Client.ID, accountFound.Client.ID)
	s.Equal(account.Balance, accountFound.Balance)
	s.Equal(account.Client.ID, accountFound.Client.ID)
	s.Equal(account.Client.Name, accountFound.Client.Name)
	s.Equal(account.Client.Email, accountFound.Client.Email)
}
