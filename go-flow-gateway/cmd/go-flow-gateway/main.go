package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bgg/go-flow-gateway/config"
	"github.com/bgg/go-flow-gateway/internal/app"
)

//	@tile			Go File Gate API
//	@version		1.0
//	@description	This is a Go File Gate API server.

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	// cfg, err := config.NewConfig("config/config.yml") for docker
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	cfg, err := config.NewConfig(fmt.Sprintf("../../config/config.%s.yml", env))

	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	app.Run(cfg)
}
