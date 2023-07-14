package device

import "sync"

type Device[T interface{}] interface {
	ID() string
	Data() *T
	SetData(data *T) T
	Read() T
	GetStatePayload() ([]byte, error)
	GetDataPayload() ([]byte, error)
	Simulator() Simulator[T]
}

type BaseDevice[T interface{}] struct {
	Device[T]
	id   string
	data *T
	mux  sync.Mutex
}

func NewBaseDevice[T interface{}](id string, data *T) *BaseDevice[T] {
	return &BaseDevice[T]{
		id:   id,
		data: data,
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
