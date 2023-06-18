package device_simulator

import (
	"math/rand"
	"time"
)

type TemperatureSensor struct {
	MinTemp     float64
	MaxTemp     float64
	Temperature float64
	Anomaly     Anomaly
	r           *rand.Rand
}

func NewTemperatureSensor(minTemp, maxTemp float64) *TemperatureSensor {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	initialTemp := minTemp + r.Float64()*(maxTemp-minTemp)
	return &TemperatureSensor{MinTemp: minTemp, MaxTemp: maxTemp, Temperature: initialTemp, Anomaly: None, r: r}
}

func (s *TemperatureSensor) Read() float64 {
	if s.Anomaly == Flatline {
		return s.Temperature
	}

	change := (s.r.Float64() - 0.5) * 3.0
	if s.Anomaly == Spike {
		change *= 10.0
	} else if s.Anomaly == Drift {
		change += 0.1
	} else if s.Anomaly == Noise {
		change *= 2.0
	}

	s.Temperature += change
	if s.Temperature < s.MinTemp {
		s.Temperature = s.MinTemp
	} else if s.Temperature > s.MaxTemp {
		s.Temperature = s.MaxTemp
	}

	return s.Temperature
}

func (s *TemperatureSensor) IntroduceAnomaly(anomaly Anomaly) {
	s.Anomaly = anomaly
}
