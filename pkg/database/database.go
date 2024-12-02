package database

import (
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

func InitDB(source string) (*sql.DB, error) {
	db, err := sql.Open("postgres", source)
	if err != nil {
		logrus.WithField("connect database", err.Error()).Fatal(err.Error())
		return nil, err
	}

	initStorage(db)
	return db, nil
}

func initStorage(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		logrus.WithField("connect database", err.Error()).Fatal(err.Error())
		return err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	logrus.WithField("connect database", "success connected").Info("success connected database")
	return nil
}
