package config

import (
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"os"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error loading .env file")
	}
}

func SetUpDatabase(ctx context.Context) (*sql.DB, error) {
	dsn := os.Getenv(DataBaseUrl)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		logrus.Fatalf("Failed to ping db: %v", err)
	}

	//migrations up for db
	if err := goose.Up(db, "migrations"); err != nil {
		logrus.Fatalf("Failed migrations for init db : %v", err)
	}

	return db, nil
}
