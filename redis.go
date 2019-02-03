package main

import (
	"errors"
	"fmt"
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

	fmt.Println(v)

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

	fmt.Println(data)

	value, ok := data.(RedisString)
	fmt.Println(value)
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

func (r *Redis) Incr(key string) (string, error) {
	return "", nil
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
}

func NewRedisString(value string) RedisData {
	return &RedisString{
		value: value,
	}
}

type RedisSet struct {
}

func NewRedisSet() RedisData {
	return &RedisSet{}
}
