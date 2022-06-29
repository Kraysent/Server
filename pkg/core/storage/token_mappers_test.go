package storage

import (
	"server/pkg/core/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokenMapperTestSuite struct {
	BaseStorageSuite
}

func (s *TokenMapperTestSuite) SetupTest() {
	s.BaseStorageSuite.SetupTest()

	testUserRows := []struct {
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

	for _, row := range testUserRows {
		_, err := s.storage.CreateUser(UserCreateParams{
			Login: row.login, Salt: row.salt, PasswordHash: row.hash,
			Description: &row.description, RegistrationDate: &row.registrationDate,
		})
		s.Require().NoError(err)
	}
}

func (s *TokenMapperTestSuite) TestCreate() {
	actual, err := s.storage.CreateToken(TokenCreateParams{
		UserID:         1,
		Value:          "aaaaaaaaa",
		StartDate:      time.Unix(1656363017, 0).UTC(),
		ExpirationDate: time.Unix(1656364017, 0).UTC(),
	})
	s.Require().NoError(err)

	expected := entities.Token{
		ID: 1, Value: "aaaaaaaaa", UserID: 1,
		StartDate: time.Unix(1656363017, 0).UTC(), ExpirationDate: time.Unix(1656364017, 0).UTC(),
	}
	assert.Equal(s.T(), &expected, actual)
}

func (s *TokenMapperTestSuite) TestFind() {
	testTokenRows := []struct {
		userID    int
		value     string
		startDate time.Time
		endDate   time.Time
	}{
		{userID: 1, value: "aaa", startDate: time.Unix(1656361017, 0).UTC(), endDate: time.Unix(1656365017, 0).UTC()},
		{userID: 2, value: "bbb", startDate: time.Unix(1656366017, 0).UTC(), endDate: time.Unix(1656368017, 0).UTC()},
		{userID: 1, value: "ccc", startDate: time.Unix(1656361017, 0).UTC(), endDate: time.Unix(1656369017, 0).UTC()},
	}

	for _, row := range testTokenRows {
		_, err := s.storage.CreateToken(TokenCreateParams{
			UserID: row.userID, Value: row.value, StartDate: row.startDate, ExpirationDate: row.endDate,
		})
		s.Require().NoError(err)
	}

	userId := 1
	time := time.Unix(1656363017, 0).UTC()
	actual, err := s.storage.FindTokens(TokenFindParams{
		UserID: &userId, Time: &time,
	})
	s.Require().NoError(err)

	expected := []entities.Token{
		{ID: 1, UserID: 1, Value: "aaa", StartDate: testTokenRows[0].startDate, ExpirationDate: testTokenRows[0].endDate},
		{ID: 3, UserID: 1, Value: "ccc", StartDate: testTokenRows[2].startDate, ExpirationDate: testTokenRows[2].endDate},
	}
	assert.Equal(s.T(), expected, actual)
}

func TestTokenMapperTestSuite(t *testing.T) {
	suite.Run(t, new(TokenMapperTestSuite))
}
