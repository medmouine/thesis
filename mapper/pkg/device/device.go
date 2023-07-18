package device

import (
	"sync"
	"time"

	"github.com/medmouine/mapper/pkg/device/simulation"
)

type Device[T interface{}] interface {
	ID() string
	Location() string
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
	location        string
	data            *T
	publishInterval *time.Duration
	mux             sync.Mutex
}

func NewBaseDevice[T interface{}](id string, pubInterval time.Duration, location ...string) *BaseDevice[T] {
	var l = "unknown"
	if len(location) > 0 && location[0] != "" {
		l = location[0]
	}

	return &BaseDevice[T]{
		id:              id,
		location:        l,
		publishInterval: &pubInterval,
		mux:             sync.Mutex{},
	}
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
func (s *BaseDevice[T]) Location() string {
	return s.location
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
