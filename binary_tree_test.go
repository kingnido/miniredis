package main

import (
	"testing"
)

func TestBinaryTree(t *testing.T) {
	t.Run("test plz", func(t *testing.T) {
		tree := NewBinaryTree()

		for i := 0; i < 15; i++ {
			tree.Add(i)
		}

		tree.Details()

		if tree == nil {
			t.Errorf("tree should not be nil")
		}
	})
}
