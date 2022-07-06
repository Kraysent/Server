package db

import (
	"context"
	"database/sql"
	"server/pkg/core/entities"
	"server/pkg/core/storage"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"golang.org/x/xerrors"
)

type UsersFindParams struct {
	ID               *int
	Login            *string
	LoginLike        *string
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

func (s *Storage) executeUserQuery(query squirrel.Sqlizer) ([]entities.User, error) {
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

func scanUser(rows pgx.Rows) (entities.User, error) {
	var result entities.User
	var cityID sql.NullInt32

	err := rows.Scan(
		&result.ID, &result.Login, &result.Salt, &result.PasswordHash,
		&result.Description, &cityID, &result.RegistrationDate,
	)
	if err != nil {
		return entities.User{}, err
	}

	result.CityID = int(cityID.Int32)

	return result, nil
}

func (s *Storage) FindUsers(params UsersFindParams) ([]entities.User, error) {
	filters := storage.NewFilters()
	filters.AddEqual(entities.UserFieldID, params.ID)
	filters.AddEqual(entities.UserFieldLogin, params.Login)
	filters.AddLike(entities.UserFieldLogin, params.LoginLike)
	filters.AddEqual(entities.UserFieldDescription, params.Description)
	filters.AddEqual(entities.UserFieldCityID, params.CityID)
	filters.AddEqual(entities.UserFieldRegistrationDate, params.RegistrationDate)

	query := squirrel.Select("*").
		From(entities.TableUsers).
		Where(filters.Condition).
		PlaceholderFormat(squirrel.Dollar)

	return s.executeUserQuery(query)
}

func (s *Storage) GetUser(login string) (*entities.User, error) {
	users, err := s.FindUsers(UsersFindParams{
		Login: &login,
	})

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
	columns := []string{
		entities.UserFieldLogin, entities.UserFieldSalt, entities.UserFieldPasswordHash,
	}
	values := []any{params.Login, params.Salt, params.PasswordHash}

	if params.Description != nil {
		columns = append(columns, entities.UserFieldDescription)
		values = append(values, &params.Description)
	}
	if params.CityID != nil {
		columns = append(columns, entities.UserFieldCityID)
		values = append(values, &params.CityID)
	}
	if params.RegistrationDate != nil {
		columns = append(columns, entities.UserFieldRegistrationDate)
		values = append(values, &params.RegistrationDate)
	}

	query := squirrel.Insert(entities.TableUsers).
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar)

	result, err := s.executeUserQuery(query)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, xerrors.New("INSERT command returned no rows.")
	}

	return &result[0], nil
}
