package main

import (
	"errors"
	"fmt"
	"strings"
)

type BinaryTree struct {
	root *node
}

type BinaryTreeValue interface {
	LessThan(BinaryTreeValue) bool
	EqualTo(BinaryTreeValue) bool
}

func NewBinaryTree() *BinaryTree {
	return &BinaryTree{}
}

func (t *BinaryTree) Add(value int) error {
	t.root = insert(t.root, value)

	return nil
}

func (t *BinaryTree) Del(value int) error {
	t.root = remove(t.root, value)

	return nil
}

func (t *BinaryTree) Rank(value int) (int, error) {
	return rank(t.root, value)
}

func (t *BinaryTree) Range(start int, stop int) []int {
	if t.root == nil {
		return []int{}
	}

	if start < 0 {
		start += weight(t.root)
	}

	if stop < 0 {
		stop += weight(t.root)
	}

	return card(t.root, weight(t.root.left), start, stop)
}

func (t *BinaryTree) Card(i int) (int, error) {
	l := t.Range(i, i)

	if len(l) == 0 {
		return 0, errors.New("invalid index")
	}

	return l[0], nil
}

func (t *BinaryTree) Details() {
	inOrder(t.root, 0, func(a *node, d int) {
		padding := strings.Repeat("    ", d)
		fmt.Printf("%s(%d)\n", padding, a.value)
	})
}

type node struct {
	value  int
	height int
	weight int

	left  *node
	right *node
}

func newNode(value int) *node {
	return &node{
		value:  value,
		height: 0,
		weight: 1,
		left:   nil,
		right:  nil,
	}
}

func insert(a *node, value int) *node {
	if a == nil {
		return newNode(value)
	}

	if value < a.value {
		a.left = insert(a.left, value)
	} else {
		a.right = insert(a.right, value)
	}

	a.height = max(height(a.left), height(a.right)) + 1
	a.weight = weight(a.left) + weight(a.right) + 1

	return rebalance(a)
}

func remove(a *node, value int) *node {
	if a == nil {
		return nil
	}

	if value == a.value {
		if a.left == nil {
			return a.right
		}

		if a.right == nil {
			return a.left
		}

		// has 2 kids
		smallest := first(a.right)
		a.right = remove(a.right, smallest)
		a.value = smallest
	} else {
		if value < a.value {
			a.left = remove(a.left, value)
		} else {
			a.right = remove(a.right, value)
		}
	}

	a.height = max(height(a.left), height(a.right)) + 1
	a.weight = weight(a.left) + weight(a.right) + 1

	return rebalance(a)
}

func first(a *node) int {
	if a.left == nil {
		return a.value
	}

	return first(a.left)
}

func balance(a *node) int {
	return height(a.left) - height(a.right)
}

func rebalance(a *node) *node {
	if balance(a) > 1 {
		if balance(a.left) >= 0 {
			return rotateRight(a)
		} else {
			return rotateLeftRight(a)
		}
	}

	if balance(a) < -1 {
		if balance(a.right) <= 0 {
			return rotateLeft(a)
		} else {
			return rotateRightLeft(a)
		}
	}

	return a
}

func height(a *node) int {
	if a == nil {
		return -1
	}

	return a.height
}

func weight(a *node) int {
	if a == nil {
		return 0
	}

	return a.weight
}

func rotateLeft(a *node) *node {
	b := a.right

	a.right = b.left
	b.left = a

	a.height = max(height(a.left), height(a.right)) + 1
	b.height = max(height(b.left), height(b.right)) + 1

	a.weight = weight(a.left) + weight(a.right) + 1
	b.weight = weight(b.left) + weight(b.right) + 1

	return b
}

func rotateRight(a *node) *node {
	b := a.left

	a.left = b.right
	b.right = a

	a.height = max(height(a.left), height(a.right)) + 1
	b.height = max(height(b.left), height(b.right)) + 1

	a.weight = weight(a.left) + weight(a.right) + 1
	b.weight = weight(b.left) + weight(b.right) + 1

	return b
}

func rotateLeftRight(a *node) *node {
	a.left = rotateLeft(a.left)
	a = rotateRight(a)

	return a
}

func rotateRightLeft(a *node) *node {
	a.right = rotateRight(a.right)
	a = rotateLeft(a)

	return a
}

func card(a *node, i int, start int, stop int) []int {
	if a == nil {
		return []int{}
	}

	l := []int{}

	if start < i && a.left != nil {
		l = card(a.left, i-weight(a.left)+weight(a.left.left), start, stop)
	}

	if start <= i && i <= stop {
		l = append(l, a.value)
	}

	if i < stop && a.right != nil {
		l = append(l, card(a.right, i+weight(a.right.left)+1, start, stop)...)
	}

	return l
}

func rank(a *node, value int) (int, error) {
	if a == nil {
		return 0, errors.New("not found")
	}

	if value == a.value {
		return weight(a.left), nil
	}

	if value < a.value {
		return rank(a.left, value)
	} else {
		r, err := rank(a.right, value)

		if err == nil {
			r = r + weight(a.left) + 1
		}

		return r, err
	}
}

func inOrder(a *node, d int, f func(*node, int)) {
	if a == nil {
		return
	}

	inOrder(a.left, d+1, f)
	f(a, d)
	inOrder(a.right, d+1, f)
}

func value(a *node) int {
	if a == nil {
		return -9999
	}

	return a.value
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
