package config

import (
	"time"

	"github.com/caarlos0/env/v9"
)

type ServerConfig struct {
	Addr         string        `env:"SERVER_ADDR" ,envDefault:"${SERVER_HOST}:${SERVER_PORT}"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT"`
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	return cfg
}
