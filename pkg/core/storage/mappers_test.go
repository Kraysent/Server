package storage

import (
	"context"
	"fmt"
	"server/pkg/core/entities"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UsersMapperTestSuite struct {
	suite.Suite
	storage *Storage
	ctx     context.Context
}

func (s *UsersMapperTestSuite) SetupSuite() {
	s.ctx = context.Background()
	zerolog.SetGlobalLevel(zerolog.Disabled)

	config := StorageConfig{
		DSN: "host=localhost port=5432 user=testserver password=passw0rd dbname=testserverdb sslmode=disable",
	}

	s.storage = NewStorage(config)
	err := s.storage.Connect()
	s.Require().NoError(err)
}

func (s *UsersMapperTestSuite) SetupTest() {
	s.Require().NoError(s.truncateAll())
}

func (s *UsersMapperTestSuite) TestCreate() {
	description := "description"
	registration_time := time.Unix(1656363017, 0).UTC()
	actual, err := s.storage.CreateUser(UserCreateParams{
		Login:            "test",
		Salt:             1,
		PasswordHash:     "hash",
		Description:      &description,
		RegistrationDate: &registration_time,
	})
	s.Require().NoError(err)

	expected := entities.User{
		ID:               1,
		Login:            "test",
		Salt:             1,
		PasswordHash:     "hash",
		Description:      description,
		CityID:           0,
		RegistrationDate: registration_time,
	}
	assert.Equal(s.T(), &expected, actual)
}

func (s *UsersMapperTestSuite) TestFind() {
	testUsersRows := []struct {
		login            string
		salt             int
		hash             string
		description      string
		registrationDate time.Time
	}{
		{login: "test1", salt: 1, hash: "hash1", description: "desc1", registrationDate: time.Unix(1656362017, 0).UTC()},
		{login: "test2", salt: 2, hash: "hash2", description: "desc2", registrationDate: time.Unix(1656363017, 0).UTC()},
		{login: "test3", salt: 3, hash: "hash3", description: "desc3", registrationDate: time.Unix(1656364017, 0).UTC()},
	}

	for _, row := range testUsersRows {
		_, err := s.storage.CreateUser(UserCreateParams{
			Login: row.login, Salt: row.salt, PasswordHash: row.hash,
			Description: &row.description, RegistrationDate: &row.registrationDate,
		})
		s.Require().NoError(err)
	}

	login := "test2"
	actual, err := s.storage.FindUsers(UsersFindParams{
		Login: &login,
	})
	s.Require().NoError(err)

	expected := []entities.User{
		{
			ID: 2, Login: "test2", Salt: 2, PasswordHash: "hash2",
			Description: "desc2", CityID: 0, RegistrationDate: time.Unix(1656363017, 0).UTC(),
		},
	}
	assert.Equal(s.T(), expected, actual)
}

func (s *UsersMapperTestSuite) TearDownSuite() {
	s.storage.Disconnect()
}

func TestMapperTestSuite(t *testing.T) {
	suite.Run(t, new(UsersMapperTestSuite))
}

func (s *UsersMapperTestSuite) truncateAll() error {
	tableNames := []string{countriesTableName, cititesTableName, tokensTableName, usersTableName}

	for _, tableName := range tableNames {
		query := squirrel.Expr(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName))
		rows, err := s.storage.runQuery(s.ctx, query)
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
