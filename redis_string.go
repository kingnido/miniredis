package main

import (
	"strconv"
	"sync"
)

type RedisString struct {
	sync.RWMutex
	value string
}

func NewRedisString(value string) *RedisString {
	return &RedisString{
		value: value,
	}
}

func (s *RedisString) Value() string {
	s.RLock()
	defer s.RUnlock()

	return s.value
}

func (s *RedisString) Incr() (int, error) {
	s.Lock()
	defer s.Unlock()

	i, err := strconv.Atoi(s.value)
	if err != nil {
		return 0, err
	}

	i += 1
	s.value = strconv.Itoa(i)

	return i, nil
}
