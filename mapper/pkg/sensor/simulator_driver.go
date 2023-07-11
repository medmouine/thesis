package temperaturesensor

import (
	"fmt"
)

func (ts *TemperatureSimulator) InitDevice(_ []byte) (err error) {
	return nil
}

func (ts *TemperatureSimulator) ReadDeviceData(_, _, _ []byte) (data interface{}, err error) {
	return ts.Read(), nil
}

func (ts *TemperatureSimulator) WriteDeviceData(data interface{}, _, _, _ []byte) (err error) {
	anomaly := ParseAnomaly(data)
	fmt.Printf("introducing anomaly: %s\n", anomaly)
	ts.IntroduceAnomaly(anomaly)
	return nil
}

func (ts *TemperatureSimulator) StopDevice() (err error) {
	fmt.Println("----------Stop Virtual Device Successful----------")
	return nil
}

func (ts *TemperatureSimulator) GetDeviceStatus(_, _, _ []byte) (status bool) {
	return true
}
