package database

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// OpenSQLX function returns a sqlx.DB connection.
func OpenSQLX(connectionString string) (*sqlx.DB, error) {
	maxConn, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	if err != nil {
		return nil, errors.New("missing env variable: DB_MAX_CONNECTIONS")
	}
	maxIdleConn, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	if err != nil {
		return nil, errors.New("missing env variable: DB_MAX_IDLE_CONNECTIONS")
	}
	maxLifetimeConn, err := time.
		ParseDuration(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))
	if err != nil {
		return nil, errors.
			New("missing env variable: DB_MAX_LIFETIME_CONNECTIONS")
	}

	db, err := sqlx.ConnectContext(
		context.Background(),
		"pgx",
		connectionString,
	)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(maxLifetimeConn)

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, err
	}

	return db, nil
}
