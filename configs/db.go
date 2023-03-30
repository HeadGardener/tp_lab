package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type DBConfig struct {
	DBName  string
	Host    string
	SSLMode string
}

func NewDBConfig(path string) (*DBConfig, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	dbname := os.Getenv("dbname")
	if dbname == "" {
		return nil, errors.New("db name is empty")
	}

	host := os.Getenv("dbhost")
	if host == "" {
		return nil, errors.New("db host is empty")
	}

	sslmode := os.Getenv("sslmode")
	if sslmode == "" {
		return nil, errors.New("sslmode is empty")
	}

	return &DBConfig{
		DBName:  dbname,
		Host:    host,
		SSLMode: sslmode,
	}, nil
}
