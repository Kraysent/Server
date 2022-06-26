package storage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	zlog "github.com/rs/zerolog/log"
)

type Storage struct {
	Config     StorageConfig
	connection *pgx.Conn
}

type StorageConfig struct {
	User     string
	Password string
	Port     int
	DBName   string
}

func (s *Storage) Connect() error {
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf(
			"host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
			s.Config.Port, s.Config.User, s.Config.Password, s.Config.DBName,
		),
	)

	s.connection = conn
	return err
}

func (s *Storage) Disconnect() error {
	return s.connection.Close(context.Background())
}

func NewStorage(config StorageConfig) *Storage {
	return &Storage{
		Config: config,
	}
}

func (s *Storage) runQuery(ctx context.Context, query squirrel.Sqlizer) (pgx.Rows, error) {
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	zlog.Debug().Str("query", sqlQuery).Interface("args", args).Msg("Run query")

	rows, err := s.connection.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}
