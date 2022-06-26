package storage

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	tokensTableName = "tokens"
)

func (s *Storage) CreateToken(user_id int, value string, ttl time.Duration) (string, error) {
	query := squirrel.Insert(tokensTableName).
		Columns("login", "value", "start_date", "expiration_date").
		Values(user_id, value, time.Now(), time.Now().Add(ttl)).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING value")

	rows, err := s.runQuery(context.Background(), query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func (s *Storage) FindValidTokens(user_id int) ([]string, error) {
	query := squirrel.Select("value").
		From(tokensTableName).
		Where(squirrel.GtOrEq{"expiration_date": time.Now()}).
		Where(squirrel.LtOrEq{"start_date": time.Now()}).
		Where(squirrel.Eq{"login": user_id}).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := s.runQuery(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		var curr string
		err = rows.Scan(&curr)
		if err != nil {
			return nil, err
		}

		result = append(result, curr)
	}

	return result, nil
}
