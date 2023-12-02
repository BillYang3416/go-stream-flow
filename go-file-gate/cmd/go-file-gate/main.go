package main

import (
	"log"

	"github.com/bgg/go-file-gate/config"
	"github.com/bgg/go-file-gate/internal/app"
)

//	@tile			Go File Gate API
//	@version		1.0
//	@description	This is a Go File Gate API server.

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	// cfg, err := config.NewConfig("config/config.yml") for docker
	cfg, err := config.NewConfig("../../config/config.yml")

	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	app.Run(cfg)
}
