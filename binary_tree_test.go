package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBinaryTree(t *testing.T) {
	t.Run("with empty tree, try to get 0-th element, expect error", func(t *testing.T) {
		tree := NewBinaryTree()

		if _, err := tree.Card(0); err == nil {
			t.Errorf("error expected")
		}
	})

	t.Run("with empty tree, try to add 1 element, expect no error", func(t *testing.T) {
		tree := NewBinaryTree()
		val := NewMember(0, "zero")

		if err := tree.Add(val); err != nil {
			t.Errorf("error not expected: %v", err)
		}
	})

	t.Run("with empty tree, try to add and get 1 element, expect element", func(t *testing.T) {
		tree := NewBinaryTree()
		val := NewMember(0, "zero")

		if err := tree.Add(val); err != nil {
			t.Errorf("error not expected: %v", err)
		}

		got, err := tree.Card(0)
		if err != nil {
			t.Errorf("error not expected: %v", err)
		}

		// comparing pointers: should return same thing
		if val != got {
			t.Errorf("expected '%p', got '%p'", val, got)
		}
	})

	t.Run("with empty tree, try to add 15 and get 7-th element, expect 7-th element", func(t *testing.T) {
		tree := NewBinaryTree()
		var seventh *Member

		for x := 0; x < 15; x++ {
			val := NewMember(x, fmt.Sprint("member", x))
			if x == 7 {
				seventh = val
			}

			if err := tree.Add(val); err != nil {
				t.Errorf("error not expected: %v", err)
			}
		}

		got, err := tree.Card(7)
		if err != nil {
			t.Errorf("error not expected: %v", err)
		}

		if seventh != got {
			t.Errorf("expected '%p', got '%p'", seventh, got)
		}
	})

	t.Run("with empty tree, try to add 21 and get 9-th to 16-th element, expect elements", func(t *testing.T) {
		tree := NewBinaryTree()
		elems := []*Member{}

		start := 9
		stop := 16
		total := 21

		ranges := []struct {
			start int
			stop  int
		}{
			{9, 16},
			{9, -5},
			{-12, 16},
			{-12, -5},
		}

		for x := 0; x < total; x++ {
			val := NewMember(x, fmt.Sprint("member", x))
			if start <= x && x <= stop {
				elems = append(elems, val)
			}

			if err := tree.Add(val); err != nil {
				t.Errorf("error not expected: %v", err)
			}
		}

		for _, r := range ranges {
			got := tree.Range(r.start, r.stop)

			if len(elems) != len(got) {
				t.Errorf("expected '%p', got '%p'", elems, got)
			}

			for i := 0; i < len(elems); i++ {
				if elems[i] != got[i] {
					t.Errorf("expected '%p', got '%p'", elems[i], got[i])
				}
			}
		}
	})

	t.Run("with empty tree, try to add 1000 and some range, expect elements", func(t *testing.T) {
		tree := NewBinaryTree()
		elems := []*Member{}

		start := 130
		stop := 650
		total := 1000
		th := 300
		var elem *Member
		var next *Member

		ranges := []struct {
			start int
			stop  int
		}{
			{start, stop},
			{start, stop - total},
			{start - total, stop},
			{start - total, stop - total},
		}

		// add elements
		for x := 0; x < total; x++ {
			val := NewMember(x, fmt.Sprint("member", x))
			if start <= x && x <= stop {
				elems = append(elems, val)
			}

			if th == x {
				elem = val
			}

			if th+1 == x {
				next = val
			}

			if err := tree.Add(val); err != nil {
				t.Errorf("error not expected: %v", err)
			}
		}

		rank, err := tree.Rank(elem)
		if err != nil {
			t.Errorf("error not expected: %v", err)
		}

		if rank != th {
			t.Errorf("expected '%d', got '%d'", th, rank)
		}

		val, err := tree.Card(th)
		if err != nil {
			t.Errorf("error not expected: %v", err)
		}

		if val != elem {
			t.Errorf("expected '%p', got '%p'", elem, val)
		}

		// check
		for _, r := range ranges {
			got := tree.Range(r.start, r.stop)

			if len(elems) != len(got) {
				t.Errorf("expected '%p', got '%p'", elems, got)
			}

			for i := 0; i < len(elems); i++ {
				if elems[i] != got[i] {
					t.Errorf("expected '%p', got '%p'", elems[i], got[i])
				}
			}
		}

		tree.Del(elem)

		rank, err = tree.Rank(elem)
		if err == nil {
			t.Errorf("error expected")
		}

		val, err = tree.Card(th)
		if err != nil {
			t.Errorf("error not expected: %v", err)
		}

		if val != next {
			t.Errorf("expected '%p', got '%p'", next, val)
		}
	})
}
