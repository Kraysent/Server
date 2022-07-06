package core

import (
	"io/ioutil"
	"server/pkg/db"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port  int    `yaml:"port"`
	Level string `yaml:"level"`
}

type Config struct {
	Storage db.StorageConfig `yaml:"storage"`
	Server  ServerConfig     `yaml:"server"`
}

func NewConfig(configFilename string) (Config, error) {
	config := Config{
		Server: ServerConfig{
			Level: "info",
		},
	}

	configFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
