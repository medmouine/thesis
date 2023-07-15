package temperature

import (
	"math/rand"
	"sync"
	"time"

	"github.com/medmouine/mapper/pkg/device/simulation"
)

type TemperatureSimulator interface {
	simulation.Simulator[simulation.VarSimulationConfig]
	TemperatureDevice
}

func NewTemperatureSimulator(id string, pi time.Duration, minTemp, maxTemp float64) TemperatureSimulator {
	var ts = new(temperatureSimulator)
	ts.temperatureDevice = newBaseTemperatureDevice(id, pi, minTemp, maxTemp)
	ts.BaseSimulator = &simulation.BaseSimulator{
		SimConfig:  simulation.DefaultVarSimConfig(),
		Simulation: ts,
	}
	ts.Device = ts
	return ts
}

type temperatureSimulator struct {
	*temperatureDevice
	*simulation.BaseSimulator
	mux sync.Mutex
}

func (ts *temperatureSimulator) Read() TemperatureData {
	ts.mux.Lock()
	defer ts.mux.Unlock()
	a := *ts.Anomaly()
	current := *ts.Data()
	t := ts.computeMinMax(current.Temperature, ts.computeVar(a))
	h := current.Humidity + ts.computeVar(a)
	d := NewData(t, h)
	return ts.SetData(d)
}

func (ts *temperatureSimulator) computeVar(a simulation.Anomaly) float64 {
	c := ts.Config()
	change := (rand.Float64() - 0.5) * c.GlobalVariance
	switch {
	default:
		return change
	case a == simulation.Flatline:
		return 0
	case a == simulation.Spike:
		return change * c.SpikeVariance
	case a == simulation.Drift:
		return change + c.DriftVariance
	case a == simulation.Noise:
		return change * c.NoiseVariance
	}
}
func (ts *temperatureSimulator) computeMinMax(t, c float64) float64 {
	switch t2 := t + c; {
	case t2 > ts.MaxTemp():
		return ts.MaxTemp()
	case t2 < ts.MinTemp():
		return ts.MinTemp()
	default:
		return t2
	}
}
