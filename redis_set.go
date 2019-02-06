package main

import (
	"errors"
	"sync"
)

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
