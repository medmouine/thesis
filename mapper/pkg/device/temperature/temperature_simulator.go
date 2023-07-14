package temperature

import (
	"math/rand"
	"sync"

	"github.com/koddr/gosl"
	"github.com/medmouine/mapper/pkg/device"
)

type StatePayload struct {
	DeviceID string  `json:"device_id"`
	MaxTemp  float64 `json:"max_temp"`
	MinTemp  float64 `json:"min_temp"`
	Anomaly  string  `json:"anomaly,omitempty"`
}

type TemperatureSimulator interface {
	device.Simulator[Data]
	TemperatureDevice
}

func NewTemperatureSimulator(id string, minTemp, maxTemp float64) TemperatureSimulator {
	ts := &temperatureSimulator{
		temperatureDevice: newBaseTemperatureDevice(id, minTemp, maxTemp),
		BaseSimulator:     &device.BaseSimulator[Data]{},
	}
	ts.Device = ts
	return ts
}

type temperatureSimulator struct {
	*temperatureDevice
	*device.BaseSimulator[Data]
	mux sync.Mutex
}

func (ts *temperatureSimulator) GetStatePayload() ([]byte, error) {
	s := &StatePayload{
		DeviceID: ts.ID(),
		MaxTemp:  ts.MaxTemp(),
		MinTemp:  ts.MinTemp(),
	}

	if sim := ts.Simulator(); sim != nil {
		s.Anomaly = sim.Anomaly().String()
	}

	return gosl.Marshal(s)
}

func (ts *temperatureSimulator) GetDataPayload() ([]byte, error) {
	d := ts.Read()
	p := &Data{
		Temperature: d.Temperature,
		Humidity:    d.Humidity,
	}

	return gosl.Marshal(p)
}

func (ts *temperatureSimulator) Read() Data {
	ts.mux.Lock()
	defer ts.mux.Unlock()
	d := &Data{
		Temperature: ts.Data().Temperature,
		Humidity:    ts.Data().Humidity,
	}
	a := *ts.Anomaly()
	if a == device.Flatline {
		return ts.SetData(d)
	}

	change := (rand.Float64() - 0.5) * 3.0
	switch {
	case a == device.Spike:
		change *= 10.0
	case a == device.Drift:
		change += 0.1
	case a == device.Noise:
		change *= 2.0
	}

	d.Temperature += change
	d.Humidity += change
	if d.Temperature < ts.MinTemp() {
		d.Temperature = ts.minTemp
	} else if d.Temperature > ts.MaxTemp() {
		d.Temperature = ts.maxTemp
	}

	return ts.SetData(d)
}
