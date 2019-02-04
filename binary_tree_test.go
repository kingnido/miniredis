package main

import (
	"math/rand"
	"testing"
)

func TestBinaryTree(t *testing.T) {
	t.Run("test plz", func(t *testing.T) {
		tree := NewBinaryTree()

		for i := 0; i < 21; i++ {
			tree.Add(i)
		}

		tree.Details()

		if tree == nil {
			t.Errorf("tree should not be nil")
		}
	})

	t.Run("test plz", func(t *testing.T) {
		tree := NewBinaryTree()

		for _, x := range rand.Perm(21) {
			tree.Add(x)
		}

		tree.Details()

		if tree == nil {
			t.Errorf("tree should not be nil")
		}
	})

	t.Run("test plz del", func(t *testing.T) {
		tree := NewBinaryTree()

		for x := 0; x < 21; x++ {
			tree.Add(x)
		}

		for x := 5; x < 11; x++ {
			tree.Del(x)
		}

		tree.Details()

		if tree == nil {
			t.Errorf("tree should not be nil")
		}
	})
}
