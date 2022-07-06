package storage

import (
	"context"
	"fmt"
	"server/pkg/core/entities"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
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
	s.Require().NoError(err)
}

func (s *BaseStorageSuite) SetupTest() {
	s.Require().NoError(s.truncateAll())
}

func (s *BaseStorageSuite) TearDownSuite() {
	s.storage.Disconnect()
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
