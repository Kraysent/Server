package actions

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/exp/slices"
	"math/rand"
	db "server/pkg/core/storage"
)

func IssueToken(storage *db.Storage, login string) (string, error) {
	err := storage.Connect()
	if err != nil {
		return "", err
	}
	defer storage.Disconnect()

	rnd := rand.Int()
	tokenBytes := md5.Sum([]byte(fmt.Sprint(rnd)))
	token := fmt.Sprintf("%x", tokenBytes)

	user, err := storage.GetUser(login)
	if err != nil {
		return "", err
	}

	value, err := storage.CreateToken(user.ID, token, storage.Config.Token)
	if err != nil {
		return "", err
	}

	return value, nil
}

func CheckUserToken(storage *db.Storage, login string, token string) (bool, error) {
	err := storage.Connect()
	if err != nil {
		return false, err
	}
	defer storage.Disconnect()

	user, err := storage.GetUser(login)
	if err != nil {
		return false, err
	}

	tokens, err := storage.FindValidTokens(user.ID)
	if err != nil {
		return false, err
	}

	return slices.Contains(tokens, token), nil
}
