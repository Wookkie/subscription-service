package config

import "flag"

type Config struct {
	Host string
	Port int
	DBConn string
}

func ReadConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.Host, "host", "0.0.0.0", "server host address")
	flag.IntVar(&cfg.Port, "port", 8080, "server port")
	flag.StringVar(&cfg.DBConn, "db", "postgres://app_user:password@localhost:5433/subscriptions?sslmode=disable", "database connection")

	flag.Parse()

	return &cfg
}
