package main

import (
	"errors"
	"sync"
)

type SafeMap struct {
	mutex *sync.RWMutex
	store map[string]interface{}
}

func NewSafeMap() (*SafeMap, error) {
	return &SafeMap{
		store: map[string]interface{}{},
		mutex: &sync.RWMutex{},
	}, nil
}

func (m *SafeMap) Add(key string, value interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.store[key] = value

	return nil
}

func (m *SafeMap) Get(key string) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, ok := m.store[key]
	if !ok {
		return nil, errors.New("Not found")
	}

	return value, nil
}

func (m *SafeMap) Del(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.store[key]; ok {
		delete(m.store, key)
		return nil
	}

	return errors.New("Not found")
}

func (m *SafeMap) DelIf(key string, value interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if stored, ok := m.store[key]; ok {
		// deletes only if stored value and param value are the same
		if stored == value {
			delete(m.store, key)
			return nil
		}
		return errors.New("Not same value")
	}

	return errors.New("Not found")
}

func (m *SafeMap) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.store)
}
