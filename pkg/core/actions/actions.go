package actions

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"server/pkg/core/entities"
	db "server/pkg/core/storage"
)

func HashPassword(password string, salt int) string {
	passwordHashHex := md5.Sum([]byte(fmt.Sprintf("%s%d", password, salt)))
	return fmt.Sprintf("%x", passwordHashHex)
}

func GetUser(storage *db.Storage, login string) (*entities.User, error) {
	err := storage.Connect()
	if err != nil {
		return nil, err
	}
	defer storage.Disconnect()

	return storage.GetUser(login)
}

func CreateUser(storage *db.Storage, login string, password string, description string) (*entities.User, error) {
	err := storage.Connect()
	if err != nil {
		return nil, err
	}
	defer storage.Disconnect()

	salt := rand.Intn(1000000)
	passwordHashHex := md5.Sum([]byte(fmt.Sprintf("%s%d", password, salt)))
	passwordHash := fmt.Sprintf("%x", passwordHashHex)

	user, err := storage.CreateUser(entities.User{
		Login:        login,
		Salt:         salt,
		PasswordHash: passwordHash,
		Description:  description,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
