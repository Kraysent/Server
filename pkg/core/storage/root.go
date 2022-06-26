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
	DSN      string `config:"dsn"`
	Host     string `config:"host"`
	User     string `config:"user"`
	Password string `config:"password"`
	Port     int    `config:"port"`
	DBName   string `config:"dbname"`
}

func (s *Storage) Connect() error {
	var connectionString string

	if connectionString = s.Config.DSN; connectionString == "" {
		connectionString = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			s.Config.Host, s.Config.Port, s.Config.User, s.Config.Password, s.Config.DBName,
		)
	}

	conn, err := pgx.Connect(context.Background(), connectionString)
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
