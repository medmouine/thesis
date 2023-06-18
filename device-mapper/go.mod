module medthesis/device-mapper

go 1.20

require (
	github.com/kubeedge/mappers-go v1.13.0
	medthesis/temperature-sensor v0.0.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.3.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/spf13/pflag v1.0.6-0.20210604193023-d5e0c0615ace // indirect
	golang.org/x/net v0.5.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
)

replace medthesis/temperature-sensor => ../temperature-sensor

replace golang.org/x/net v0.0.0-20220225172249-27dd8689420f => golang.org/x/net v0.5.0 // indirect
