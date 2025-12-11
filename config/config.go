package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbAddr        string
	DbName        string
	DbUser        string
	DbPassword    string
	Port          string
	JWTExpiry     int
	RefreshExpiry int
}

func GetConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}
	c := Config{
		DbAddr:        os.Getenv("DB_ADDR"),
		DbName:        os.Getenv("DB_NAME"),
		DbUser:        os.Getenv("DB_USER"),
		DbPassword:    os.Getenv("DB_PASSWORD"),
		Port:          os.Getenv("PORT"),
		JWTExpiry:     3600,
		RefreshExpiry: 604800,
	}

	if c.Port == "" {
		c.Port = "8080"
	}

	if c.DbAddr == "" {
		c.DbAddr = "localhost:3306"
	}

	if c.DbName == "" || c.DbUser == "" || c.DbPassword == "" {
		return c, fmt.Errorf("missing dbName, dbUser and dbPassword environment variables")
	}

	return c, nil
}
