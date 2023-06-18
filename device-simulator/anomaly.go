package device_simulator

type Anomaly int

const (
	None Anomaly = iota
	Spike
	Drift
	Noise
	Flatline
)
