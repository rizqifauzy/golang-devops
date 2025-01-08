package main

import (
	"database/sql"
	"fmt"
)

func dbConn() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config("PG_HOST"), config("PG_USER"), config("PG_PASSWORD"), config("PG_DB_NAME"), config("PG_PORT"), config("PG_SSL_MODE"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %s", err)
	}
	return db, nil
}
