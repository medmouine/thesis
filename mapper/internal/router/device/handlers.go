package device

import (
	"net/http"

	"github.com/koddr/gosl"
	temperaturesensor "github.com/medmouine/device-mapper/pkg/sensor"
	log "github.com/sirupsen/logrus"
)

type StateResponse struct {
	DeviceID string                    `json:"device_id"`
	Anomaly  temperaturesensor.Anomaly `json:"anomaly,omitempty"`
	MaxTemp  float64                   `json:"max_temp,omitempty"`
	MinTemp  float64                   `json:"min_temp,omitempty"`
}

func getState(driver *temperaturesensor.TemperatureSimulator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		state := &StateResponse{
			DeviceID: driver.ID,
			Anomaly:  driver.Anomaly,
			MaxTemp:  driver.MaxTemp,
			MinTemp:  driver.MinTemp,
		}

		if body, err := gosl.Marshal(state); err != nil {
			log.Errorf("Error marshaling device state: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(body)
			if err != nil {
				log.Errorf("Error writing response: %v", err)
			}
		}
	}
}

type DataResponse struct {
	Temperature float64 `json:"temperature"`
}

func getData(driver temperaturesensor.Sensor) func(w http.ResponseWriter, r *http.Request) {
	res := &DataResponse{
		Temperature: driver.Read(),
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		if body, err := gosl.Marshal(res); err != nil {
			log.Errorf("Error marshaling device data: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(body)
			if err != nil {
				log.Errorf("Error writing response: %v", err)
			}
		}
	}
}

func setAnomaly(driver *temperaturesensor.TemperatureSimulator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		a := r.URL.Query().Get("anomaly")
		log.Infof("Setting anomalyz to %v (current %s)", a, driver.Anomaly)
		driver.IntroduceAnomaly(temperaturesensor.ParseAnomaly(a))
		w.WriteHeader(http.StatusOK)
		state := &StateResponse{
			DeviceID: driver.ID,
			Anomaly:  driver.Anomaly,
		}

		if body, err := gosl.Marshal(state); err != nil {
			log.Errorf("Error marshaling response state: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(body)
			if err != nil {
				log.Errorf("Error writing response: %v", err)
			}
		}
	}
}
