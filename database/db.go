package database

import (
	"database/sql"
	"fmt"
	"online-university/config"

	_ "github.com/lib/pq" // драйвер PostgreSQL (неиспользуемый импорт, но необходим)
)

type DB struct {
	Conn *sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{Conn: conn}, nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
