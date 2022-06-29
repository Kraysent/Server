package storage

import (
	"context"
	"server/pkg/core/entities"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type TokenCreateParams struct {
	UserID         int
	Value          string
	StartDate      time.Time
	ExpirationDate time.Time
}

type TokenFindParams struct {
	UserID *int
	Value  *string
	//Will select all tokens for which Time is in range (start_date, expiration_date)
	Time *time.Time
}

func scanToken(rows pgx.Rows) (entities.Token, error) {
	var result entities.Token

	err := rows.Scan(
		&result.ID,
		&result.UserID,
		&result.Value,
		&result.StartDate,
		&result.ExpirationDate,
	)
	if err != nil {
		return entities.Token{}, err
	}

	return result, nil
}

func (s *Storage) CreateToken(params TokenCreateParams) (*entities.Token, error) {
	columns := []string{"user_id", "value", "start_date", "expiration_date"}
	values := []any{params.UserID, params.Value, params.StartDate, params.ExpirationDate}

	query := squirrel.Insert(tokensTableName).
		Columns(columns...).
		Values(values...).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING *")

	rows, err := s.RunQuery(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result entities.Token
	for rows.Next() {
		result, err = scanToken(rows)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *Storage) FindTokens(params TokenFindParams) ([]entities.Token, error) {
	filters := NewFilters()
	filters.AddEqual("user_id", params.UserID)
	filters.AddEqual("value", params.Value)
	filters.AddInRange("start_date", "expiration_date", params.Time)

	query := squirrel.Select("*").
		From(tokensTableName).
		Where(filters.Condition).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := s.RunQuery(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.Token, 0)
	for rows.Next() {
		curr, err := scanToken(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, curr)
	}

	return result, nil
}
