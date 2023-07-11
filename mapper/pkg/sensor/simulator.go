package temperaturesensor

import (
	"math/rand"
	"sync"
	"time"
)

type Simulator interface {
	IntroduceAnomaly(anomaly Anomaly)
}

type TemperatureSimulator struct {
	Simulator
	*temperatureSensor
	Anomaly Anomaly
	r       *rand.Rand
	mux     sync.Mutex
}

func NewTemperatureSimulator(id string, minTemp, maxTemp float64) *TemperatureSimulator {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initialTemp := minTemp + r.Float64()*(maxTemp-minTemp)
	simulator := &TemperatureSimulator{
		temperatureSensor: newTemperatureSensor(id, minTemp, maxTemp, initialTemp),
		Anomaly:           None,
		r:                 r,
		mux:               sync.Mutex{},
	}
	simulator.temperatureSensor.Sensor = simulator
	return simulator
}

func (ts *TemperatureSimulator) Read() float64 {
	if ts.Anomaly == Flatline {
		return ts.Temperature
	}

	change := (ts.r.Float64() - 0.5) * 3.0
	switch {
	case ts.Anomaly == Spike:
		change *= 10.0
	case ts.Anomaly == Drift:
		change += 0.1
	case ts.Anomaly == Noise:
		change *= 2.0
	}

	ts.Temperature += change
	if ts.Temperature < ts.MinTemp {
		ts.Temperature = ts.MinTemp
	} else if ts.Temperature > ts.MaxTemp {
		ts.Temperature = ts.MaxTemp
	}

	return ts.Temperature
}

func (ts *TemperatureSimulator) IntroduceAnomaly(anomaly Anomaly) {
	ts.mux.Lock()
	ts.Anomaly = anomaly
	ts.mux.Unlock()
}
