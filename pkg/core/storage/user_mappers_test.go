package storage

import (
	"server/pkg/core/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UsersMapperTestSuite struct {
	BaseStorageSuite
}

func (s *UsersMapperTestSuite) TestCreate() {
	description := "description"
	registration_time := time.Unix(1656363017, 0).UTC()
	actual, err := s.storage.CreateUser(UserCreateParams{
		Login: "test", Salt: 1, PasswordHash: "hash",
		Description: &description, RegistrationDate: &registration_time,
	})
	require.NoError(s.T(), err)

	expected := entities.User{
		ID: 1, Login: "test", Salt: 1, PasswordHash: "hash",
		Description: description, CityID: 0, RegistrationDate: registration_time,
	}
	assert.Equal(s.T(), &expected, actual)
}

func (s *UsersMapperTestSuite) TestFind() {
	s.createTestUser("test1", 1, "hash1", "desc1", time.Unix(1656362017, 0).UTC())
	s.createTestUser("test2", 2, "hash2", "desc2", time.Unix(1656363017, 0).UTC())
	s.createTestUser("test3", 3, "hash3", "desc3", time.Unix(1656364017, 0).UTC())

	login := "test2"
	actual, err := s.storage.FindUsers(UsersFindParams{
		Login: &login,
	})
	require.NoError(s.T(), err)

	expected := []entities.User{
		{
			ID: 2, Login: "test2", Salt: 2, PasswordHash: "hash2",
			Description: "desc2", CityID: 0, RegistrationDate: time.Unix(1656363017, 0).UTC(),
		},
	}
	assert.Equal(s.T(), expected, actual)
}

func TestUserMapperTestSuite(t *testing.T) {
	suite.Run(t, new(UsersMapperTestSuite))
}
