package main

import (
	"errors"
	"sync"
)

var (
	MemberNotFound = errors.New("member not found")
)

type RedisSet struct {
	sync.RWMutex

	store map[string]*Member
	tree  *BinaryTree
}

func NewRedisSet() *RedisSet {
	return &RedisSet{
		store: map[string]*Member{},
		tree:  NewBinaryTree(),
	}
}

func (s *RedisSet) Set(score int, member string) error {
	s.Lock()
	defer s.Unlock()

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
	s.RLock()
	defer s.RUnlock()

	return len(s.store)
}

func (s *RedisSet) Rank(member string) (int, error) {
	s.RLock()
	defer s.RUnlock()

	m, ok := s.store[member]
	if !ok {
		return -1, MemberNotFound
	}

	return s.tree.Rank(m)
}

func (s *RedisSet) Range(start int, stop int) []string {
	s.RLock()
	list := s.tree.Range(start, stop)
	s.RUnlock()

	r := make([]string, 0, len(list))
	for _, v := range list {
		s, _ := v.(*Member)
		r = append(r, s.Member)
	}

	return r
}
