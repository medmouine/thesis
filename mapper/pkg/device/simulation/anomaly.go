package simulation

import "strings"

type Anomaly int

const (
	None Anomaly = iota
	Spike
	Drift
	Noise
	Flatline
)

func (a Anomaly) String() string {
	switch a {
	case Spike:
		return "Spike"
	case Drift:
		return "Drift"
	case Noise:
		return "Noise"
	case Flatline:
		return "Flatline"
	case None:
	default:
	}
	return "None"
}

func ParseAnomaly(a interface{}) Anomaly {
	switch t := a.(type) {
	case Anomaly:
		return t
	case string:
		return parseString(t)
	default:
		return None
	}
}

func parseString(s string) Anomaly {
	switch strings.ToLower(s) {
	case "none":
		return None
	case "spike":
		return Spike
	case "drift":
		return Drift
	case "noise":
		return Noise
	case "flatline":
		return Flatline
	default:
		return None
	}
}
