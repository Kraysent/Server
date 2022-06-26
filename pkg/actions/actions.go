package actions

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"server/pkg/entities"
	"server/pkg/storage"
)

func HashPassword(password string, salt int) string {
	passwordHashHex := md5.Sum([]byte(fmt.Sprintf("%s%d", password, salt)))
	return fmt.Sprintf("%x", passwordHashHex)
}

func CreateUser(login string, password string, description string) (*entities.User, error) {
	salt := rand.Intn(1000000)
	passwordHashHex := md5.Sum([]byte(fmt.Sprintf("%s%d", password, salt)))
	passwordHash := fmt.Sprintf("%x", passwordHashHex)

	user, err := storage.Create(entities.User{
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
