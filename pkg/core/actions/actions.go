package actions

import (
	"server/pkg/core/entities"
	db "server/pkg/core/storage"
)

func GetUser(storage *db.Storage, login string) (*entities.User, error) {
	err := storage.Connect()
	if err != nil {
		return nil, err
	}
	defer storage.Disconnect()

	return storage.GetUser(db.UsersFindParams{
		Login: &login,
	})
}

func CreateUser(storage *db.Storage, login string, password string, description string) (*entities.User, error) {
	err := storage.Connect()
	if err != nil {
		return nil, err
	}
	defer storage.Disconnect()

	salt := GenerateSalt(nil)
	passwordHash := HashPassword(password, salt)

	user, err := storage.CreateUser(db.UserCreateParams{
		Login:        login,
		Salt:         salt,
		PasswordHash: passwordHash,
		Description:  &description,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
