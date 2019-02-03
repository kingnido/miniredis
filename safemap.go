package main

import (
	"errors"
	"sync"
)

type Key string
type Value interface{}

type SafeMap struct {
	store map[Key]Value
	mutex sync.Mutex
}

func NewSafeMap() (*SafeMap, error) {
	return &SafeMap{
		store: map[Key]Value{},
		mutex: sync.Mutex{},
	}, nil
}

func (m *SafeMap) Add(key Key, value Value) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.store[key] = value

	return nil
}

func (m *SafeMap) Get(key Key) (Value, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	value, ok := m.store[key]
	if !ok {
		return nil, errors.New("Not found")
	}

	return value, nil
}

func (m *SafeMap) Del(key Key) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.store[key]; ok {
		delete(m.store, key)
		return nil
	}

	return errors.New("Not found")
}

func (m *SafeMap) DelIf(key Key, value Value) error {
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
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return len(m.store)
}
