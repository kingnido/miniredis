package main

import (
	"errors"
	"sync"
)

var (
	KeyNotFoundError      = errors.New("key not found")
	KeyValueNotFoundError = errors.New("key-value pairs not found")
)

type SafeMap struct {
	sync.RWMutex
	store map[string]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		store: map[string]interface{}{},
	}
}

func (m *SafeMap) Set(key string, value interface{}) {
	m.Lock()
	defer m.Unlock()

	m.store[key] = value
}

func (m *SafeMap) Get(key string) (interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	value, ok := m.store[key]
	if !ok {
		return nil, KeyNotFoundError
	}

	return value, nil
}

func (m *SafeMap) GetOrDefault(key string, value interface{}) interface{} {
	m.Lock()
	defer m.Unlock()

	v, ok := m.store[key]
	if ok {
		return v
	}
	m.store[key] = value

	return value
}

func (m *SafeMap) Del(key string) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.store[key]; ok {
		delete(m.store, key)
		return nil
	}

	return KeyNotFoundError
}

func (m *SafeMap) DelIf(key string, value interface{}) error {
	m.Lock()
	defer m.Unlock()

	if stored, ok := m.store[key]; ok {
		if stored == value {
			// deletes only if stored value and param value are the same
			delete(m.store, key)
			return nil
		}

		return KeyValueNotFoundError
	}

	return KeyNotFoundError
}

func (m *SafeMap) Size() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.store)
}
