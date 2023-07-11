package config

import (
	"sync"
)

// Config struct for describe configuration of the app.
type Config struct {
	Server *ServerConfig
	Mqtt   *MqttConfig
}

var (
	once     sync.Once // create sync.Once primitive
	instance *Config   // create nil Config struct
)

// NewConfig function to prepare config variables from .env file and return config.
func NewConfig() *Config {
	once.Do(func() {
		mqttConfig, err := NewMqttConfig()
		serverConfig := NewServerConfig()
		if err != nil {
			panic(err)
		}
		instance = &Config{
			Server: serverConfig,
			Mqtt:   mqttConfig,
		}
	})
	// Return configured config instance.
	return instance
}
