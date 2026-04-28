package main

import (
	"github.com/rs/zerolog/log"

	"github.com/Wookkie/subscription-service/internal/config"
	"github.com/Wookkie/subscription-service/internal/server"
)

func main() {
	cfg := config.ReadConfig()

	log.Info().Msg("starting server")

	api := server.New(cfg)

	if err := api.Run(); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}
