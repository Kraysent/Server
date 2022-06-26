package core

import (
	"context"
	db "server/pkg/core/storage"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type ServerConfig struct {
	Port  int    `config:"port,required"`
	Level string `config:"level"`
}

type Config struct {
	Storage db.StorageConfig `config:"storage"`
	Server  ServerConfig     `config:"server"`
}

func NewConfig(configFilename string) (Config, error) {
	config := Config{
		Server: ServerConfig{
			Level: "info",
		},
	}

	loader := confita.NewLoader(
		file.NewBackend("configs/dev.yaml"),
		env.NewBackend(),
	)
	err := loader.Load(context.Background(), &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
