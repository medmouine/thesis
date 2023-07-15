package simulation

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDevice struct{}

func (m *MockDevice) Config() *VarSimulationConfig {
	return nil
}

func (m *MockDevice) SetConfig(_ VarSimulationConfig) {}

func (m *MockDevice) IntroduceAnomaly(_ Anomaly) {}
func (m *MockDevice) Anomaly() *Anomaly          { return nil }

func TestBaseSimulator_IntroduceAnomaly(t *testing.T) {
	mockDevice := &MockDevice{}
	bs := &BaseSimulator{
		Simulation: mockDevice,
		SimConfig:  DefaultVarSimConfig(),
		mux:        sync.Mutex{},
	}

	bs.IntroduceAnomaly(Spike)

	assert.Equal(t, Spike, *bs.Anomaly())
}

func TestBaseSimulator_Anomaly(t *testing.T) {
	mockDevice := &MockDevice{}
	bs := &BaseSimulator{
		Simulation: mockDevice,
		SimConfig:  DefaultVarSimConfig(),
		mux:        sync.Mutex{},
	}

	anomaly := bs.Anomaly()

	assert.Equal(t, None, *anomaly)

	bs.IntroduceAnomaly(Flatline)

	anomaly = bs.Anomaly()

	assert.Equal(t, Flatline, *anomaly)
}
