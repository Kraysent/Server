package actions

import (
	"server/pkg/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TokenActionsTestSuite struct {
	BaseActionsTestSuite
}

type testTokenRow struct {
	userID    int
	value     string
	startDate time.Time
	endDate   time.Time
}

func (s *TokenActionsTestSuite) SetupTest() {
	s.BaseActionsTestSuite.SetupTest()

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
		_, err := s.action.Storage.CreateUser(db.UserCreateParams{
			Login: row.login, Salt: row.salt, PasswordHash: row.hash,
			Description: &row.description, RegistrationDate: &row.registrationDate,
		})
		require.NoError(s.T(), err)
	}
}

func (s *TokenActionsTestSuite) TestCheckTokenCaseExists() {
	testTokenRows := []testTokenRow{
		{userID: 1, value: "aaa", startDate: time.Unix(1656361017, 0).UTC(), endDate: time.Unix(1656365017, 0).UTC()},
		{userID: 2, value: "bbb", startDate: time.Unix(1656366017, 0).UTC(), endDate: time.Unix(1656368017, 0).UTC()},
		{userID: 1, value: "ccc", startDate: time.Unix(1656361017, 0).UTC(), endDate: time.Unix(1656369017, 0).UTC()},
	}

	for _, row := range testTokenRows {
		_, err := s.action.Storage.CreateToken(db.TokenCreateParams{
			UserID: row.userID, Value: row.value, StartDate: row.startDate, ExpirationDate: row.endDate,
		})
		require.NoError(s.T(), err)
	}

	actual, err := s.action.CheckUserToken("test1", "aaa", time.Unix(1656362017, 0).UTC())
	require.NoError(s.T(), err)

	assert.True(s.T(), actual)
}

func (s *TokenActionsTestSuite) TestCheckTokenCaseDoesNotExist() {
	testTokenRows := []testTokenRow{
		{userID: 1, value: "aaa", startDate: time.Unix(1656361017, 0).UTC(), endDate: time.Unix(1656365017, 0).UTC()},
		{userID: 2, value: "bbb", startDate: time.Unix(1656366017, 0).UTC(), endDate: time.Unix(1656368017, 0).UTC()},
		{userID: 1, value: "ccc", startDate: time.Unix(1656361017, 0).UTC(), endDate: time.Unix(1656369017, 0).UTC()},
	}

	for _, row := range testTokenRows {
		_, err := s.action.Storage.CreateToken(db.TokenCreateParams{
			UserID: row.userID, Value: row.value, StartDate: row.startDate, ExpirationDate: row.endDate,
		})
		require.NoError(s.T(), err)
	}

	actual, err := s.action.CheckUserToken("test1", "aba", time.Unix(1656362017, 0).UTC())
	require.NoError(s.T(), err)

	assert.False(s.T(), actual)
}

func TestTokenActionsTestSuite(t *testing.T) {
	suite.Run(t, new(TokenActionsTestSuite))
}
