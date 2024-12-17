package database

import (
	"database/sql"
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/entity"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
}

func (s *AccountDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)

	s.db = db

	db.Exec("CREATE TABLE accounts (id varchar(255), balance float, created_at timestamp default current_timestamp, updated_at timestamp default current_timestamp on update current_timestamp)")

	s.accountDB = NewAccountDB(db)
}

func (s *AccountDBTestSuite) TearDownTest() {
	defer s.db.Close()

	s.db.Exec("DROP TABLE accounts")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSave() {
	account := entity.NewAccount(uuid.NewString(), 0.0)

	err := s.accountDB.Save(account)
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestFindByID() {
	account := entity.NewAccount(uuid.NewString(), 0.0)

	err := s.accountDB.Save(account)
	s.Nil(err)

	accountFound, err := s.accountDB.FindByID(account.ID)
	s.Nil(err)
	s.Equal(account.ID, accountFound.ID)
	s.Equal(account.Balance, accountFound.Balance)
}
