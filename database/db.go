package database

import (
	"database/sql"
	"fmt"
	"online-university/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}
