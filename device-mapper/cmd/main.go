package main

import (
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/service"
	"github.com/medmouine/temperature-sensor"
)

func main() {
	d := temperaturesensor.NewTemperatureSimulator("temperature-sensor-1", 0, 100)
	service.Bootstrap("Simulator", d)
}
