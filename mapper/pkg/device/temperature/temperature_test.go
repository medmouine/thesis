package temperature

import (
	"testing"

	"github.com/medmouine/mapper/pkg/device"
	"github.com/stretchr/testify/assert"
)

func TestNewTemperatureSimulator(t *testing.T) {
	tsim := NewTemperatureSimulator("id", -10, 50)

	// Verify that the tsim has been properly initialized
	assert.NotNil(t, tsim)
	assert.Equal(t, "id", tsim.ID())
	assert.Equal(t, device.None, *tsim.Anomaly())
	assert.LessOrEqual(t, tsim.Data().Temperature, 50.0)
	assert.GreaterOrEqual(t, tsim.Data().Temperature, -10.0)
	assert.NotNil(t, tsim.Data().Humidity)
}

func TestIntroduceAnomaly(t *testing.T) {
	simulator := NewTemperatureSimulator("id", -10, 50)

	// Introduce Spike anomaly and test
	simulator.IntroduceAnomaly(device.Spike)
	assert.Equal(t, device.Spike, *simulator.Anomaly())

	// Introduce Drift anomaly and test
	simulator.IntroduceAnomaly(device.Drift)
	assert.Equal(t, device.Drift, *simulator.Anomaly())

	// Introduce Noise anomaly and test
	simulator.IntroduceAnomaly(device.Noise)
	assert.Equal(t, device.Noise, *simulator.Anomaly())

	// Introduce Flatline anomaly and test
	simulator.IntroduceAnomaly(device.Flatline)
	assert.Equal(t, device.Flatline, *simulator.Anomaly())

	// Return to None anomaly and test
	simulator.IntroduceAnomaly(device.None)
	assert.Equal(t, device.None, *simulator.Anomaly())
}

// Since the Read() method involves randomness, we can't predict the exact value.
// But we can check if the value is within the expected range.
func TestSimulatorRead(t *testing.T) {
	simulator := NewTemperatureSimulator("id", -10, 50)

	// For None anomaly
	simulator.IntroduceAnomaly(device.None)
	temp := simulator.Read()
	assert.LessOrEqual(t, temp.Temperature, 50.0)
	assert.GreaterOrEqual(t, temp.Temperature, -10.0)

	// For Flatline anomaly
	simulator.IntroduceAnomaly(device.Flatline)
	tempFlat := simulator.Read()
	assert.Equal(t, tempFlat, *simulator.Data())

	// For device.Spike anomaly
	simulator.IntroduceAnomaly(device.Spike)
	tempSpike := simulator.Read()
	assert.NotEqual(t, tempSpike, tempFlat)

	// For Drift anomaly
	simulator.IntroduceAnomaly(device.Drift)
	tempDrift := simulator.Read()
	assert.NotEqual(t, tempDrift, tempSpike)

	// For Noise anomaly
	simulator.IntroduceAnomaly(device.Noise)
	tempNoise := simulator.Read()
	assert.NotEqual(t, tempNoise, tempDrift)
}

// Since the Read() method involves randomness, we can't predict the exact value.
// But we can check if the value is within the expected range.
func TestGetSimulator(t *testing.T) {
	d := NewTemperatureSimulator("id", -10, 50).(TemperatureDevice)
	s := d.Simulator()
	assert.NotNil(t, s)
	assert.Equal(t, d, s)
}

//
//func TestNewTemperatureSensor(t *testing.T) {
//	baseDevice := newBaseTemperatureDevice("id", -10, 50)
//
//	assert.Equal(t, "id", baseDevice.ID)
//	assert.Equal(t, -10.0, baseDevice.MinTemp)
//	assert.Equal(t, 50.0, baseDevice.maxTemp)
//	assert.Equal(t, 20.0, baseDevice.Data)
//}
//
//func TestRead(t *testing.T) {
//	baseDevice := newTemperatureSensor("Id", -10, 50, 20)
//
//	// The Read() method should panic with "not implemented" message
//	assert.PanicsWithValue(t, "not implemented", func() { baseDevice.Read() }, "The code did not panic")
//
//	// or you can use require instead of assert to stop the test if it fails
//	require.PanicsWithValue(t, "not implemented", func() { baseDevice.Read() }, "The code did not panic")
//}
