package storage

import (
	"context"
	"fmt"
	"server/pkg/core/entities"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BaseStorageSuite struct {
	suite.Suite
	storage *Storage
	ctx     context.Context
}

func (s *BaseStorageSuite) SetupSuite() {
	s.ctx = context.Background()
	zerolog.SetGlobalLevel(zerolog.Disabled)

	config := StorageConfig{
		DSN: "host=localhost port=5432 user=testserver password=passw0rd dbname=testserverdb sslmode=disable",
	}

	s.storage = NewStorage(config)
	err := s.storage.Connect()
	require.NoError(s.T(), err)
}

func (s *BaseStorageSuite) SetupTest() {
	require.NoError(s.T(), s.truncateAll())
}

func (s *BaseStorageSuite) TearDownSuite() {
	s.storage.Disconnect()
}

func (s *BaseStorageSuite) createTestUser(login string, salt int, hash, description string, registrationDate time.Time) {
	_, err := s.storage.CreateUser(UserCreateParams{
		Login: login, Salt: salt, PasswordHash: hash,
		Description: &description, RegistrationDate: &registrationDate,
	})
	require.NoError(s.T(), err)
}

func (s *BaseStorageSuite) createTestToken(userID int, value string, startDate, endDate time.Time) {
	_, err := s.storage.CreateToken(TokenCreateParams{
		UserID: userID, Value: value, StartDate: startDate, ExpirationDate: endDate,
	})
	require.NoError(s.T(), err)
}

func (s *BaseStorageSuite) truncateAll() error {
	tableNames := []string{
		countriesTableName, citiesTableName, entities.TableTokens, entities.TableUsers,
	}

	for _, tableName := range tableNames {
		query := squirrel.Expr(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName))
		rows, err := s.storage.RunQuery(s.ctx, query)
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
