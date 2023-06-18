package device_simulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemperatureSensor(t *testing.T) {
	sensor := NewTemperatureSensor(-10, 50)

	assert.InDelta(t, 20.0, sensor.Temperature, 30.0) // should be within -10 and 50
	assert.Equal(t, None, sensor.Anomaly)             // should be None by default
}

func TestReadNormal(t *testing.T) {
	sensor := NewTemperatureSensor(0, 100)
	sensor.Anomaly = None

	first := sensor.Read()
	second := sensor.Read()

	// Normal readings should not change more than 1.5 degrees
	assert.InDelta(t, first, second, 1.5)
}

func TestReadSpike(t *testing.T) {
	sensor := NewTemperatureSensor(0, 100)
	sensor.Anomaly = Spike

	first := sensor.Read()
	second := sensor.Read()

	// Spiked readings should change up to 15.0 degrees
	assert.True(t, (second <= first+15.0 && second >= first-15.0))
}

func TestReadFlatline(t *testing.T) {
	sensor := NewTemperatureSensor(0, 100)
	sensor.Anomaly = Flatline

	first := sensor.Read()
	second := sensor.Read()

	// Flatlined readings should not change
	assert.Equal(t, first, second)
}

func TestIntroduceAnomaly(t *testing.T) {
	sensor := NewTemperatureSensor(0, 100)
	sensor.IntroduceAnomaly(Spike)

	assert.Equal(t, Spike, sensor.Anomaly)
}
