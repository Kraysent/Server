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
			fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", user, password, port, database),
		)
		if err != nil {
			log.Fatalf("Unable to connect to database. Error: %s", err)
		}
		defer conn.Close(context.Background())

		query, _, err := squirrel.Select("login").
			From("users").
			ToSql()
		if err != nil {
			log.Fatalf("unable to create request. Error: %s", err)
		}

		log.Printf("Running SQL query '%s'", query)

		rows, err := conn.Query(context.Background(), query)
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
