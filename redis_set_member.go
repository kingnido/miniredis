package main

import (
	"fmt"
	"reflect"
)

type Member struct {
	Score  int
	Member string
}

func NewMember(score int, member string) *Member {
	return &Member{
		Score:  score,
		Member: member,
	}
}

func (a *Member) LessThan(o BinaryTreeValue) bool {
	b, ok := o.(*Member)
	if !ok {
		panic(fmt.Sprintf("Type mismatch. Expected *Member, got %s\n",
			reflect.TypeOf(o)))
	}

	if a.Score == b.Score {
		return a.Member < b.Member
	}

	return a.Score < b.Score
}

func (a *Member) EqualTo(o BinaryTreeValue) bool {
	b, ok := o.(*Member)
	if !ok {
		panic(fmt.Sprintf("Type mismatch. Expected *Member, got %s\n",
			reflect.TypeOf(o)))
	}

	return a.Score == b.Score && a.Member == b.Member
}
