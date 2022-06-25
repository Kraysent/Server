package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func AddUserRequest(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	if login == "" {
		login = "user"
	}

	password := r.URL.Query().Get("password")
	if password == "" {
		password = "p@ssw0rd"
	}

	description := r.URL.Query().Get("description")

	salt := rand.Intn(1000000)
	passwordHashHex := md5.Sum([]byte(fmt.Sprintf("%s%d", password, salt)))
	passwordHash := fmt.Sprintf("%x", passwordHashHex)

	_, err := Create(User{
		Login:        login,
		Salt:         salt,
		PasswordHash: passwordHash,
		Description:  description,
	})
	if err != nil {
		log.Fatal(err)
	}
}
