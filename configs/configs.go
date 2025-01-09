package configs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DB MySQLConfig
}

type MySQLConfig struct {
	Username string
	Password string
}

func LoadConfig() (*Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = godotenv.Load(filepath.Join(pwd, "../../.env"))
	if err != nil {
		log.Fatalf("Error loading .env file: %s.", err)
	}
	cfg := &Config{
		DB: MySQLConfig{
			Username: os.Getenv("DBUSER"),
			Password: os.Getenv("DBPWD"),
		},
	}
	return cfg, nil
}
