package main

import (
	"strconv"
	"sync"
)

type RedisString struct {
	value string
	mutex *sync.RWMutex
}

func NewRedisString(value string) *RedisString {
	return &RedisString{
		value: value,
		mutex: &sync.RWMutex{},
	}
}

func (s *RedisString) Incr() (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	i, err := strconv.Atoi(s.value)
	if err != nil {
		return 0, err
	}

	i += 1
	s.value = strconv.Itoa(i)

	return i, nil
}
