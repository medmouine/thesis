package device

import "time"

type BaseData struct {
	Epoch    int64  `json:"timestamp"`
	DeviceID string `json:"device_id"`
	Location string `json:"location,omitempty"`
}

func NewBaseData(DeviceID, Location string) *BaseData {
	return &BaseData{
		Epoch:    time.Now().Unix(),
		DeviceID: DeviceID,
		Location: Location,
	}
}
