package database

import (
	"database/sql"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/pkg/log"
)

func InitDB(source string) (*sql.DB, error) {
	logger := log.InitLogger()
	db, err := sql.Open("postgres", source)
	if err != nil {
		logger.WithError(err).Fatal("failed connect database")
		return nil, err
	}

	initStorage(db)
	return db, nil
}

func initStorage(db *sql.DB) error {
	logger := log.InitLogger()
	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("failed ping database")
		return err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	logger.Info("success connected database")
	return nil
}
