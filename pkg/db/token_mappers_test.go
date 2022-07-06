package db

import (
	"server/pkg/core/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TokenMapperTestSuite struct {
	BaseStorageSuite
}

func (s *TokenMapperTestSuite) SetupTest() {
	s.BaseStorageSuite.SetupTest()

	s.createTestUser("test1", 1, "hash1", "desc1", time.Unix(1656362017, 0).UTC())
	s.createTestUser("test2", 2, "hash2", "desc2", time.Unix(1656363017, 0).UTC())
	s.createTestUser("test3", 3, "hash3", "desc3", time.Unix(1656364017, 0).UTC())
}

func (s *TokenMapperTestSuite) TestCreate() {
	actual, err := s.storage.CreateToken(TokenCreateParams{
		UserID: 1, Value: "aaaaaaaaa", StartDate: time.Unix(1656363017, 0).UTC(),
		ExpirationDate: time.Unix(1656364017, 0).UTC(),
	})
	require.NoError(s.T(), err)

	expected := entities.Token{
		ID: 1, Value: "aaaaaaaaa", UserID: 1, StartDate: time.Unix(1656363017, 0).UTC(),
		ExpirationDate: time.Unix(1656364017, 0).UTC(),
	}
	assert.Equal(s.T(), &expected, actual)
}

func (s *TokenMapperTestSuite) TestFind() {
	s.createTestToken(1, "aaa", time.Unix(1656361017, 0).UTC(), time.Unix(1656365017, 0).UTC())
	s.createTestToken(2, "bbb", time.Unix(1656366017, 0).UTC(), time.Unix(1656368017, 0).UTC())
	s.createTestToken(1, "ccc", time.Unix(1656361017, 0).UTC(), time.Unix(1656369017, 0).UTC())

	userId := 1
	t := time.Unix(1656363017, 0).UTC()
	actual, err := s.storage.FindTokens(TokenFindParams{
		UserID: &userId, Time: &t,
	})
	s.Require().NoError(err)

	expected := []entities.Token{
		{
			ID: 1, UserID: 1, Value: "aaa", StartDate: time.Unix(1656361017, 0).UTC(),
			ExpirationDate: time.Unix(1656365017, 0).UTC(),
		},
		{
			ID: 3, UserID: 1, Value: "ccc", StartDate: time.Unix(1656361017, 0).UTC(),
			ExpirationDate: time.Unix(1656369017, 0).UTC(),
		},
	}
	assert.Equal(s.T(), expected, actual)
}

func TestTokenMapperTestSuite(t *testing.T) {
	suite.Run(t, new(TokenMapperTestSuite))
}
