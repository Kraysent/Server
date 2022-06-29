package actions

import (
	"server/pkg/core/entities"
	db "server/pkg/core/storage"
)

func (a *StorageAction) GetUser(login string) (*entities.User, error) {
	return a.Storage.GetUser(db.UsersFindParams{
		Login: &login,
	})
}

func (a *StorageAction) CreateUser(login string, password string, description string) (*entities.User, error) {
	salt := GenerateSalt(nil)
	passwordHash := HashPassword(password, salt)

	user, err := a.Storage.CreateUser(db.UserCreateParams{
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
