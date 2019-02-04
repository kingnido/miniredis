package main

import (
	"fmt"
	"strings"
)

type BinaryTree struct {
	root *node
}

func NewBinaryTree() *BinaryTree {
	return &BinaryTree{}
}

func (t *BinaryTree) Add(value int) error {
	t.root = insert(t.root, value)

	return nil
}

func (t *BinaryTree) InOrder(f func(value int)) {
	inOrder(t.root, func(a *node) {
		f(a.value)
	})
}

func (t *BinaryTree) Details() {
	inOrder(t.root, func(a *node) {
		padding := strings.Repeat("    ", t.root.height-a.height)
		fmt.Printf("%s{ v: %d, h: %d, w: %d}\n", padding, a.value, a.height, a.weight)
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

func insert(root *node, value int) *node {
	if root == nil {
		return newNode(value)
	}

	if value < root.value {
		root.left = insert(root.left, value)
	} else {
		root.right = insert(root.right, value)
	}

	root.height = max(height(root.left), height(root.right)) + 1
	root.weight = weight(root.left) + weight(root.right) + 1

	balance := height(root.left) - height(root.right)

	if balance > 1 {
		if value < root.left.value {
			return rotateRight(root)
		} else {
			return rotateLeftRight(root)
		}
	}

	if balance < -1 {
		if value > root.right.value {
			return rotateLeft(root)
		} else {
			return rotateRightLeft(root)
		}
	}

	return root
}

func delete(a *node, value int) *node {
	return nil
}

func height(root *node) int {
	if root == nil {
		return -1
	}

	return root.height
}

func weight(root *node) int {
	if root == nil {
		return 0
	}

	return root.weight
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

func inOrder(a *node, f func(*node)) {
	if a == nil {
		return
	}

	inOrder(a.left, f)
	f(a)
	inOrder(a.right, f)
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
