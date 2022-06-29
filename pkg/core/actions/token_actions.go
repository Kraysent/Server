package actions

import (
	"golang.org/x/exp/slices"
	db "server/pkg/core/storage"
)

func (a *StorageAction) IssueToken(login string) (string, error) {
	token := GenerateToken(nil)

	user, err := a.Storage.GetUser(db.UsersFindParams{
		Login: &login,
	})
	if err != nil {
		return "", err
	}

	value, err := a.Storage.CreateToken(user.ID, token, a.Storage.Config.Token)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (a *StorageAction) CheckUserToken(login string, token string) (bool, error) {
	user, err := a.Storage.GetUser(db.UsersFindParams{
		Login: &login,
	})
	if err != nil {
		return false, err
	}

	tokens, err := a.Storage.FindValidTokens(user.ID)
	if err != nil {
		return false, err
	}

	return slices.Contains(tokens, token), nil
}
