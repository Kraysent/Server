package storage

import (
	"context"
	"database/sql"
	"server/pkg/core/entities"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type UsersFindParams struct {
	ID               *int
	Login            *string
	Description      *string
	CityID           *int
	RegistrationDate *time.Time
}

type UserCreateParams struct {
	Login            string
	Salt             int
	PasswordHash     string
	Description      *string
	CityID           *int
	RegistrationDate *time.Time
}

func scanUser(rows pgx.Rows) (entities.User, error) {
	var result entities.User
	var cityID sql.NullInt32

	err := rows.Scan(
		&result.ID,
		&result.Login,
		&result.Salt,
		&result.PasswordHash,
		&result.Description,
		&cityID,
		&result.RegistrationDate,
	)
	if err != nil {
		return entities.User{}, err
	}

	result.CityID = int(cityID.Int32)

	return result, err
}

func (s *Storage) FindUsers(params UsersFindParams) ([]entities.User, error) {
	filters := NewFilters()
	filters.AddEqual("id", params.ID)
	filters.AddEqual("login", params.Login)
	filters.AddEqual("description", params.Description)
	filters.AddEqual("city_id", params.CityID)
	filters.AddEqual("registration_date", params.RegistrationDate)

	query := squirrel.Select("*").
		From(usersTableName).
		Where(filters.Condition).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := s.RunQuery(context.Background(), query)
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

func (s *Storage) GetUser(params UsersFindParams) (*entities.User, error) {
	users, err := s.FindUsers(params)

	if err != nil {
		return nil, err
	}

	if len(users) >= 1 {
		return &users[0], nil
	} else {
		return nil, nil
	}
}

func (s *Storage) CreateUser(params UserCreateParams) (*entities.User, error) {
	columns := []string{"login", "salt", "password_hash"}
	values := []any{params.Login, params.Salt, params.PasswordHash}

	if params.Description != nil {
		columns = append(columns, "description")
		values = append(values, &params.Description)
	}
	if params.CityID != nil {
		columns = append(columns, "city")
		values = append(values, &params.CityID)
	}
	if params.RegistrationDate != nil {
		columns = append(columns, "registration_date")
		values = append(values, &params.RegistrationDate)
	}

	query := squirrel.Insert(usersTableName).
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar)

	rows, err := s.RunQuery(context.Background(), query)
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
