package device

import (
	"sync"
	"time"

	"github.com/medmouine/mapper/pkg/device/simulation"
)

type Device[T interface{}] interface {
	ID() string
	Data() *T
	SetData(data *T) T
	Read() T
	GetStatePayload() ([]byte, error)
	GetDataPayload() ([]byte, error)
	Simulator() simulation.Simulator[simulation.VarSimulationConfig]
	PublishInterval() time.Duration
	SetPublishInterval(d time.Duration)
}

type BaseDevice[T interface{}] struct {
	Device[T]
	id              string
	data            *T
	publishInterval *time.Duration
	mux             sync.Mutex
}

func (s *BaseDevice[T]) SetPublishInterval(d time.Duration) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.publishInterval = &d
}

func (s *BaseDevice[T]) PublishInterval() time.Duration {
	s.mux.Lock()
	defer s.mux.Unlock()
	return *s.publishInterval
}

func NewBaseDevice[T interface{}](id string, data *T, pubInterval time.Duration) *BaseDevice[T] {
	return &BaseDevice[T]{
		id:              id,
		data:            data,
		publishInterval: &pubInterval,
	}
}

func (s *BaseDevice[T]) ID() string {
	return s.id
}

func (s *BaseDevice[T]) SetData(data *T) T {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.data = data
	return *s.data
}

func (s *BaseDevice[T]) Data() *T {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.data
}
