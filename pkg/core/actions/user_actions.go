package actions

import (
	"server/pkg/core/entities"
	"server/pkg/db"
)

func (a *StorageAction) GetUser(login string) (*entities.User, error) {
	return a.Storage.GetUser(login)
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

func (a *StorageAction) FindUserByPattern(loginLike string) ([]entities.User, error) {
	return a.Storage.FindUsers(db.UsersFindParams{
		LoginLike: &loginLike,
	})
}
