package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	logger "github.com/zsandibe/medods-service/pkg"
)

type Config struct {
	Server   serverConfig
	Postgres postgresConfig
	Token    tokenConfig
}

type postgresConfig struct {
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     string `envconfig:"DB_PORT" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
}

type serverConfig struct {
	Port string `envconfig:"SERVER_PORT" required:"true"`
}

type tokenConfig struct {
	AccessKey       string        `envconfig:"ACCESS_KEY" required:"true"`
	AccessTokenAge  time.Duration `envconfig:"ACCESS_TOKEN_AGE" required:"true"`
	RefreshTokenAge time.Duration `envconfig:"REFRESH_TOKEN_AGE" required:"true"`
}

func NewConfig(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		logger.Errorf("godotenv.Load(): %v", err)
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	var config Config

	if err := envconfig.Process("", &config); err != nil {
		logger.Errorf("envconfig.Process(): %v", err)
		return nil, fmt.Errorf("error processing .env file: %v", err)
	}
	return &config, nil
}
