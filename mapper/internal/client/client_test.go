package client

import (
	"context"
	"errors"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/medmouine/mapper/pkg/device"
	"github.com/medmouine/mapper/pkg/device/temperature"
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
	c := NewClient[temperature.Data](new(MockDevice), opts)

	assert.NotNil(t, c)
	assert.Equal(t, opts.DataTopic, c.Options.DataTopic)
	assert.Equal(t, opts.StateTopics, c.Options.StateTopics)
	assert.Equal(t, opts.SubTopics, c.Options.SubTopics)
}

func TestClient_Connect(t *testing.T) {
	mockMQTT := new(MockMQTT)
	mockToken := new(MockToken)
	mockDevice := new(MockDevice)

	mockToken.On("Wait").Return(true)
	mockToken.On("Error").Return(nil)
	mockMQTT.On("Connect").Return(mockToken)
	mockMQTT.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(mockToken)

	c := NewClient[temperature.Data](mockDevice,
		&Options{
			MqttOptions:     MQTT.NewClientOptions(),
			SubTopics:       []string{"topic1", "topic2"},
			DataTopic:       "data",
			StateTopics:     []string{"state"},
			PublishInterval: 1 * time.Second,
		})
	c.mqtt = mockMQTT

	err := c.Connect()
	assert.Nil(t, err)
}

func TestClient_ConnectFailure(t *testing.T) {
	mockMQTT := new(MockMQTT)
	mockToken := new(MockToken)
	mockDevice := new(MockDevice)
	mockToken.On("Wait").Return(true)
	mockToken.On("Error").Return(errors.New("connection error"))
	mockMQTT.On("Connect").Return(mockToken)

	c := NewClient[temperature.Data](mockDevice,
		&Options{
			MqttOptions:     MQTT.NewClientOptions(),
			SubTopics:       []string{"topic1", "topic2"},
			DataTopic:       "data",
			StateTopics:     []string{"state"},
			PublishInterval: 1 * time.Second,
		})
	c.mqtt = mockMQTT

	err := c.Connect()
	assert.NotNil(t, err)
}

func TestClient_SubscribeFailure(t *testing.T) {
	mockMQTT := new(MockMQTT)
	mockToken := new(MockToken)
	mockDevice := new(MockDevice)
	mockToken.On("Wait").Return(true)
	mockToken.On("Error").Return(nil).Once()                    // For Connect()
	mockToken.On("Error").Return(errors.New("subscribe error")) // For Subscribe()
	mockMQTT.On("Connect").Return(mockToken)
	mockMQTT.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(mockToken)

	c := NewClient[temperature.Data](mockDevice,
		&Options{
			MqttOptions:     MQTT.NewClientOptions(),
			SubTopics:       []string{"topic1", "topic2"},
			DataTopic:       "data",
			StateTopics:     []string{"state"},
			PublishInterval: 1 * time.Second,
		})
	c.mqtt = mockMQTT

	err := c.Connect()
	assert.NotNil(t, err)
}

func TestClient_StreamData(t *testing.T) {
	mockMQTT := new(MockMQTT)
	mockToken := new(MockToken)
	mockDevice := new(MockDevice)
	mockToken.On("Wait").Return(true)
	mockToken.On("Error").Return(nil)
	mockMQTT.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockToken)
	mockDevice.On("Read").Return(temperature.Data{Temperature: 100, Humidity: 50})

	c := NewClient[temperature.Data](mockDevice,
		&Options{
			MqttOptions:     MQTT.NewClientOptions(),
			SubTopics:       []string{"topic1", "topic2"},
			DataTopic:       "data",
			StateTopics:     []string{"state"},
			PublishInterval: 4 * time.Millisecond,
		})
	c.mqtt = mockMQTT

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	go c.StreamData(ctx)()

	// give some time for go routine to execute
	time.Sleep(20 * time.Millisecond)

	mockMQTT.AssertNumberOfCalls(t, "Publish", 3)
}

func TestUpdateLocalState(t *testing.T) {
	mockDevice := new(MockDevice)
	client := NewClient[temperature.Data](mockDevice, &Options{
		MqttOptions:     MQTT.NewClientOptions(),
		SubTopics:       []string{"topic1", "topic2"},
		DataTopic:       "data",
		StateTopics:     []string{"state"},
		PublishInterval: 1 * time.Second,
	})

	// Define your payload here
	payload := []byte(`{"report_interval": "5s"}`)

	// Call the function
	client.UpdateLocalState(payload)

	// Assert the publishInterval has been updated correctly
	assert.Equal(t, 5*time.Second, client.Options.PublishInterval)
}

func TestHandle(t *testing.T) {
	mockDevice := new(MockDevice)
	c := NewClient[temperature.Data](mockDevice,
		&Options{
			MqttOptions:     MQTT.NewClientOptions(),
			SubTopics:       []string{"topic1", "topic2"},
			DataTopic:       "data",
			StateTopics:     []string{"state"},
			PublishInterval: 1 * time.Second,
		})
	c.Options.MqttOptions.ClientID = "test_client"

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
	assert.Equal(t, 5*time.Second, c.Options.PublishInterval)
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

type MockDevice struct {
	device.Device[temperature.Data]
	mock.Mock
}

func (m *MockDevice) Read() temperature.Data {
	args := m.Called()
	return args.Get(0).(temperature.Data)
}

func (m *MockDevice) Data() *temperature.Data {
	args := m.Called()
	return args.Get(0).(*temperature.Data)
}

func (m *MockDevice) Simulator() device.Simulator[temperature.Data] {
	args := m.Called()
	return args.Get(0).(device.Simulator[temperature.Data])
}
