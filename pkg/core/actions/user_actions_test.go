package actions

import (
	"context"
	"fmt"
	"server/pkg/core/entities"
	db "server/pkg/core/storage"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	countriesTableName = "countries"
	cititesTableName   = "cities"
	tokensTableName    = "tokens"
	usersTableName     = "users"
)

type UserActionsTestSuite struct {
	suite.Suite
	action StorageAction
	ctx    context.Context
}

func (s *UserActionsTestSuite) SetupSuite() {
	s.ctx = context.Background()
	zerolog.SetGlobalLevel(zerolog.Disabled)

	storage := db.NewStorage(db.StorageConfig{
		DSN: "host=localhost port=5432 user=testserver password=passw0rd dbname=testserverdb sslmode=disable",
	})
	s.action = NewStorageAction(storage)
	err := s.action.Storage.Connect()
	s.Require().NoError(err)
}

func (s *UserActionsTestSuite) SetupTest() {
	s.Require().NoError(s.truncateAll())
}

func (s *UserActionsTestSuite) TestCreate() {
	actual, err := s.action.CreateUser("test", "password", "description")
	s.Require().NoError(err)

	expected := entities.User{
		ID:          1,
		Login:       "test",
		Description: "description",
		CityID:      0,
	}
	assert.Equal(s.T(), expected, entities.User{
		ID: actual.ID, Login: actual.Login, Description: actual.Description, CityID: actual.CityID,
	})
}

func (s *UserActionsTestSuite) TestGet() {
	testUsersRows := []struct {
		login       string
		password    string
		description string
	}{
		{login: "test1", password: "password1", description: "desc1"},
		{login: "test2", password: "password2", description: "desc2"},
		{login: "test3", password: "password3", description: "desc3"},
	}

	for _, row := range testUsersRows {
		_, err := s.action.CreateUser(row.login, row.password, row.description)
		s.Require().NoError(err)
	}

	actual, err := s.action.GetUser("test3")
	s.Require().NoError(err)

	expected := entities.User{
		ID: 3, Login: "test3", Description: "desc3", CityID: 0,
	}

	assert.Equal(s.T(), expected, entities.User{
		ID: actual.ID, Login: actual.Login, Description: actual.Description, CityID: actual.CityID,
	})
}

func (s *UserActionsTestSuite) TearDownSuite() {
	s.action.Storage.Disconnect()
}

func TestMapperTestSuite(t *testing.T) {
	suite.Run(t, new(UserActionsTestSuite))
}

func (s *UserActionsTestSuite) truncateAll() error {
	tableNames := []string{countriesTableName, cititesTableName, tokensTableName, usersTableName}

	for _, tableName := range tableNames {
		query := squirrel.Expr(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName))
		rows, err := s.action.Storage.RunQuery(s.ctx, query)
		if err != nil {
			return err
		}
		if err := rows.Err(); err != nil {
			return err
		}
		rows.Close()
	}

	return nil
}
