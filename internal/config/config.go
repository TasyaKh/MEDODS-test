package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Port                  string `env:"PORT" envDefault:"8080"`
	JWTSecret             string `env:"JWT_SECRET,required"`
	JWTAccessExpirySEC    int    `env:"JWT_ACCESS_EXPIRY_SEC" envDefault:"900"`
	RefreshTokenExpirySEC int    `env:"REFRESH_TOKEN_EXPIRY_SEC" envDefault:"2592000"`
	WebhookURL            string `env:"WEBHOOK_URL"`

	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDBName   string `env:"POSTGRES_DB_NAME,required"`

	AutoMigrate bool `env:"AUTO_MIGRATE" envDefault:"false"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
