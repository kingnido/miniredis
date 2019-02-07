package main

import (
	"errors"
	"time"
)

var (
	InvalidExpirationError = errors.New("expiration time should be greater than 1s")
	NotAStringError        = errors.New("value is not a string")
	NotASetError           = errors.New("value is not a set")
)

type Redis struct {
	store *SafeMap
}

func NewRedis() *Redis {
	return &Redis{
		store: NewSafeMap(),
	}
}

func (r *Redis) Set(key string, value string) {
	r.store.Set(key, NewRedisString(value))
}

func (r *Redis) SetExpire(key string, value string, expires int) error {
	if expires <= 0 {
		return InvalidExpirationError
	}

	k := key
	v := NewRedisString(value)

	// add the key-pair value and run timeout task
	r.store.Set(k, v)
	go func(timeout <-chan time.Time) {
		<-timeout
		r.store.DelIf(k, v)
	}(time.NewTimer(time.Duration(expires) * time.Second).C)

	return nil
}

func (r *Redis) Get(key string) (string, error) {
	data, err := r.store.Get(key)
	if err != nil {
		return "", err
	}

	value, ok := data.(*RedisString)
	if !ok {
		return "", NotAStringError
	}

	return value.Value(), nil
}

func (r *Redis) Del(key string) error {
	return r.store.Del(key)
}

func (r *Redis) DbSize() int {
	return r.store.Size()
}

func (r *Redis) Incr(key string) (int, error) {
	data := r.store.GetOrDefault(key, NewRedisString("0"))

	value, ok := data.(*RedisString)
	if !ok {
		return 0, NotAStringError
	}

	return value.Incr()
}

func (r *Redis) ZAdd(key string, score int, member string) (int, error) {
	data := r.store.GetOrDefault(key, NewRedisSet())

	set, ok := data.(*RedisSet)
	if !ok {
		return 0, NotASetError
	}

	set.Set(score, member)

	return 1, nil
}

func (r *Redis) ZCard(key string) (int, error) {
	var set *RedisSet

	data, err := r.store.Get(key)
	if err != nil {
		return -1, KeyNotFoundError
	}

	set, ok := data.(*RedisSet)
	if !ok {
		return 0, NotASetError
	}

	return set.Card(), nil
}

func (r *Redis) ZRank(key string, member string) (int, error) {
	var set *RedisSet

	data, err := r.store.Get(key)
	if err != nil {
		return -1, err
	}

	set, ok := data.(*RedisSet)
	if !ok {
		return 0, NotASetError
	}

	return set.Rank(member)
}

func (r *Redis) ZRange(key string, start int, stop int) ([]string, error) {
	var set *RedisSet

	data, err := r.store.Get(key)
	if err != nil {
		return nil, KeyNotFoundError
	}

	set, ok := data.(*RedisSet)
	if !ok {
		return nil, NotASetError
	}

	return set.Range(start, stop), nil
}
