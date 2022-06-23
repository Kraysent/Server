package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func GetUsersRequestFunction(user string, password string, port int, database string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := pgx.Connect(
			context.Background(),
			fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable", port, user, password, database),
		)
		if err != nil {
			log.Fatalf("Unable to connect to database. Error: %s", err)
		}
		defer conn.Close(context.Background())

		query := squirrel.Select("login").From("users")

		login := r.URL.Query().Get("login")
		if login != "" {
			query = query.Where(squirrel.Eq{
				"login": []string{login},
			})
		}

		query = query.PlaceholderFormat(squirrel.Dollar)

		rows, err := RunQuery(conn, context.Background(), query)
		if err != nil {
			log.Fatalf("unable to select from database. Error: %s", err)
		}

		defer rows.Close()

		var curr string
		result := make([]string, 0)
		for rows.Next() {
			err := rows.Scan(&curr)
			if err != nil {
				log.Fatalf("unable to read from database. Error: %s", err)
			}
			result = append(result, curr)
		}

		fmt.Fprintf(w, "%s", strings.Join(result, "\n"))
	}
}
