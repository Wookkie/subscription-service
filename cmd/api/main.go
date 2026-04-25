package main

import (
	"fmt"

	"github.com/Wookkie/subscription-service/internal/config"
)

func main() {
	cfg := config.ReadConfig()

	fmt.Printf("Host: %s\nPort: %d\n", cfg.Host, cfg.Port)
}
