package client

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-multierror"
	"github.com/koddr/gosl"
	"github.com/medmouine/mapper/pkg/device"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/stream"
)

type Client[T interface{}] struct {
	mqtt    MQTT.Client
	d       device.Device[T]
	Options *Options
	mux     sync.Mutex
}

type Options struct {
	MqttOptions *MQTT.ClientOptions
	SubTopics   []string
	DataTopic   string
	StateTopics []string
}

func NewClient[T interface{}](d device.Device[T], opts *Options) *Client[T] {
	clt := MQTT.NewClient(opts.MqttOptions)

	return &Client[T]{
		mqtt:    clt,
		d:       d,
		Options: opts,
	}
}

func (c *Client[T]) Connect() error {
	if token := c.mqtt.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("Error connecting to MQTT broker: %v", token.Error())
		return token.Error()
	}

	log.Infof("Connected to MQTT brokers [%v]", c.Options.MqttOptions.Servers)

	if err := c.Subscribe(); err != nil {
		log.Errorf("Error subscribing to topics: %v", err)
		return err
	}
	return nil
}

func (c *Client[T]) Subscribe() error {
	var errs error
	lo.ForEach(lo.Union(c.Options.SubTopics, c.Options.StateTopics), func(topic string, i int) {
		if token := c.mqtt.Subscribe(topic, 0, c.handle()); token.Wait() && token.Error() != nil {
			errs = multierror.Append(errs, token.Error())
			return
		}
		log.Infof("Subscribed to topic [%v]", topic)
	})
	return errs
}

func (c *Client[T]) StreamData(ctx context.Context) func() {
	var t = c.Options.DataTopic
	t = strings.ReplaceAll(t, "+", c.d.ID())
	pubTask := func() stream.Callback {
		c.mux.Lock()
		d := c.d.Read()
		c.mux.Unlock()
		payload, err := gosl.Marshal(&d)
		log.Infof("Publishing data [%v] to topic [%v]", string(payload), t)
		if err != nil {
			return c.handlePublishError(fmt.Errorf("error during data marshal: %w", err), t, d)
		}
		if token := c.mqtt.Publish(t, 0, false, payload); token.Wait() && token.Error() != nil {
			return c.handlePublishError(token.Error(), t, string(payload))
		}
		return func() {
			log.Infof("Successfully published data [%v] to topic [%v]", string(payload), t)
		}
	}
	upstream := stream.New()
	return func() {
		log.Infof("Starting data stream on topic [%v]", t)

		for {
			select {
			case <-ctx.Done():
				log.Infof("Stopping data stream")
				upstream.Wait()
				return
			default:
				time.Sleep(c.d.PublishInterval())
				upstream.Go(pubTask)
			}
		}
	}
}

func (c *Client[T]) UpdateLocalState(payload []byte) {
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
	c.d.SetPublishInterval(d)
}

func (c *Client[T]) handle() func(MQTT.Client, MQTT.Message) {
	var t = c.Options.StateTopics[0]
	t = strings.ReplaceAll(t, "+", c.d.ID())

	return func(clt MQTT.Client, msg MQTT.Message) {
		log.Infof("Received message [%v] on topic [%v]", string(msg.Payload()), msg.Topic())
		switch msg.Topic() {
		case t:
			log.Infof("Received state update message [%v] on topic [%v]", string(msg.Payload()), msg.Topic())
			c.UpdateLocalState(msg.Payload())
		default:
			log.Infof("Received message [%v] on topic [%v]", string(msg.Payload()), msg.Topic())
		}
	}
}

func (c *Client[T]) handlePublishError(err error, topic string, payload interface{}) func() {
	return func() {
		log.Errorf("Error publishing payload %v to topic [%s]: %v", err, topic, payload)
	}
}
