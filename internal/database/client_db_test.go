package database

import (
	"database/sql"
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	clientDB *ClientDB
}

func (s *ClientDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)

	s.db = db

	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date default current_timestamp)")

	s.clientDB = NewClientDB(db)
}

func (s *ClientDBTestSuite) TearDownTest() {
	defer s.db.Close()

	s.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (s *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("name", "email")

	err := s.clientDB.Save(client)
	s.Nil(err)

	clientFound, err := s.clientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, clientFound.ID)
	s.Equal(client.Name, clientFound.Name)
	s.Equal(client.Email, clientFound.Email)
}

func (s *ClientDBTestSuite) TestSave() {
	client := &entity.Client{
		ID:    "123",
		Name:  "name",
		Email: "email",
	}

	err := s.clientDB.Save(client)
	s.Nil(err)
}
