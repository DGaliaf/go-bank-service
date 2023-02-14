package main

import (
	"avito-tech/app/internal/app"
	"avito-tech/app/internal/config"
	"context"
	"log"
)

// @title           Swagger Example API
// @version         1.1
// @description     This is a sample banking service.

// @host      localhost:30001
// @BasePath  /
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("config initializing")
	cfg := config.GetConfig()

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, err)
	}

	log.Println("Running Application")
	a.Run(ctx)
}
