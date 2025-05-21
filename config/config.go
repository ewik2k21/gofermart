package config

import "github.com/sirupsen/logrus"
import "github.com/joho/godotenv"

func LoadEnviroment() {
	if err := godotenv.Load("app.env"); err != nil {
		logrus.Fatalf("Error loading .env file")
	}
}

// set up db  ...
