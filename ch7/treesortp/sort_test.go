// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package treesort_test

import (
	"math/rand"
	"sort"
	"testing"

	treesort "gopl.io/ch7/treesortp"
)

func TestSort(t *testing.T) {
	data := make([]int, 50)
	for i := range data {
		data[i] = rand.Int() % 50
	}
	treesort.Sort(data)
	if !sort.IntsAreSorted(data) {
		t.Errorf("not sorted: %v", data)
	}
}
func TestTreeString(t *testing.T) {
	tests := []struct {
		name   string
		build  func() *treesort.Tree
		expect string
	}{
		{
			name: "nil tree",
			build: func() *treesort.Tree {
				return nil
			},
			expect: "",
		},
		{
			name: "single node",
			build: func() *treesort.Tree {
				var t treesort.Tree
				t.Value = 42
				return &t
			},
			expect: " 42 ",
		},
		{
			name: "left and right children",
			build: func() *treesort.Tree {
				root := &treesort.Tree{Value: 2}
				root.Left = &treesort.Tree{Value: 1}
				root.Right = &treesort.Tree{Value: 3}
				return root
			},
			expect: " 1  2  3 ",
		},
		{
			name: "unbalanced tree",
			build: func() *treesort.Tree {
				root := &treesort.Tree{Value: 10}
				root.Left = &treesort.Tree{Value: 5}
				root.Left.Left = &treesort.Tree{Value: 2}
				root.Right = &treesort.Tree{Value: 20}
				return root
			},
			expect: " 2  5  10  20 ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.build().String()
			if got != tt.expect {
				t.Errorf("String() = %q, want %q", got, tt.expect)
			}
		})
	}
}
