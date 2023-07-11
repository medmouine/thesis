package mqtt

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"github.com/sourcegraph/conc/stream"
)

type ClientStatus string

const (
	INIT       ClientStatus = "init"
	CONNECTING ClientStatus = "connecting"
	CONNECTED  ClientStatus = "connected"
)

type SubTopic struct {
	Handler     func(MQTT.Client, MQTT.Message)
	Topic       string
	initialised bool
}

type Client struct {
	mqttClient  MQTT.Client
	Subs        []SubTopic
	DataTopic   string
	StateTopics []string
	Status      ClientStatus
	driver      interface{}
}

type ClientOptions struct {
	*MQTT.ClientOptions
	SubTopics   []string
	DataTopic   string
	StateTopics []string
}

func NewClient(opts *ClientOptions) (*Client, error) {
	clt := MQTT.NewClient(opts.ClientOptions)
	subs := lo.Map(opts.SubTopics, func(topic string, _ int) SubTopic {
		return SubTopic{
			Topic:       topic,
			initialised: false,
		}
	})

	return &Client{
		clt,
		subs,
		opts.DataTopic,
		opts.StateTopics,
		INIT,
		nil,
	}, nil
}

func (c *Client) Init() error {
	p := pool.New().WithErrors()
	p.Go(c.conn)
	RETRIES := 3
	for i := 0; i < RETRIES; i++ {
		if err := p.Wait(); err != nil {
			fmt.Printf("connection error, retry after 5s: %v\n", err)
			time.Sleep(5 * time.Second)
		}
	}
	if c.Status != CONNECTED {
		return fmt.Errorf("connection error")
	}

	p.Go(c.Subscribe)
	err := p.Wait()
	if err != nil {
		if !lo.NoneBy(c.Subs, func(sub SubTopic) bool { return sub.initialised }) {
			err = multierror.Prefix(err, "some subs failed to initialise")
		}
		return err
	}

	//pubStream := stream.New()
	//go func() {
	//pubStream.Go()
	//}()
	return nil
}

func (c *Client) SetStatus(state ClientStatus) {
	c.Status = state
}

// func handleError(i interface{}) {}

func (c *Client) conn() error {
	c.SetStatus(CONNECTING)
	token := c.mqttClient.Connect()
	token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error())
	}
	fmt.Println("connect success")
	c.SetStatus(CONNECTED)
	return token.Error()
}

func (c *Client) Publish(topic string) stream.Callback {
	return func() {
		var errs error
		token := c.mqttClient.Publish(topic, 0, false, c.driver)
		token.Wait()
		errs = multierror.Append(errs, token.Error())
		fmt.Printf("publish error: %v", errs)
	}
}

func (c *Client) Subscribe() error {
	var errs error
	for _, topic := range c.Subs {
		token := c.mqttClient.Subscribe(topic.Topic, 0, handle)
		token.Wait()
		errs = multierror.Append(errs, token.Error())
		if token.Error() == nil {
			el, i, _ := lo.FindIndexOf(c.Subs, func(sub SubTopic) bool {
				return sub.Topic == topic.Topic
			})
			if i == -1 {
				el.initialised = true
				el.Handler = handle
			}
		}
	}
	return errs
}

func handle(client MQTT.Client, msg MQTT.Message) {}
