package conn

import (
	"database/sql"
	"fmt"
	"telebot-invent/config"
)

func DbConn() (*sql.DB, error) {
	pg_host := config.Config("PG_HOST")
	pg_user := config.Config("PG_USER")
	pg_password := config.Config("PG_PASSWORD")
	pg_dbname := config.Config("PG_DBNAME")
	pg_port := config.Config("PG_PORT")
	pg_sslmode := config.Config("PG_SSL_MODE")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", pg_host, pg_user, pg_password, pg_dbname, pg_port, pg_sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %s", err)
	}
	return db, nil
}
