package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/medmouine/thesis/mapper/internal/router/healthcheck"
	"github.com/medmouine/thesis/mapper/internal/router/state"
	"github.com/medmouine/thesis/mapper/pkg/device"
)

// GetRoutes function for getting routes.
func GetRoutes[T interface{}](m *chi.Mux, d device.Device[T]) {
	healthcheck.Routes(m) // health check routes
	state.Routes(m, d)
	m.NotFound(http.NotFound) // not found routes
}
