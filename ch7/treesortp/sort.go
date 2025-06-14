// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 101.

// Package treesort provides insertion sort using an unbalanced binary tree.
// Practice 7.3: Implement String() for the Tree type.
package treesort

import "strconv"

// !+
type Tree struct {
	Value       int
	Left, Right *Tree
}

// Sort sorts values in place.
func Sort(values []int) {
	var root *Tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.
func appendValues(values []int, t *Tree) []int {
	if t != nil {
		values = appendValues(values, t.Left)
		values = append(values, t.Value)
		values = appendValues(values, t.Right)
	}
	return values
}

func add(t *Tree, value int) *Tree {
	if t == nil {
		// Equivalent to return &tree{value: value}.
		t = new(Tree)
		t.Value = value
		return t
	}
	if value < t.Value {
		t.Left = add(t.Left, value)
	} else {
		t.Right = add(t.Right, value)
	}
	return t
}

//!-

func (t *Tree) String() string {
	if t == nil {
		return ""
	}
	return t.Left.String() + " " + strconv.Itoa(t.Value) + " " + t.Right.String()
}
