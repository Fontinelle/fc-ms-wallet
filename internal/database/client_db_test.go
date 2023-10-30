package database

import (
	"database/sql"
	"testing"

	"github.com/fontinelle/fc-ms-wallet/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	clientDb *ClientDB
}

func (suite *ClientDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.Nil(err)
	suite.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	suite.clientDb = NewClientDB(db)
}

func (suite *ClientDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (suite *ClientDBTestSuite) TestSave() {
	client := &entity.Client{
		ID:    "1",
		Name:  "John Doe",
		Email: "j@j",
	}
	err := suite.clientDb.Save(client)
	suite.Nil(err)
}

func (suite *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("John Doe", "j@j")
	suite.clientDb.Save(client)

	clientDb, err := suite.clientDb.Get(client.ID)
	suite.Nil(err)
	suite.Equal(client.ID, clientDb.ID)
	suite.Equal(client.Name, clientDb.Name)
	suite.Equal(client.Email, clientDb.Email)
}
