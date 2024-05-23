package fxstate

import (
	"fmt"
	"sync"
)

type State interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

type state struct {
	sync.Mutex

	data map[string]interface{}
}

func New(initial map[string]interface{}) State {
	return &state{data: initial}
}

func (s *state) Set(key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()

	s.data[key] = value

	return nil
}

func (s *state) Get(key string) (interface{}, error) {
	s.Lock()
	defer s.Unlock()

	value, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("value not found by key = %v", key)
	}

	return value, nil
}
