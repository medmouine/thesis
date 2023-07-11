package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/medmouine/device-mapper/cmd"
	"github.com/medmouine/device-mapper/internal/client"
	"github.com/medmouine/device-mapper/internal/config"
	"github.com/medmouine/device-mapper/internal/router"
	temperaturesensor "github.com/medmouine/device-mapper/pkg/sensor"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfig()
	driver := temperaturesensor.NewTemperatureSimulator(cfg.Mqtt.ClientID, 0, 100)
	api := setupAPI(cfg.Server, driver)
	clt := client.NewClient(cfg.Mqtt.ToClientOptions(), driver)

	mapper := &cmd.Mapper{
		Config:       cfg,
		Client:       clt,
		DeviceDriver: driver,
		API:          api,
	}

	// Run instance.
	if err := mapper.Run(); err != nil {
		log.Fatal(err)
	}
}

func setupAPI(config *config.ServerConfig, driver *temperaturesensor.TemperatureSimulator) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(config.ReadTimeout))
	router.GetRoutes(r, driver)

	return r
}
