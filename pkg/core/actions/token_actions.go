package actions

import (
	db "server/pkg/core/storage"
	"time"
)

func (a *StorageAction) IssueToken(login string) (string, error) {
	token := GenerateToken(nil)

	user, err := a.Storage.GetUser(db.UsersFindParams{
		Login: &login,
	})
	if err != nil {
		return "", err
	}

	value, err := a.Storage.CreateToken(db.TokenCreateParams{
		UserID: user.ID, Value: token, StartDate: time.Now(), ExpirationDate: time.Now().Add(a.Storage.Config.Token),
	})
	if err != nil {
		return "", err
	}

	return value.Value, nil
}

func (a *StorageAction) CheckUserToken(login string, token string, currTime time.Time) (bool, error) {
	user, err := a.Storage.GetUser(db.UsersFindParams{
		Login: &login,
	})
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}

	tokens, err := a.Storage.FindTokens(db.TokenFindParams{
		UserID: &user.ID, Value: &token, Time: &currTime,
	})
	if err != nil {
		return false, err
	}

	return len(tokens) != 0, nil
}
