package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/jackc/pgx/v4"
)

func AddUserRequestFunction(user string, password string, port int, database string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := pgx.Connect(
			context.Background(),
			fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", user, password, port, database),
		)
		if err != nil {
			log.Fatalf("Unable to connect to database. Error: %s", err)
		}
		defer conn.Close(context.Background())

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
		password_hash_hex := md5.Sum([]byte(fmt.Sprintf("%s%d", password, salt)))
		password_hash := fmt.Sprintf("%x", password_hash_hex)

		query := fmt.Sprintf(
			"INSERT INTO users (login, salt, password_hash, description) VALUES ('%s', '%d', '%s', '%s');",
			login, salt, password_hash, description,
		)
		fmt.Println(query)
		err = conn.QueryRow(context.Background(), query).Scan()
		if err != nil {
			log.Fatalf("Unable to insert to database. Error: %s", err)
		}

		fmt.Fprintf(w, "Wrote user '%s' data to DB.", login)
	}
}
