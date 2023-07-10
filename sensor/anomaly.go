package temperaturesensor

import "fmt"

type Anomaly int

const (
	None Anomaly = iota
	Spike
	Drift
	Noise
	Flatline
)

func FromString(anomalyStr string) (Anomaly, error) {
	switch anomalyStr {
	case "Spike":
		return Spike, nil
	case "Drift":
		return Drift, nil
	case "Noise":
		return Noise, nil
	case "Flatline":
		return Flatline, nil
	case "None":
		return None, nil
	default:
		return None, fmt.Errorf("unknown anomaly: %s", anomalyStr)
	}
}
