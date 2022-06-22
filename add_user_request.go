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

		query := "INSERT INTO users (login, salt, password_hash, description) VALUES ($1, $2, $3, $4);"
		fmt.Println(query)
		commandTag, err := conn.Exec(context.Background(), query, login, salt, password_hash, description)
		if err != nil {
			log.Fatalf("unable to insert to database. Error: %s", err)
		}
		if commandTag.RowsAffected() != 1 {
			log.Fatal("no rows were inserted")
		}

		fmt.Fprintf(w, "Wrote user '%s' data to DB.", login)
	}
}
