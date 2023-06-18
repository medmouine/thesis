package temperature_sensor

type Sensor interface {
	Read() float64
}

type temperatureSensor struct {
	Sensor
	Id          string
	MinTemp     float64
	MaxTemp     float64
	Temperature float64
}

func newTemperatureSensor(id string, minTemp, maxTemp, initialTemperature float64) *temperatureSensor {
	return &temperatureSensor{
		Id:          id,
		MinTemp:     minTemp,
		MaxTemp:     maxTemp,
		Temperature: initialTemperature,
	}
}

func (s *temperatureSensor) Read() float64 {
	panic("not implemented")
}
