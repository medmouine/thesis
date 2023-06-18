package globals

import (
	"github.com/kubeedge/mappers-go/mappers/common"
	"medthesis/device-mapper/driver"
)

// SimDevice is the ble device configuration and client information.
type SimDevice struct {
	Instance  common.DeviceInstance
	Simulator *driver.Tempsim
}

var MqttClient *common.MqttClient
