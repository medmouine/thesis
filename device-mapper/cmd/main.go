package main

import (
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/service"
	"medthesis/device-mapper/driver"
)

// main Tempsim device program entry
func main() {
	d := &driver.Tempsim{}
	service.Bootstrap("Tempsim", d)
}
