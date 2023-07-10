package main

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/medmouine/device-mapper/cmd"
	"github.com/medmouine/device-mapper/internal/config"
	"github.com/medmouine/device-mapper/internal/mqtt"
	"github.com/medmouine/device-mapper/internal/router"
)

func main() {
	// Create router.
	r := chi.NewRouter()

	// Create config.
	c := config.NewConfig()

	// Set a logger middleware.
	r.Use(middleware.Logger)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(c.Server.ReadTimeout))

	// Get router with all routes.
	router.GetRoutes(r)

	mqttc, err := mqtt.NewClient(c.Mqtt.ToClientOptions())
	if err != nil {
		log.Fatal(err)
	}

	// Run server instance.
	if err := cmd.Run(c, r, mqttc); err != nil {
		log.Fatal(err)
		return
	}
}
