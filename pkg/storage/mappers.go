package storage

import (
	"context"
	"fmt"
	"log"
	"server/pkg/entities"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

const (
	postgresUsername = "testserver"
	postgresPassword = "passw0rd"
	postgresPort     = 5432
	postgresDBName   = "serverdb"
	usersTableName   = "users"
)

func runQuery(connection *pgx.Conn, ctx context.Context, query squirrel.Sqlizer) (pgx.Rows, error) {
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

func scanUser(rows pgx.Rows) (entities.User, error) {
	var result entities.User
	err := rows.Scan(&result.ID, &result.Login, &result.Salt, &result.PasswordHash, &result.Description)

	if err != nil {
		return entities.User{}, err
	}

	return result, err
}

func Find(login string) ([]entities.User, error) {
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf(
			"host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
			postgresPort, postgresUsername, postgresPassword, postgresDBName,
		),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := squirrel.Select("*").
		From(usersTableName).
		Where(squirrel.Eq{"login": login}).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := runQuery(conn, context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entities.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func Get(login string) (*entities.User, error) {
	users, err := Find(login)

	if len(users) >= 1 {
		return &users[0], err
	} else {
		return nil, err
	}
}

func Create(user entities.User) (*entities.User, error) {
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf(
			"host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
			postgresPort, postgresUsername, postgresPassword, postgresDBName,
		),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := squirrel.Insert(usersTableName).
		Columns("login", "salt", "password_hash", "description").
		Values(user.Login, user.Salt, user.PasswordHash, user.Description).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id,login,salt,password_hash,description")

	rows, err := runQuery(conn, context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result entities.User
	for rows.Next() {
		result, err = scanUser(rows)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}
