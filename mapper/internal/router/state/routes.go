package state

import (
	"github.com/go-chi/chi/v5"
	"github.com/medmouine/mapper/pkg/device"
)

const (
	groupURL     = "/temperature"
	stateURL     = "/state"
	dataURL      = "/data"
	simulatorURL = "/simulator"
)

func Routes[T interface{}](m *chi.Mux, d device.Device[T]) {
	m.Route(groupURL, func(r chi.Router) {
		r.Get(stateURL, getState(d))
		r.Get(dataURL, getData(d))
		r.Put(simulatorURL, putConfig(d))
	})
}
