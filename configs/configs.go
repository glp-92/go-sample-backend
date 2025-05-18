package configs

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   MySQLConfig
	Auth AuthConfig
	API  APIConfig
}

type MySQLConfig struct {
	Username string
	Password string
	Address  string
	DBName   string
}

type AuthConfig struct {
	JWTSignKey                []byte
	JWTAccessTokenExpiration  int
	JWTRefreshTokenExpiration int
}

type APIConfig struct {
	APIPort string
}

func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid number on reading exp time: %s", s)
	}
	return i
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
			Address:  os.Getenv("DBADDR"),
			DBName:   os.Getenv("DBNAME"),
		},
		Auth: AuthConfig{
			JWTSignKey:                []byte(os.Getenv("JWTSIGNKEY")),
			JWTAccessTokenExpiration:  ToInt(os.Getenv("JWTACCESSTOKENEXPIRATION")),
			JWTRefreshTokenExpiration: ToInt(os.Getenv("JWTREFRESHTOKENEXPIRATION")),
		},
		API: APIConfig{
			APIPort: os.Getenv("API_PORT"),
		},
	}
	return cfg, nil
}
