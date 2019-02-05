package main

import (
	"fmt"
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

		a, b := tree.Rank(11)
		fmt.Println(a, b)

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

	t.Run("test plz rank", func(t *testing.T) {
		tree := NewBinaryTree()

		for x := 0; x < 21; x++ {
			tree.Add(x)
		}

		tree.Details()

		for x := 0; x < 21; x++ {
			r, err := tree.Rank(x)
			fmt.Printf("value: %d, rank: %d, error: %v\n", x, r, err)
		}

		if tree == nil {
			t.Errorf("tree should not be nil")
		}
	})

	t.Run("test plz range", func(t *testing.T) {
		tree := NewBinaryTree()

		for x := 0; x < 21; x++ {
			tree.Add(x)
		}

		tree.Details()

		r := tree.Range(4, 17)

		fmt.Println(r)

		if tree == nil {
			t.Errorf("tree should not be nil")
		}
	})

	t.Run("test plz card", func(t *testing.T) {
		tree := NewBinaryTree()

		for x := 0; x < 21; x++ {
			tree.Add(x)
		}

		tree.Details()

		r, err := tree.Card(8)

		fmt.Println(r, err)

		if tree == nil {
			t.Errorf("tree should not be nil")
		}
	})

}
