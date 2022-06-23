package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type User struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	Salt         int    `json:"salt"`
	PasswordHash string `json:"password_hash"`
	Description  string `json:"description"`
}

func RunQuery(connection *pgx.Conn, ctx context.Context, query squirrel.Sqlizer) (pgx.Rows, error) {
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	log.Printf("Running SQL query '%s' %% %s", sqlQuery, args)

	rows, err := connection.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}

func AddUserRequestFunction(user string, password string, port int, database string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := pgx.Connect(
			context.Background(),
			fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable", port, user, password, database),
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

		query := squirrel.Insert("users").
			Columns("login", "salt", "password_hash", "description").
			Values(login, salt, password_hash, description).
			PlaceholderFormat(squirrel.Dollar).
			Suffix("RETURNING id,login")

		rows, err := RunQuery(conn, context.Background(), query)
		if err != nil {
			log.Fatalf("unable to insert to database. Error: %s", err)
		}

		result := User{}
		for rows.Next() {
			err := rows.Scan(&result.Id, &result.Login)
			if err != nil {
				log.Fatalf("unable to insert to database. Error: %s", err)
			}

			fmt.Fprintf(w, "Wrote user '%s' data to DB (id: %d).", result.Login, result.Id)
		}
	}
}
