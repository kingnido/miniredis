package main

import (
	"errors"
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

func (r *Redis) ZAdd(key string, score int, member string) (int, error) {
	k := Key(key)

	var set *RedisSet

	data, err := r.store.Get(k)
	if err != nil {
		// empty. create new set
		set = NewRedisSet()
		set.Add(score, member)
		r.store.Add(k, set)
	} else {
		// has something. try to cast
		set, ok := data.(*RedisSet)
		if !ok {
			return 0, errors.New("Invalid type")
		}
		// add member to set. it is in the store already
		set.Add(score, member)
	}

	return 1, nil
}

func (r *Redis) ZCard(key string) (int, error) {
	k := Key(key)
	var set *RedisSet

	data, err := r.store.Get(k)
	if err != nil {
		return -1, errors.New("not found")
	}

	set, ok := data.(*RedisSet)
	if !ok {
		return 0, errors.New("Invalid type")
	}

	return set.Card(), nil
}

func (r *Redis) ZRank(key string, member string) (int, error) {
	k := Key(key)
	var set *RedisSet

	data, err := r.store.Get(k)
	if err != nil {
		return -1, errors.New("not found")
	}

	set, ok := data.(*RedisSet)
	if !ok {
		return 0, errors.New("Invalid type")
	}

	return set.Rank(member)
}

func (r *Redis) ZRange(key string, start int, stop int) ([]string, error) {
	k := Key(key)
	var set *RedisSet

	data, err := r.store.Get(k)
	if err != nil {
		return nil, errors.New("not found")
	}

	set, ok := data.(*RedisSet)
	if !ok {
		return nil, errors.New("Invalid type")
	}

	return set.Range(start, stop), nil
}
