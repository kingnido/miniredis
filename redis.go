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

type RedisData interface {
}

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

type RedisSet struct {
	store map[string]*Member
	tree  *BinaryTree
	mutex *sync.RWMutex
}

func NewRedisSet() *RedisSet {
	return &RedisSet{
		store: map[string]*Member{},
		tree:  NewBinaryTree(),
		mutex: &sync.RWMutex{},
	}
}

func (s *RedisSet) Add(score int, member string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var m *Member
	var ok bool

	if m, ok = s.store[member]; ok {
		// member in the set. update
		s.tree.Del(m)
	}

	m = NewMember(score, member)
	s.store[member] = m
	s.tree.Add(m)

	return nil
}

func (s *RedisSet) Card() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.store)
}

func (s *RedisSet) Rank(member string) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	m, ok := s.store[member]
	if !ok {
		return -1, errors.New("not found")
	}

	return s.tree.Rank(m)
}

func (s *RedisSet) Range(start int, stop int) []string {
	s.mutex.RLock()
	list := s.tree.Range(start, stop)
	s.mutex.RUnlock()

	r := make([]string, 0, len(list))
	for _, v := range list {
		s, _ := v.(*Member)
		r = append(r, s.Member)
	}

	return r
}
