package db

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

func (s *UsersMapperTestSuite) TestCreateNew() {
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

func (s *UsersMapperTestSuite) TestCreateCollision() {
	s.generateTestUsers()

	description := "description"
	registration_time := time.Unix(1656363018, 0).UTC()
	_, err := s.storage.CreateUser(UserCreateParams{
		Login: "test1", Salt: 5, PasswordHash: "hash",
		Description: &description, RegistrationDate: &registration_time,
	})
	assert.Error(s.T(), err)
}

func (s *UsersMapperTestSuite) TestFindNoConditions() {
	s.generateTestUsers()

	actual, err := s.storage.FindUsers(UsersFindParams{})
	require.NoError(s.T(), err)

	expected := []entities.User{
		{
			ID: 1, Login: "test1", Salt: 1, PasswordHash: "hash1",
			Description: "desc1", RegistrationDate: time.Unix(1656362017, 0).UTC(),
		},
		{
			ID: 2, Login: "test2", Salt: 2, PasswordHash: "hash2",
			Description: "desc2", RegistrationDate: time.Unix(1656363017, 0).UTC(),
		},
		{
			ID: 3, Login: "test3", Salt: 3, PasswordHash: "hash3",
			Description: "desc3", RegistrationDate: time.Unix(1656364017, 0).UTC(),
		},
	}

	assert.Equal(s.T(), expected, actual)
}

func (s *UsersMapperTestSuite) TestFindWithFieldConditions() {
	s.generateTestUsers()

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

func (s *UsersMapperTestSuite) generateTestUsers() {
	s.createTestUser("test1", 1, "hash1", "desc1", time.Unix(1656362017, 0).UTC())
	s.createTestUser("test2", 2, "hash2", "desc2", time.Unix(1656363017, 0).UTC())
	s.createTestUser("test3", 3, "hash3", "desc3", time.Unix(1656364017, 0).UTC())
}

func TestUserMapperTestSuite(t *testing.T) {
	suite.Run(t, new(UsersMapperTestSuite))
}
