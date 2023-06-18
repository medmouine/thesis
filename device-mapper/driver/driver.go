package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

type TempsimProtocolConfig struct {
	ProtocolName       string `json:"protocolName"`
	ProtocolConfigData `json:"configData"`
}

type ProtocolConfigData struct {
	DeviceID int `json:"deviceID,omitempty"`
}

type TempsimProtocolCommonConfig struct {
	CommonCustomizedValues `json:"customizedValues"`
}

type CommonCustomizedValues struct {
	ProtocolID int `json:"protocolID"`
}
type TempsimVisitorConfig struct {
	ProtocolName      string `json:"protocolName"`
	VisitorConfigData `json:"configData"`
}

type VisitorConfigData struct {
	DataType string `json:"dataType"`
}

// Tempsim Realize the structure of random number
type Tempsim struct {
	mutex                 sync.Mutex
	virtualProtocolConfig TempsimProtocolConfig
	protocolCommonConfig  TempsimProtocolCommonConfig
	visitorConfig         TempsimVisitorConfig
	client                map[int]int64
}

// InitDevice Sth that need to do in the first
// If you need mount a persistent connection, you should provIDe parameters in configmap's protocolCommon.
// and handle these parameters in the following function
func (ts *Tempsim) InitDevice(protocolCommon []byte) (err error) {
	if protocolCommon != nil {
		if err = json.Unmarshal(protocolCommon, &ts.protocolCommonConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolCommonConfig error: %v\n", err)
			return err
		}
	}
	fmt.Printf("InitDevice%d...\n", ts.protocolCommonConfig.ProtocolID)
	return nil
}

// SetConfig Parse the configmap's raw json message
func (ts *Tempsim) SetConfig(protocolCommon, visitor, protocol []byte) (dataType string, deviceID int, err error) {
	// TODO parse simulator config
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	ts.NewClient()
	if protocolCommon != nil {
		if err = json.Unmarshal(protocolCommon, &ts.protocolCommonConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolCommonConfig error: %v\n", err)
			return "", 0, err
		}
	}
	if visitor != nil {
		if err = json.Unmarshal(visitor, &ts.visitorConfig); err != nil {
			fmt.Printf("Unmarshal visitorConfig error: %v\n", err)
			return "", 0, err
		}
	}

	if protocol != nil {
		if err = json.Unmarshal(protocol, &ts.virtualProtocolConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolConfig error: %v\n", err)
			return "", 0, err
		}
	}
	dataType = ts.visitorConfig.DataType
	deviceID = ts.virtualProtocolConfig.DeviceID
	return
}

// ReadDeviceData  is an interface that reads data from a specific device, data's dataType is consistent with configmap
func (ts *Tempsim) ReadDeviceData(protocolCommon, visitor, protocol []byte) (data interface{}, err error) {
	// TODO integrate simulator
	// Parse raw json message to get a Tempsim instance
	DataTye, DeviceID, err := ts.SetConfig(protocolCommon, visitor, protocol)
	if err != nil {
		return nil, err
	}
	if DataTye == "int" {
		if ts.client[DeviceID] == 0 {
			return 0, errors.New("ts.limit should not be 0")
		}
		return rand.Intn(int(ts.client[DeviceID])), nil
	} else if DataTye == "float" {
		if ts.client[DeviceID] == 0 {
			return 0, errors.New("ts.limit should not be 0")
		}
		// Simulate device that have time delay
		// time.Sleep(time.Second)
		return rand.Float64(), nil
	} else {
		return "", errors.New("dataType don't exist")
	}
}

// WriteDeviceData is an interface that write data to a specific device, data's dataType is consistent with configmap
func (ts *Tempsim) WriteDeviceData(data interface{}, protocolCommon, visitor, protocol []byte) (err error) {
	// Parse raw json message to get a Tempsim instance
	_, DeviceID, err := ts.SetConfig(protocolCommon, visitor, protocol)
	if err != nil {
		return err
	}
	ts.client[DeviceID] = data.(int64)
	return nil
}

// StopDevice is an interface to disconnect a specific device
// This function is called when mapper stops serving
func (ts *Tempsim) StopDevice() (err error) {
	// in this func, u can get ur device-instance in the client map, and give a safety exit
	fmt.Println("----------Stop Virtual Device Successful----------")
	return nil
}

// NewClient create device-instance, if device-instance exit, set the limit as 100.
// Control a group of devices through singleton pattern
func (ts *Tempsim) NewClient() {
	if ts.client == nil {
		ts.client = make(map[int]int64)
	}
	if _, ok := ts.client[ts.virtualProtocolConfig.DeviceID]; ok {
		if ts.client[ts.virtualProtocolConfig.DeviceID] == 0 {
			ts.client[ts.virtualProtocolConfig.DeviceID] = 100
		}
	}
}

// GetDeviceStatus is an interface to get the device status true is OK , false is DISCONNECTED
func (ts *Tempsim) GetDeviceStatus(protocolCommon, visitor, protocol []byte) (status bool) {
	_, _, err := ts.SetConfig(protocolCommon, visitor, protocol)
	return err == nil
}
