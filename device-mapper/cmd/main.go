package main

import (
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/service"
	"medthesis/temperature_sensor"
)

// main TemperatureSensorDriver device program entry
func main() {
	d := temperature_sensor.NewTemperatureSimulator("temperature-sensor-1", 0, 100)
	service.Bootstrap("Simulator", d)
}
