package temperature_sensor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemperatureSimulator(t *testing.T) {
	simulator := NewTemperatureSimulator("id", -10, 50)

	// Verify that the simulator has been properly initialized
	assert.NotNil(t, simulator)
	assert.Equal(t, "id", simulator.Id)
	assert.Equal(t, None, simulator.Anomaly)
	assert.LessOrEqual(t, simulator.Temperature, 50.0)
	assert.GreaterOrEqual(t, simulator.Temperature, -10.0)
}

func TestIntroduceAnomaly(t *testing.T) {
	simulator := NewTemperatureSimulator("id", -10, 50)

	// Introduce Spike anomaly and test
	simulator.IntroduceAnomaly(Spike)
	assert.Equal(t, Spike, simulator.Anomaly)

	// Introduce Drift anomaly and test
	simulator.IntroduceAnomaly(Drift)
	assert.Equal(t, Drift, simulator.Anomaly)

	// Introduce Noise anomaly and test
	simulator.IntroduceAnomaly(Noise)
	assert.Equal(t, Noise, simulator.Anomaly)

	// Introduce Flatline anomaly and test
	simulator.IntroduceAnomaly(Flatline)
	assert.Equal(t, Flatline, simulator.Anomaly)

	// Return to None anomaly and test
	simulator.IntroduceAnomaly(None)
	assert.Equal(t, None, simulator.Anomaly)
}

// Since the Read() method involves randomness, we can't predict the exact value.
// But we can check if the value is within the expected range.
func TestSimulatorRead(t *testing.T) {
	simulator := NewTemperatureSimulator("id", -10, 50)

	// For None anomaly
	simulator.IntroduceAnomaly(None)
	temp := simulator.Read()
	assert.LessOrEqual(t, temp, 50.0)
	assert.GreaterOrEqual(t, temp, -10.0)

	// For Flatline anomaly
	simulator.IntroduceAnomaly(Flatline)
	tempFlat := simulator.Read()
	assert.Equal(t, tempFlat, simulator.Temperature)

	// For Spike anomaly
	simulator.IntroduceAnomaly(Spike)
	tempSpike := simulator.Read()
	assert.NotEqual(t, tempSpike, tempFlat)

	// For Drift anomaly
	simulator.IntroduceAnomaly(Drift)
	tempDrift := simulator.Read()
	assert.NotEqual(t, tempDrift, tempSpike)

	// For Noise anomaly
	simulator.IntroduceAnomaly(Noise)
	tempNoise := simulator.Read()
	assert.NotEqual(t, tempNoise, tempDrift)
}
