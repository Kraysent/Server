package db

import (
	"context"
	"server/pkg/core/entities"
	"server/pkg/core/storage"
	"time"

	"github.com/Masterminds/squirrel"
	"golang.org/x/xerrors"
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

func (s *Storage) executeTokenQuery(query squirrel.Sqlizer) ([]entities.Token, error) {
	rows, err := s.RunQuery(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.Token, 0)
	for rows.Next() {
		var curr entities.Token
		if err := rows.Scan(
			&curr.ID, &curr.UserID, &curr.Value, &curr.StartDate, &curr.ExpirationDate,
		); err != nil {
			return nil, err
		}

		result = append(result, curr)
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

	result, err := s.executeTokenQuery(query)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, xerrors.New("INSERT command returned no rows.")
	}

	return &result[0], nil
}

func (s *Storage) FindTokens(params TokenFindParams) ([]entities.Token, error) {
	filters := storage.NewFilters()
	filters.AddEqual(entities.TokenFieldUserID, params.UserID)
	filters.AddEqual(entities.TokenFieldValue, params.Value)
	filters.AddInRange(entities.TokenFieldStartDate, entities.TokenFieldExpirationDate, params.Time)

	query := squirrel.Select("*").
		From(entities.TableTokens).
		Where(filters.Condition).
		PlaceholderFormat(squirrel.Dollar)

	return s.executeTokenQuery(query)
}
