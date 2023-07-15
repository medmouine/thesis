package temperature

import (
	"time"

	"github.com/medmouine/mapper/pkg/device"
	"github.com/medmouine/mapper/pkg/device/simulation"
)

type TemperatureData struct {
	*device.Timeseries
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

func NewData(t, h float64) *TemperatureData {
	return &TemperatureData{
		Timeseries: &device.Timeseries{
			Epoch: time.Now().Unix(),
		},
		Temperature: t,
		Humidity:    h,
	}
}

type StatePayload struct {
	DeviceID string                         `json:"device_id"`
	MaxTemp  float64                        `json:"max_temp"`
	MinTemp  float64                        `json:"min_temp"`
	Config   simulation.VarSimulationConfig `json:"simulation_config,omitempty" `
	Anomaly  string                         `json:"anomaly,omitempty"`
}
