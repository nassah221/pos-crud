package config

import (
	"log"
	"os"
	"sales-api/constants"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver    string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	BindAddr    string
	JWTSecret   string
	JWTDuration time.Duration
}

func LoadConfig() (Config, error) {
	config := Config{}

	if err := godotenv.Load(".env"); err != nil {
		return config, err
	}

	config.DBDriver = "mysql"
	config.BindAddr = ":3030"

	config.DBHost = os.Getenv(constants.DBHost)
	config.DBPort = os.Getenv(constants.DBPort)
	config.DBUser = os.Getenv(constants.DBUser)
	config.DBPassword = os.Getenv(constants.DBPassword)
	config.DBName = os.Getenv(constants.DBName)
	if os.Getenv(constants.JWTSecret) == "" {
		config.JWTSecret = "01234567890123456789012345678901"
	} else {
		config.JWTSecret = os.Getenv(constants.JWTSecret)
	}

	if os.Getenv(constants.JWTDuration) == "" {
		config.JWTDuration = time.Hour * 15
	} else {
		t, err := time.ParseDuration(os.Getenv(constants.JWTDuration))
		if err != nil {
			log.Panic(err)
		}
		config.JWTDuration = t
	}

	return config, nil
}
