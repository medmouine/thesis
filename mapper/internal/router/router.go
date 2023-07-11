package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/medmouine/device-mapper/internal/router/device"
	"github.com/medmouine/device-mapper/internal/router/healthcheck"
	temperaturesensor "github.com/medmouine/device-mapper/pkg/sensor"
)

// GetRoutes function for getting routes.
func GetRoutes(m *chi.Mux, d *temperaturesensor.TemperatureSimulator) {
	healthcheck.Routes(m) // health check routes
	device.Routes(m, d)
	m.NotFound(http.NotFound) // not found routes
}
