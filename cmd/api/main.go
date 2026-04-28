package main

import (
	"github.com/rs/zerolog/log"

	"github.com/Wookkie/subscription-service/internal/config"
	"github.com/Wookkie/subscription-service/internal/database"
	"github.com/Wookkie/subscription-service/internal/server"
)

// @title Subscription Service API
// @version 1.0
// @description Service for managing user subscriptions
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.ReadConfig()

	if err := database.ApplyMigrations(cfg.DBConn); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}

	log.Info().Msg("starting server")

	api := server.New(cfg)

	if err := api.Run(); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}
