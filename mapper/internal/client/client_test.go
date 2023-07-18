package client

import (
	"context"
	"errors"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/medmouine/mapper/pkg/device"
	"github.com/medmouine/mapper/pkg/device/simulation"
	"github.com/medmouine/mapper/pkg/device/temperature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewClient(t *testing.T) {
	// Here we just test if the function creates a new client correctly.
	opts := &Options{
		MqttOptions: MQTT.NewClientOptions(),
		SubTopics:   []string{"topic1", "topic2"},
		DataTopic:   "data",
		StateTopics: []string{"state"},
	}
	c := NewClient[temperature.TemperatureData](new(MockDevice), opts)

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

	c := NewClient[temperature.TemperatureData](mockDevice,
		&Options{
			MqttOptions: MQTT.NewClientOptions(),
			SubTopics:   []string{"topic1", "topic2"},
			DataTopic:   "data",
			StateTopics: []string{"state"},
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

	c := NewClient[temperature.TemperatureData](mockDevice,
		&Options{
			MqttOptions: MQTT.NewClientOptions(),
			SubTopics:   []string{"topic1", "topic2"},
			DataTopic:   "data",
			StateTopics: []string{"state"},
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

	c := NewClient[temperature.TemperatureData](mockDevice,
		&Options{
			MqttOptions: MQTT.NewClientOptions(),
			SubTopics:   []string{"topic1", "topic2"},
			DataTopic:   "data",
			StateTopics: []string{"state"},
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
	mockDevice.On("Read").Return(temperature.TemperatureData{Temperature: 100, Humidity: 50})
	mockDevice.On("PublishInterval").Return(4 * time.Millisecond)
	mockDevice.On("ID").Return("dev1")
	c := NewClient[temperature.TemperatureData](mockDevice,
		&Options{
			MqttOptions: MQTT.NewClientOptions(),
			SubTopics:   []string{"topic1", "topic2"},
			DataTopic:   "data/+/topic",
			StateTopics: []string{"state"},
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
	mockDevice.On("PublishInterval").Return(time.Duration(1))
	mockDevice.On("SetPublishInterval", mock.Anything).Return(nil)
	client := NewClient[temperature.TemperatureData](mockDevice, &Options{
		MqttOptions: MQTT.NewClientOptions(),
		SubTopics:   []string{"topic1", "topic2"},
		DataTopic:   "data",
		StateTopics: []string{"state"},
	})

	// Define your payload here
	payload := []byte(`{"report_interval": "5s"}`)

	// Call the function
	client.UpdateLocalState(payload)

	// Assert the publishInterval has been updated correctly
	mockDevice.AssertCalled(t, "SetPublishInterval", 5*time.Second)
}

func TestHandle(t *testing.T) {
	mockDevice := new(MockDevice)
	mockDevice.On("PublishInterval").Return(time.Duration(1))
	mockDevice.On("SetPublishInterval", mock.Anything).Return(nil)
	c := NewClient[temperature.TemperatureData](mockDevice,
		&Options{
			MqttOptions: MQTT.NewClientOptions(),
			SubTopics:   []string{"topic1", "topic2"},
			DataTopic:   "data",
			StateTopics: []string{"state"},
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

	// Assert the publishInterval has been updated correctly'
	mockDevice.AssertCalled(t, "SetPublishInterval", 5*time.Second)
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
	device.Device[temperature.TemperatureData]
	mock.Mock
}

func (m *MockDevice) ID() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockDevice) SetPublishInterval(pi time.Duration) {
	m.Called(pi)
}

func (m *MockDevice) PublishInterval() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockDevice) Read() temperature.TemperatureData {
	args := m.Called()
	return args.Get(0).(temperature.TemperatureData)
}

func (m *MockDevice) Data() *temperature.TemperatureData {
	args := m.Called()
	return args.Get(0).(*temperature.TemperatureData)
}

func (m *MockDevice) Simulator() simulation.Simulator[simulation.VarSimulationConfig] {
	args := m.Called()
	return args.Get(0).(simulation.Simulator[simulation.VarSimulationConfig])
}
