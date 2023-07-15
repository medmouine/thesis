package main

import (
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/medmouine/mapper/cmd"
	"github.com/medmouine/mapper/internal/client"
	"github.com/medmouine/mapper/internal/config"
	"github.com/medmouine/mapper/internal/router"
	"github.com/medmouine/mapper/pkg/device"
	"github.com/medmouine/mapper/pkg/device/temperature"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfig()
	d := temperature.NewTemperatureSimulator(cfg.Mqtt.ClientID, cfg.Mqtt.PublishInterval, cfg.Mqtt.MinTemp, cfg.Mqtt.MaxTemp)
	api := setupAPI[temperature.TemperatureData](cfg.Server, d)
	clt := client.NewClient[temperature.TemperatureData](d, cfg.Mqtt.ToClientOptions())

	mapper := &cmd.Mapper{
		Config: cfg,
		Client: clt,
		Device: d,
		API:    api,
	}

	// Run instance.
	if err := mapper.Run(); err != nil {
		log.Fatal(err)
	}
}

func setupAPI[T interface{}](config *config.ServerConfig, d device.Device[T]) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(config.ReadTimeout))
	router.GetRoutes(r, d)

	return r
}
