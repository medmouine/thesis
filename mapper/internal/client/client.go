package client

import (
	"context"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-multierror"
	"github.com/koddr/gosl"
	sensor "github.com/medmouine/device-mapper/pkg/sensor"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/stream"
)

type Client struct {
	mqtt            MQTT.Client
	Subs            []string
	DataTopic       string
	StateTopics     []string
	Status          Status
	driver          sensor.Sensor
	opts            *Options
	publishInterval time.Duration
	mux             sync.Mutex
}

type Options struct {
	MqttOptions     *MQTT.ClientOptions
	SubTopics       []string
	DataTopic       string
	StateTopics     []string
	PublishInterval time.Duration
}

func NewClient(opts *Options, driver sensor.Sensor) *Client {
	clt := MQTT.NewClient(opts.MqttOptions)

	return &Client{
		clt,
		opts.SubTopics,
		opts.DataTopic,
		opts.StateTopics,
		Init,
		driver,
		opts,
		opts.PublishInterval,
		sync.Mutex{},
	}
}

func (c *Client) SetStatus(state Status) {
	c.Status = state
}

// func handleError(i interface{}) {}

func (c *Client) Connect() error {
	c.SetStatus(Connecting)
	if token := c.mqtt.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("Error connecting to MQTT broker: %v", token.Error())
		c.SetStatus(ConnError)
		return token.Error()
	}

	log.Infof("Connected to MQTT brokers [%v]", c.opts.MqttOptions.Servers)
	c.SetStatus(Connected)

	if err := c.Subscribe(); err != nil {
		log.Errorf("Error subscribing to topics: %v", err)
		return err
	}
	return nil
}

func (c *Client) Subscribe() error {
	var errs error
	lo.ForEach(lo.Union(c.Subs, c.StateTopics), func(topic string, i int) {
		if token := c.mqtt.Subscribe(topic, 0, c.handle()); token.Wait() && token.Error() != nil {
			errs = multierror.Append(errs, token.Error())
			return
		}
		log.Infof("Subscribed to topic [%v]", topic)
	})
	return errs
}

func (c *Client) StreamData(ctx context.Context) func() {
	pubTask := func() stream.Callback {
		c.mux.Lock()
		d := c.driver.Read()
		if token := c.mqtt.Publish(c.DataTopic, 0, false, d); token.Wait() && token.Error() != nil {
			log.Errorf("Error publishing data: %v", token.Error())
		}
		defer c.mux.Unlock()
		return func() {
			log.Infof("Published data [%v] to topic [%v]", d, c.DataTopic)
		}
	}
	upstream := stream.New()
	return func() {
		for {
			select {
			case <-ctx.Done():
				log.Infof("Stopping data stream")
				upstream.Wait()
				return
			default:
				time.Sleep(c.publishInterval)
				upstream.Go(pubTask)
			}
		}
	}
}

func (c *Client) UpdateLocalState(payload []byte) {
	s, err := gosl.Unmarshal(payload, &StateUpdate{})
	if err != nil {
		log.Errorf("Error unmarshalling state payload: %v", err)
		return
	}
	log.Infof("Received new state: %v", s)

	d, err := time.ParseDuration(s.ReportInterval)
	if err != nil {
		log.Warnf("Received invalid report interval duration: %v", s.ReportInterval)
		return
	}
	c.mux.Lock()
	c.publishInterval = d
	c.mux.Unlock()
}

func (c *Client) handle() func(MQTT.Client, MQTT.Message) {
	return func(clt MQTT.Client, msg MQTT.Message) {
		id := c.opts.MqttOptions.ClientID
		log.Infof("Received message [%v] on topic [%v]", string(msg.Payload()), msg.Topic())
		switch msg.Topic() {
		case UpdateStateTopic.Fmt(id):
			log.Infof("Received state update message [%v] on topic [%v]", string(msg.Payload()), msg.Topic())
			c.UpdateLocalState(msg.Payload())
		default:
			log.Infof("Received message [%v] on topic [%v]", string(msg.Payload()), msg.Topic())
		}
	}
}
