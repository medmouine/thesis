package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/medmouine/device-mapper/internal/router/healthcheck"
)

// GetRoutes function for getting routes.
func GetRoutes(m *chi.Mux) {
	healthcheck.Routes(m)     // health check routes
	m.NotFound(http.NotFound) // not found routes
}
