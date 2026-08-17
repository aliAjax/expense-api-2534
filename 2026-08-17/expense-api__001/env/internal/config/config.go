package config

import (
	"os"
	"strings"
)

type Config struct {
	Port   string
	DBPath string
}

func Load() Config {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	dbPath := strings.TrimSpace(os.Getenv("DB_PATH"))
	if dbPath == "" {
		dbPath = "./data/expense.db"
	}

	return Config{
		Port:   port,
		DBPath: dbPath,
	}
}
