package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error loading .env file")
	}
}

func SetUpDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv(DataBaseUrl)
	dbPool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	err = dbPool.Ping(ctx)
	if err != nil {
		logrus.Fatalf("Failed to ping db: %v", err)
	}

	return dbPool, nil
}
