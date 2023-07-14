package state

import (
	"net/http"

	"github.com/medmouine/mapper/pkg/device"
	log "github.com/sirupsen/logrus"
)

func getState[T interface{}](d device.Device[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		if body, err := d.GetStatePayload(); err != nil {
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
	Humidity    float64 `json:"humidity"`
}

func getData[T interface{}](d device.Device[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		if body, err := d.GetDataPayload(); err != nil {
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

func putConfig[T interface{}](d device.Device[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sim := d.Simulator()
		if sim == nil {
			log.Errorf("Cannot set anomaly on device %s: no simulator", d.ID())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		a := r.URL.Query().Get("anomaly")
		log.Infof("Setting anomaly to %v (current %s)", a, sim.Anomaly())
		sim.IntroduceAnomaly(device.ParseAnomaly(a))
		w.WriteHeader(http.StatusOK)

		if body, err := d.GetStatePayload(); err != nil {
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
