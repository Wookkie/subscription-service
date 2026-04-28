package config

import (
	"os"
	"strconv"
)

type Config struct {
	Host   string
	Port   int
	DBConn string
}

func ReadConfig() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))

	return &Config{
		Host:   getEnv("HOST", "0.0.0.0"),
		Port:   port,
		DBConn: buildDBConn(),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func buildDBConn() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "app_user")
	pass := getEnv("DB_PASSWORD", "password")
	db := getEnv("DB_NAME", "subscriptions")
	ssl := getEnv("DB_SSLMODE", "disable")

	return "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + db + "?sslmode=" + ssl
}
