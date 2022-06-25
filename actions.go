package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func CreateUser(login string, password string, description string) (*User, error) {
	salt := rand.Intn(1000000)
	passwordHashHex := md5.Sum([]byte(fmt.Sprintf("%s%d", password, salt)))
	passwordHash := fmt.Sprintf("%x", passwordHashHex)

	user, err := Create(User{
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