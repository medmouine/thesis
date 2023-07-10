package config

import (
	"net/url"
	"os"

	"github.com/medmouine/device-mapper/internal/mqtt"
	"github.com/samber/lo"
)

type MqttConfig struct {
	ClientID    string
	SubTopics   []string
	DataTopics  []string
	StateTopics []string
	BrokerURL   url.URL
}

func NewMqttConfig() *MqttConfig {
	DataTopics := []string{"data"}
	StateTopics := []string{"state"}
	SubTocis := []string{"sub"}

	id := os.Getenv("MQTT_CLIENT_ID")
	if lo.IsEmpty(id) {
		panic("wrong client id (check your .env)")
	}
	brokerURL, err := url.Parse(os.Getenv("MQTT_BROKER_URL"))
	if err != nil {
		panic("wrong broker url (check your .env)")
	}

	return &MqttConfig{
		ClientID:    id,
		SubTopics:   SubTocis,
		DataTopics:  DataTopics,
		StateTopics: StateTopics,
		BrokerURL:   *brokerURL,
	}
}

func (c *MqttConfig) ToClientOptions() *mqtt.ClientOptions {
	mqttOpts := &mqtt.ClientOptions{
		SubTopics:   c.SubTopics,
		DataTopics:  c.DataTopics,
		StateTopics: c.StateTopics,
	}
	mqttOpts.SetClientID(c.ClientID)
	mqttOpts.AddBroker(c.BrokerURL.Host)
	return mqttOpts
}
