package client

import (
	"context"
	"errors"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	device "github.com/medmouine/device-mapper/pkg/sensor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewClient(t *testing.T) {
	// Here we just test if the function creates a new client correctly.
	opts := &Options{
		MqttOptions:     MQTT.NewClientOptions(),
		SubTopics:       []string{"topic1", "topic2"},
		DataTopic:       "data",
		StateTopics:     []string{"state"},
		PublishInterval: 1 * time.Second,
	}
	c := NewClient(opts, new(MockTemperatureDriver))

	assert.NotNil(t, c)
	assert.Equal(t, opts.DataTopic, c.DataTopic)
	assert.Equal(t, opts.StateTopics, c.StateTopics)
	assert.Equal(t, Init, c.Status)
}

func TestClient_Connect(t *testing.T) {
	mockMQTT := new(MockMQTT)
	mockToken := new(MockToken)
	mockToken.On("Wait").Return(true)
	mockToken.On("Error").Return(nil)
	mockMQTT.On("Connect").Return(mockToken)
	mockMQTT.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(mockToken)

	c := NewClient(&Options{
		MqttOptions:     MQTT.NewClientOptions(),
		SubTopics:       []string{"topic1", "topic2"},
		DataTopic:       "data",
		StateTopics:     []string{"state"},
		PublishInterval: 1 * time.Second,
	}, &device.TemperatureSimulator{})
	c.mqtt = mockMQTT

	err := c.Connect()
	assert.Nil(t, err)
	assert.Equal(t, Connected, c.Status)
}

func TestClient_ConnectFailure(t *testing.T) {
	mockMQTT := new(MockMQTT)
	mockToken := new(MockToken)
	mockToken.On("Wait").Return(true)
	mockToken.On("Error").Return(errors.New("connection error"))
	mockMQTT.On("Connect").Return(mockToken)

	c := NewClient(&Options{
		MqttOptions:     MQTT.NewClientOptions(),
		SubTopics:       []string{"topic1", "topic2"},
		DataTopic:       "data",
		StateTopics:     []string{"state"},
		PublishInterval: 1 * time.Second,
	}, &device.TemperatureSimulator{})
	c.mqtt = mockMQTT

	err := c.Connect()
	assert.NotNil(t, err)
	assert.Equal(t, ConnError, c.Status)
}

func TestClient_SubscribeFailure(t *testing.T) {
	mockMQTT := new(MockMQTT)
	mockToken := new(MockToken)
	mockToken.On("Wait").Return(true)
	mockToken.On("Error").Return(nil).Once()                    // For Connect()
	mockToken.On("Error").Return(errors.New("subscribe error")) // For Subscribe()
	mockMQTT.On("Connect").Return(mockToken)
	mockMQTT.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(mockToken)

	c := NewClient(&Options{
		MqttOptions:     MQTT.NewClientOptions(),
		SubTopics:       []string{"topic1", "topic2"},
		DataTopic:       "data",
		StateTopics:     []string{"state"},
		PublishInterval: 1 * time.Second,
	}, &device.TemperatureSimulator{})
	c.mqtt = mockMQTT

	err := c.Connect()
	assert.NotNil(t, err)
	assert.Equal(t, Connected, c.Status)
}

func TestClient_StreamData(t *testing.T) {
	mockMQTT := new(MockMQTT)
	mockToken := new(MockToken)
	mockTemperatureDriver := new(MockTemperatureDriver)
	mockToken.On("Wait").Return(true)
	mockToken.On("Error").Return(nil)
	mockMQTT.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockToken)
	mockTemperatureDriver.On("Read").Return(100.0)

	c := NewClient(&Options{
		MqttOptions:     MQTT.NewClientOptions(),
		SubTopics:       []string{"topic1", "topic2"},
		DataTopic:       "data",
		StateTopics:     []string{"state"},
		PublishInterval: 4 * time.Millisecond,
	}, mockTemperatureDriver)
	c.mqtt = mockMQTT

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	go c.StreamData(ctx)()

	// give some time for go routine to execute
	time.Sleep(20 * time.Millisecond)

	mockMQTT.AssertNumberOfCalls(t, "Publish", 3)
}

func TestUpdateLocalState(t *testing.T) {
	client := NewClient(&Options{
		MqttOptions:     MQTT.NewClientOptions(),
		SubTopics:       []string{"topic1", "topic2"},
		DataTopic:       "data",
		StateTopics:     []string{"state"},
		PublishInterval: 1 * time.Second,
	}, &device.TemperatureSimulator{})

	// Define your payload here
	payload := []byte(`{"report_interval": "5s"}`)

	// Call the function
	client.UpdateLocalState(payload)

	// Assert the publishInterval has been updated correctly
	assert.Equal(t, 5*time.Second, client.publishInterval)
}

func TestHandle(t *testing.T) {
	c := NewClient(&Options{
		MqttOptions:     MQTT.NewClientOptions(),
		SubTopics:       []string{"topic1", "topic2"},
		DataTopic:       "data",
		StateTopics:     []string{"state"},
		PublishInterval: 1 * time.Second,
	}, &device.TemperatureSimulator{})
	c.opts.MqttOptions.ClientID = "test_client"

	// Mock MQTT client and message
	c.mqtt = new(MockMQTT)
	mqttMessage := new(MockMessage)

	// Define your payload here
	payload := []byte(`{"report_interval": "5s"}`)

	// Assume MQTT message topic is the id
	mqttMessage.On("Topic").Return(UpdateStateTopic.Fmt("test_client"))
	mqttMessage.On("Payload").Return(payload)

	// Call the function
	handleFunc := c.handle()
	handleFunc(c.mqtt, mqttMessage)

	// Assert the publishInterval has been updated correctly
	assert.Equal(t, 5*time.Second, c.publishInterval)
}

// MockMQTT is a mock struct that satisfies the MQTT.Client interface
type MockMQTT struct {
	MQTT.Client
	mock.Mock
}

func (m *MockMQTT) Connect() MQTT.Token {
	args := m.Called()
	return args.Get(0).(MQTT.Token)
}

func (m *MockMQTT) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	args := m.Called(topic, qos, retained, payload)
	return args.Get(0).(MQTT.Token)
}

func (m *MockMQTT) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	args := m.Called(topic, qos, callback)
	return args.Get(0).(MQTT.Token)
}

type MockToken struct {
	MQTT.Token
	mock.Mock
}

type MockMessage struct {
	mock.Mock
	MQTT.Message
}

func (m *MockMessage) Payload() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *MockMessage) Topic() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockToken) Wait() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockToken) Error() error {
	args := m.Called()
	return args.Error(0)
}

type MockTemperatureDriver struct {
	device.Sensor
	mock.Mock
}

func (m *MockTemperatureDriver) Read() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}
