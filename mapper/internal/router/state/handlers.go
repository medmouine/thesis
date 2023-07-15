package state

import (
	"net/http"

	"github.com/medmouine/mapper/pkg/device"
	"github.com/medmouine/mapper/pkg/device/simulation"
	"github.com/samber/lo"
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
			log.Errorf("no simulator for device [%s]", d.ID())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if a := r.URL.Query().Get("anomaly"); a != "" {
			sim.IntroduceAnomaly(simulation.ParseAnomaly(a))
		}

		varsValues := getMultiQueryParams(r, []string{"gv", "sv", "dv", "nv"})
		simulation.HandleSimConfigUpdate(sim, varsValues)

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

func getMultiQueryParams(r *http.Request, k []string) map[string]string {
	return lo.Associate(k, func(s string) (string, string) {
		return s, getQueryParam(r, s)
	})
}

func getQueryParam(r *http.Request, k string) string {
	return r.URL.Query().Get(k)
}
