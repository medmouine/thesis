package temperature

import (
	"testing"
	"time"

	"github.com/medmouine/mapper/pkg/device/simulation"
	"github.com/stretchr/testify/assert"
)

func TestNewTemperatureSimulator(t *testing.T) {
	tsim := NewTemperatureSimulator("id", time.Duration(time.Second), -10, 50)

	// Verify that the tsim has been properly initialized
	assert.NotNil(t, tsim)
	assert.Equal(t, "id", tsim.ID())
	assert.Equal(t, simulation.None, *tsim.Anomaly())
	assert.LessOrEqual(t, tsim.Data().Temperature, 50.0)
	assert.GreaterOrEqual(t, tsim.Data().Temperature, -10.0)
	assert.NotNil(t, tsim.Data().Humidity)
}

func TestIntroduceAnomaly(t *testing.T) {
	simulator := NewTemperatureSimulator("id", time.Duration(time.Second), -10, 50)

	// Introduce Spike anomaly and test
	simulator.IntroduceAnomaly(simulation.Spike)
	assert.Equal(t, simulation.Spike, *simulator.Anomaly())

	// Introduce Drift anomaly and test
	simulator.IntroduceAnomaly(simulation.Drift)
	assert.Equal(t, simulation.Drift, *simulator.Anomaly())

	// Introduce Noise anomaly and test
	simulator.IntroduceAnomaly(simulation.Noise)
	assert.Equal(t, simulation.Noise, *simulator.Anomaly())

	// Introduce Flatline anomaly and test
	simulator.IntroduceAnomaly(simulation.Flatline)
	assert.Equal(t, simulation.Flatline, *simulator.Anomaly())

	// Return to None anomaly and test
	simulator.IntroduceAnomaly(simulation.None)
	assert.Equal(t, simulation.None, *simulator.Anomaly())
}

// Since the Read() method involves randomness, we can't predict the exact value.
// But we can check if the value is within the expected range.
func TestSimulatorRead(t *testing.T) {
	simulator := NewTemperatureSimulator("id", time.Second, -10, 50)

	// For None anomaly
	simulator.IntroduceAnomaly(simulation.None)
	temp := simulator.Read()
	assert.LessOrEqual(t, temp.Temperature, 50.0)
	assert.GreaterOrEqual(t, temp.Temperature, -10.0)

	// For Flatline anomaly
	simulator.IntroduceAnomaly(simulation.Flatline)
	tempFlat := simulator.Read()
	assert.Equal(t, tempFlat, *simulator.Data())

	// For device.Spike anomaly
	simulator.IntroduceAnomaly(simulation.Spike)
	tempSpike := simulator.Read()
	assert.NotEqual(t, tempSpike, tempFlat)

	// For Drift anomaly
	simulator.IntroduceAnomaly(simulation.Drift)
	tempDrift := simulator.Read()
	assert.NotEqual(t, tempDrift, tempSpike)

	// For Noise anomaly
	simulator.IntroduceAnomaly(simulation.Noise)
	tempNoise := simulator.Read()
	assert.NotEqual(t, tempNoise, tempDrift)
}

// Since the Read() method involves randomness, we can't predict the exact value.
// But we can check if the value is within the expected range.
func TestGetSimulator(t *testing.T) {
	d := NewTemperatureSimulator("id", time.Duration(1*time.Second), -10, 50).(TemperatureDevice)
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
//	assert.Equal(t, 20.0, baseDevice.TemperatureData)
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
