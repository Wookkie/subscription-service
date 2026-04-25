package main

import (
	"log"

	"github.com/Wookkie/subscription-service/internal/config"
	"github.com/Wookkie/subscription-service/internal/server"
)

func main() {
	cfg := config.ReadConfig()

	api := server.New(cfg)

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
