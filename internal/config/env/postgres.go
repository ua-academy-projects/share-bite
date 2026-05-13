package env

import (
	"fmt"
)

const (
	postgresSSLAllow   = "allow"
	postgresSSLDisable = "disable"
)

type postgresConfig struct {
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     string `env:"POSTGRES_PORT,required"`
	PostgresSSL      string `env:"POSTGRES_SSL,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`

	PostgresMigrationsDir string `env:"POSTGRES_MIGRATIONS_DIR,required"`
}

func NewPostgresConfig(opts ...Options) (*postgresConfig, error) {
	config := new(postgresConfig)
	if err := Parse(config, opts...); err != nil {
		return nil, err
	}

	if config.PostgresSSL != postgresSSLAllow && config.PostgresSSL != postgresSSLDisable {
		return nil, fmt.Errorf(`unknown ssl option: %s (only %s or %s)`, config.PostgresSSL, postgresSSLAllow, postgresSSLDisable)
	}

	return config, nil
}

func (c *postgresConfig) Dsn() string {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=public",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresSSL,
	)

	return connString
}

func (c *postgresConfig) MigrationsDir() string {
	return c.PostgresMigrationsDir
}
