package main

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

type Redis struct {
	store *SafeMap
}

func NewRedis() (*Redis, error) {
	store, _ := NewSafeMap()
	return &Redis{
		store: store,
	}, nil
}

func (r *Redis) Set(key string, value string) error {
	k := Key(key)
	v := NewRedisString(value)

	return r.store.Add(k, v)
}

func (r *Redis) SetExpire(key string, value string, expires int) error {
	if expires <= 0 {
		return errors.New("Expiration time should be greater than 1s")
	}

	k := Key(key)
	v := NewRedisString(value)

	go func(timeout <-chan time.Time) {
		<-timeout
		r.store.DelIf(k, v)
	}(time.NewTimer(time.Duration(expires) * time.Second).C)

	return r.store.Add(k, v)
}

func (r *Redis) Get(key string) (string, error) {
	k := Key(key)

	data, err := r.store.Get(k)
	if err != nil {
		return "", err
	}

	value, ok := data.(*RedisString)
	if !ok {
		return "", errors.New("Invalid type")
	}

	return value.value, nil
}

func (r *Redis) Del(key string) error {
	k := Key(key)

	return r.store.Del(k)
}

func (r *Redis) DbSize() int {
	return r.store.Size()
}

func (r *Redis) Incr(key string) (int, error) {
	k := Key(key)

	data, err := r.store.Get(k)
	if err != nil {
		data = NewRedisString("0")
		r.store.Add(k, data)
	}

	value, ok := data.(*RedisString)
	if !ok {
		return 0, errors.New("Invalid type")
	}

	return value.Incr()
}

func (r *Redis) ZAdd(key string, score string, member string) (int, error) {
	return 0, nil
}

func (r *Redis) ZCard(key string) (int, error) {
	return 0, nil
}

func (r *Redis) ZRank(key string, member string) (int, error) {
	return 0, nil
}

func (r *Redis) ZRange(key string, start string, stop string) ([]string, error) {
	return nil, nil
}

type RedisData interface {
}

type RedisString struct {
	value string
	mutex sync.Mutex
}

func NewRedisString(value string) RedisData {
	return &RedisString{
		value: value,
		mutex: sync.Mutex{},
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

type RedisSet struct {
}

func NewRedisSet() RedisData {
	return &RedisSet{}
}
