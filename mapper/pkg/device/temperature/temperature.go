package temperature

import (
	"math/rand"

	"github.com/medmouine/mapper/pkg/device"
)

type Data struct {
	Temperature float64
	Humidity    float64
}

type TemperatureDevice interface {
	device.Device[Data]
	MinTemp() float64
	MaxTemp() float64
}

type temperatureDevice struct {
	*device.BaseDevice[Data]
	minTemp float64
	maxTemp float64
}

func (ts *temperatureDevice) MinTemp() float64 {
	return ts.minTemp
}

func (ts *temperatureDevice) MaxTemp() float64 {
	return ts.maxTemp
}

func newBaseTemperatureDevice(id string, minTemp float64, maxTemp float64) *temperatureDevice {
	data := &Data{
		Temperature: minTemp + rand.Float64()*(maxTemp-minTemp),
		Humidity:    0,
	}
	return &temperatureDevice{
		minTemp:    minTemp,
		maxTemp:    maxTemp,
		BaseDevice: device.NewBaseDevice(id, data),
	}
}
