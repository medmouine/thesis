package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/eclipse/paho.golang/autopaho"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/medmouine/mapper/internal/client"
)

type MqttConfig struct {
	ClientID        string        `env:"MQTT_CLIENT_ID"`
	SubTopics       []string      `env:"MQTT_SUB_TOPICS" envSeparator:":"`
	StateTopics     []string      `env:"MQTT_STATE_TOPICS" envSeparator:":"`
	DataTopic       string        `env:"MQTT_DATA_TOPIC"`
	BrokerURL       string        `env:"MQTT_BROKER_URL"`
	PublishInterval time.Duration `env:"MQTT_PUBLISH_INTERVAL" envDefault:"5s"`
	MaxTemp         float64       `env:"MAX_TEMP" envDefault:"100"`
	MinTemp         float64       `env:"MIN_TEMP" envDefault:"10"`
}

func NewMqttConfig() (*MqttConfig, error) {
	cfg := &MqttConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse mqtt config: %w", err)
	}

	return cfg, nil
}

func (c *MqttConfig) ToClientOptions() *client.Options {
	mqttOpts := &client.Options{
		MqttOptions: MQTT.NewClientOptions(),
		SubTopics:   c.SubTopics,
		DataTopic:   c.DataTopic,
		StateTopics: c.StateTopics,
	}
	mqttOpts.MqttOptions.SetClientID(c.ClientID)
	mqttOpts.MqttOptions.AddBroker(c.BrokerURL)
	mqttOpts.MqttOptions.SetKeepAlive(60 * time.Second)
	mqttOpts.MqttOptions.SetPingTimeout(1 * time.Second)

	return mqttOpts
}

func (c *MqttConfig) ToV5Config() autopaho.ClientConfig {
	return autopaho.ClientConfig{}
}
