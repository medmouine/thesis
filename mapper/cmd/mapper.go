package cmd

import (
	"context"
	"net/http"

	"github.com/medmouine/device-mapper/internal/client"
	"github.com/medmouine/device-mapper/internal/config"
	sensor "github.com/medmouine/device-mapper/pkg/sensor"
)

type Mapper struct {
	Config       *config.Config
	Client       *client.Client
	DeviceDriver sensor.Sensor
	API          http.Handler
}

// Run function to start server instance with config & router.
func (m *Mapper) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := getHTTPServer(m.Config, m.API)

	if err := m.Client.Connect(); err != nil {
		return err
	}

	go m.Client.StreamData(ctx)

	return server.ListenAndServe()
}

func getHTTPServer(c *config.Config, r http.Handler) *http.Server {
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
