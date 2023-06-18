module medthesis/device-mapper

go 1.17

require (
	github.com/kubeedge/mappers-go v1.13.0
	medthesis/temperature_sensor v0.0.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.3.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/spf13/pflag v1.0.6-0.20210604193023-d5e0c0615ace // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
)

replace medthesis/temperature_sensor => ../temperature-sensor