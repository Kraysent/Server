package actions

import (
	"server/pkg/core/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserActionsTestSuite struct {
	BaseActionsTestSuite
}

func (s *UserActionsTestSuite) TestCreate() {
	actual, err := s.action.CreateUser("test", "password", "description")
	s.Require().NoError(err)

	expected := entities.User{
		ID:          1,
		Login:       "test",
		Description: "description",
		CityID:      0,
	}
	assert.Equal(s.T(), expected, entities.User{
		ID: actual.ID, Login: actual.Login, Description: actual.Description, CityID: actual.CityID,
	})
}

func (s *UserActionsTestSuite) TestGet() {
	testUsersRows := []struct {
		login       string
		password    string
		description string
	}{
		{login: "test1", password: "password1", description: "desc1"},
		{login: "test2", password: "password2", description: "desc2"},
		{login: "test3", password: "password3", description: "desc3"},
	}

	for _, row := range testUsersRows {
		_, err := s.action.CreateUser(row.login, row.password, row.description)
		s.Require().NoError(err)
	}

	actual, err := s.action.GetUser("test3")
	s.Require().NoError(err)

	expected := entities.User{
		ID: 3, Login: "test3", Description: "desc3", CityID: 0,
	}

	assert.Equal(s.T(), expected, entities.User{
		ID: actual.ID, Login: actual.Login, Description: actual.Description, CityID: actual.CityID,
	})
}

func TestUserActionsTestSuite(t *testing.T) {
	suite.Run(t, new(UserActionsTestSuite))
}
