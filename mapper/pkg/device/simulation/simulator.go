package simulation

import (
	"sync"
)

type Simulator[T interface{}] interface {
	IntroduceAnomaly(anomaly Anomaly)
	Anomaly() *Anomaly
	Config() *T
	SetConfig(T)
}

type BaseSimulator struct {
	Simulation Simulator[VarSimulationConfig]
	SimConfig  *VarSimulationConfig
	anomaly    *Anomaly
	mux        sync.Mutex
}

func (vs *BaseSimulator) Simulator() Simulator[VarSimulationConfig] {
	return vs.Simulation
}

func (vs *BaseSimulator) IntroduceAnomaly(a Anomaly) {
	vs.mux.Lock()
	vs.anomaly = &a
	vs.mux.Unlock()
}

func (vs *BaseSimulator) Config() *VarSimulationConfig {
	vs.mux.Lock()
	defer vs.mux.Unlock()
	return vs.SimConfig
}

func (vs *BaseSimulator) SetConfig(c VarSimulationConfig) {
	vs.mux.Lock()
	defer vs.mux.Unlock()
	vs.SimConfig = &c
}

func (vs *BaseSimulator) Anomaly() *Anomaly {
	vs.mux.Lock()
	defer vs.mux.Unlock()
	if vs.anomaly == nil {
		vs.anomaly = new(Anomaly)
	}
	return vs.anomaly
}
