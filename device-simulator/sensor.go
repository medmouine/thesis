package device_simulator

type Sensor interface {
	Read() float64
	IntroduceAnomaly(anomaly Anomaly)
}
