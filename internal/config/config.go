package config

import "flag"

type Config struct {
	Host string
	Port int
}

func ReadConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.Host, "host", "0.0.0.0", "flag for configure host")
	flag.IntVar(&cfg.Port, "port", 8080, "flag for configure host")

	flag.Parse()

	return &cfg
}
