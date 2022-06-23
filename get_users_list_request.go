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

		query := squirrel.Select("login").
			From("users")

		login := r.URL.Query().Get("login")
		if login != "" {
			query = query.Where(squirrel.Eq{
				"login": []string{ login },
			})
		}

		sqlQuery, args, err := query.PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			log.Fatalf("unable to construct SQL query. Error: %s", err)
		}

		log.Printf("Running SQL query '%s' %% %s", sqlQuery, args)

		rows, err := conn.Query(context.Background(), sqlQuery, args...)
		if err != nil {
			log.Fatalf("unable to select from database. Error: %s", err)
		}
		if err := rows.Err(); err != nil {
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
