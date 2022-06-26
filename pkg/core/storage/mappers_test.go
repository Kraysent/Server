package storage

import (
	"server/pkg/core/entities"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MapperTestSuite struct {
	suite.Suite
	storage *Storage
}

func (s *MapperTestSuite) SetupTest() {
	config := StorageConfig{
		DSN: "host=localhost port=5432 user=testserver password=passw0rd dbname=testserverdb sslmode=disable",
	}
	zerolog.SetGlobalLevel(zerolog.Disabled)

	s.storage = NewStorage(config)
	s.storage.Connect()

	_, _ = s.storage.CreateUser(entities.User{
		Login:        "test1",
		Salt:         1,
		PasswordHash: "weak_hash",
		Description:  "weak desc",
	})
	_, _ = s.storage.CreateUser(entities.User{
		Login:        "test2",
		Salt:         2,
		PasswordHash: "middle_hash",
		Description:  "middle desc",
	})
	_, _ = s.storage.CreateUser(entities.User{
		Login:        "test3",
		Salt:         3,
		PasswordHash: "strong_hash",
		Description:  "strong_desc",
	})
}

func (s *MapperTestSuite) TestGet() {
	actual, _ := s.storage.GetUser("test2")
	assert.Equal(s.T(), entities.User{
		ID:           2,
		Login:        "test2",
		Salt:         2,
		PasswordHash: "middle_hash",
		Description:  "middle desc",
	}, *actual)
}

func (s *MapperTestSuite) TearDownSuite() {
	s.storage.Disconnect()
}

func TestMapperTestSuite(t *testing.T) {
	suite.Run(t, new(MapperTestSuite))
}
