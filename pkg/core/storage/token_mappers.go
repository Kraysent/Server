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
	columns := []string{entities.TokenFieldUserID, entities.TokenFieldValue, entities.TokenFieldStartDate, entities.TokenFieldExpirationDate}
	values := []any{params.UserID, params.Value, params.StartDate, params.ExpirationDate}

	query := squirrel.Insert(entities.TableTokens).
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
	filters.AddEqual(entities.TokenFieldUserID, params.UserID)
	filters.AddEqual(entities.TokenFieldValue, params.Value)
	filters.AddInRange(entities.TokenFieldStartDate, entities.TokenFieldExpirationDate, params.Time)

	query := squirrel.Select("*").
		From(entities.TableTokens).
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
