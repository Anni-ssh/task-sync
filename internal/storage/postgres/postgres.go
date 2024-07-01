package postgres

import (
	"database/sql"
	"fmt"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func (cfg Config) Prepare() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	const operation = "storage.NewPostgresDB"

	db, err := sql.Open("postgres", cfg.Prepare())
	if err != nil {
		return nil, fmt.Errorf("%s- failed to open database connection: %w", operation, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s - failed to ping database: %w", operation, err)
	}

	return db, nil

}
