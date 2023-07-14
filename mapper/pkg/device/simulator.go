package device

import "sync"

type Simulator[T interface{}] interface {
	IntroduceAnomaly(anomaly Anomaly)
	Anomaly() *Anomaly
}

type BaseSimulator[T interface{}] struct {
	Device  Simulator[T]
	anomaly *Anomaly
	mux     sync.Mutex
}

func (s *BaseSimulator[T]) Simulator() Simulator[T] {
	return s.Device
}

func (s *BaseSimulator[T]) IntroduceAnomaly(a Anomaly) {
	s.mux.Lock()
	s.anomaly = &a
	s.mux.Unlock()
}

func (s *BaseSimulator[T]) Anomaly() *Anomaly {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.anomaly == nil {
		s.anomaly = new(Anomaly)
	}
	return s.anomaly
}
