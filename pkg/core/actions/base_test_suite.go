package actions

import (
	"context"
	"fmt"

	"server/pkg/db"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	countriesTableName = "countries"
	citiesTableName    = "cities"
	tokensTableName    = "tokens"
	usersTableName     = "users"
)

type BaseActionsTestSuite struct {
	suite.Suite
	action StorageAction
	ctx    context.Context
}

func (s *BaseActionsTestSuite) SetupSuite() {
	s.ctx = context.Background()
	zerolog.SetGlobalLevel(zerolog.Disabled)

	storage := db.NewStorage(db.StorageConfig{
		DSN: "host=localhost port=5432 user=testserver password=passw0rd dbname=testserverdb sslmode=disable",
	})
	s.action = NewStorageAction(storage)
	err := s.action.Storage.Connect()
	require.NoError(s.T(), err)
}

func (s *BaseActionsTestSuite) SetupTest() {
	require.NoError(s.T(), s.truncateAll())
}

func (s *BaseActionsTestSuite) TearDownSuite() {
	require.NoError(s.T(), s.action.Storage.Disconnect())
}

func (s *BaseActionsTestSuite) truncateAll() error {
	tableNames := []string{countriesTableName, citiesTableName, tokensTableName, usersTableName}

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
