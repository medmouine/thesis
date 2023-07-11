package device

import (
	"github.com/go-chi/chi/v5"
	temperaturesensor "github.com/medmouine/device-mapper/pkg/sensor"
)

const (
	groupURL   = "/device"
	stateURL   = "/state"
	dataURL    = "/data"
	anomalyURL = "/anomaly"
)

func Routes(m *chi.Mux, driver *temperaturesensor.TemperatureSimulator) {
	m.Route(groupURL, func(r chi.Router) {
		r.Get(stateURL, getState(driver))
		r.Get(dataURL, getData(driver))
		r.Put(anomalyURL, setAnomaly(driver))
	})
}
