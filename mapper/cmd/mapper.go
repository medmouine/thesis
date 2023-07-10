package cmd

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/medmouine/device-mapper/internal/config"
	"github.com/medmouine/device-mapper/internal/mqtt"
)

// Run function to start server instance with config & router.
func Run(c *config.Config, r *chi.Mux, mqttClient *mqtt.Client) error {
	server := getHttpServer(c, r)

	if err := mqttClient.Init(); err != nil {
		return err
	}
	// Start server.
	return server.ListenAndServe()
}

func getHttpServer(c *config.Config, r *chi.Mux) *http.Server {
	// Prepare server with CloudFlare recommendation timeouts config.
	// See: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := &http.Server{
		Handler:      r,
		Addr:         c.Server.Addr,
		ReadTimeout:  c.Server.ReadTimeout,
		WriteTimeout: c.Server.WriteTimeout,
		IdleTimeout:  c.Server.IdleTimeout,
	}
	return server
}
