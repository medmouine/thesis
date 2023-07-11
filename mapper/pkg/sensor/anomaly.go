package temperaturesensor

type Anomaly string

const (
	None     Anomaly = "None"
	Spike    Anomaly = "Spike"
	Drift    Anomaly = "Drift"
	Noise    Anomaly = "Noise"
	Flatline Anomaly = "Flatline"
)

func ParseAnomaly(i interface{}) Anomaly {
	switch a := i.(type) {
	case string:
		return parseString(a)
	case Anomaly:
		return a
	default:
		return None
	}
}

func parseString(s string) Anomaly {
	switch s {
	case "None":
		return None
	case "Spike":
		return Spike
	case "Drift":
		return Drift
	case "Noise":
		return Noise
	case "Flatline":
		return Flatline
	default:
		return None
	}
}
