package storage

import (
	"context"
	"server/pkg/core/entities"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

const (
	usersTableName = "users"
)

func scanUser(rows pgx.Rows) (entities.User, error) {
	var result entities.User
	err := rows.Scan(
		&result.ID,
		&result.Login,
		&result.Salt,
		&result.PasswordHash,
		&result.Description,
		&result.CityID,
		&result.RegistrationDate,
	)

	if err != nil {
		return entities.User{}, err
	}

	return result, err
}

func (s *Storage) Find(login string) ([]entities.User, error) {
	query := squirrel.Select("*").
		From(usersTableName).
		Where(squirrel.Eq{"login": login}).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := s.runQuery(context.Background(), query)
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

func (s *Storage) Get(login string) (*entities.User, error) {
	users, err := s.Find(login)

	if len(users) >= 1 {
		return &users[0], err
	} else {
		return nil, err
	}
}

func (s *Storage) Create(user entities.User) (*entities.User, error) {
	query := squirrel.Insert(usersTableName).
		Columns("login", "salt", "password_hash", "description").
		Values(user.Login, user.Salt, user.PasswordHash, user.Description).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id,login,salt,password_hash,description,city,registration_date")

	rows, err := s.runQuery(context.Background(), query)
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
