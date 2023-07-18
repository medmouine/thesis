package temperature

import (
	"math/rand"
	"os"
	"time"

	"github.com/medmouine/mapper/pkg/device/simulation"
)

type TemperatureSimulator interface {
	simulation.Simulator[simulation.VarSimulationConfig]
	TemperatureDevice
}

func NewTemperatureSimulator(id string, pi time.Duration, minTemp, maxTemp float64) TemperatureSimulator {
	location := os.Getenv("DEVICE_LOCATION")
	ts := &temperatureSimulator{
		temperatureDevice: newBaseTemperatureDevice(id, pi, minTemp, maxTemp, location),
		BaseSimulator: &simulation.BaseSimulator{
			SimConfig: simulation.DefaultVarSimConfig(),
		},
	}
	ts.Simulation = ts
	ts.Device = ts

	ts.Read()
	return ts
}

type temperatureSimulator struct {
	*temperatureDevice
	*simulation.BaseSimulator
}

func (ts *temperatureSimulator) Read() TemperatureData {
	a := *ts.Anomaly()
	current := ts.Data()
	if current == nil {
		return ts.SetData(NewData(0, 0, ts))
	}
	t := ts.computeMinMax(current.Temperature, ts.computeVar(a))
	h := ts.computeMinMax(current.Humidity, ts.computeVar(a))
	return ts.SetData(NewData(t, h, ts))
}

func (ts *temperatureSimulator) computeVar(a simulation.Anomaly) float64 {
	c := ts.Config()
	change := (rand.Float64() - 0.5) * c.GlobalVariance
	switch a {
	default:
		return change
	case simulation.Flatline:
		return 0
	case simulation.Spike:
		return change * c.SpikeVariance
	case simulation.Drift:
		return change + c.DriftVariance
	case simulation.Noise:
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
