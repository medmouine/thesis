package simulation

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

type VarSimulationConfig struct {
	GlobalVariance float64 `json:"variance"`
	SpikeVariance  float64 `json:"spike_variance"`
	DriftVariance  float64 `json:"drift_variance"`
	NoiseVariance  float64 `json:"noise_variance"`
}

func DefaultVarSimConfig() *VarSimulationConfig {
	return &VarSimulationConfig{
		GlobalVariance: 3.0,
		SpikeVariance:  10.0,
		DriftVariance:  .1,
		NoiseVariance:  2.0,
	}
}

func HandleSimConfigUpdate(sim Simulator[VarSimulationConfig], params map[string]string) {
	if len(params) == 0 {
		return
	}
	current := sim.Config()
	sim.SetConfig(VarSimulationConfig{
		GlobalVariance: parseVarsValues(params["gv"], current.GlobalVariance),
		SpikeVariance:  parseVarsValues(params["sv"], current.SpikeVariance),
		DriftVariance:  parseVarsValues(params["dv"], current.DriftVariance),
		NoiseVariance:  parseVarsValues(params["nv"], current.NoiseVariance),
	})
}

func parseVarsValues(v string, defaultV float64) float64 {
	if v == "" {
		return defaultV
	}
	fv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Errorf("Error parsing float value %s: %v", v, err)
		return defaultV
	}
	return fv
}
