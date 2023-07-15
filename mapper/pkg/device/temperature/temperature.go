package temperature

import (
	"math/rand"
	"time"

	"github.com/koddr/gosl"
	"github.com/medmouine/mapper/pkg/device"
)

type TemperatureDevice interface {
	device.Device[TemperatureData]
	MinTemp() float64
	MaxTemp() float64
}

func (ts *temperatureDevice) MinTemp() float64 {
	return ts.minTemp
}

func (ts *temperatureDevice) MaxTemp() float64 {
	return ts.maxTemp
}

func (ts *temperatureDevice) GetStatePayload() ([]byte, error) {
	s := &StatePayload{
		DeviceID: ts.ID(),
		MaxTemp:  ts.MaxTemp(),
		MinTemp:  ts.MinTemp(),
	}

	if sim := ts.Simulator(); sim != nil {
		s.Anomaly = sim.Anomaly().String()
		s.Config = *sim.Config()
	}

	return gosl.Marshal(s)
}

func (ts *temperatureDevice) GetDataPayload() ([]byte, error) {
	d := ts.Read()
	return gosl.Marshal(&d)
}

type temperatureDevice struct {
	*device.BaseDevice[TemperatureData]
	minTemp float64
	maxTemp float64
}

func newBaseTemperatureDevice(id string, pubInterval time.Duration, minTemp float64, maxTemp float64) *temperatureDevice {
	t := minTemp + rand.Float64()*(maxTemp-minTemp)
	h := 0.0 + rand.Float64()*(100.0-0.0)
	data := NewData(t, h)
	d := device.NewBaseDevice[TemperatureData](id, data, pubInterval)
	return &temperatureDevice{
		BaseDevice: d,
		minTemp:    minTemp,
		maxTemp:    maxTemp,
	}
}
