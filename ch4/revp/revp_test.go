package revp_test

import (
	"reflect"
	"testing"

	"gopl.io/ch4/revp"
)

// Assume these functions are implemented in the revp package
// ReverseArrayPtr(a *[6]int)
// ReverseSlice(s []int)
func TestReverseArrayPtr(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5}
	expected := [...]int{5, 4, 3, 2, 1, 0}
	revp.ReverseArrayPtr(&a)
	if a != expected {
		t.Errorf("ReverseArrayPtr failed: got %v, want %v", a, expected)
	}
}

func TestReverseSlice(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	expected := []int{5, 4, 3, 2, 1, 0}
	revp.ReverseSlice(s)
	for i := range s {
		if s[i] != expected[i] {
			t.Errorf("ReverseSlice failed: got %v, want %v", s, expected)
			break
		}
	}
}

func TestRotate(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	got := revp.Rotate(s, 2)
	want := []int{2, 3, 4, 5, 0, 1}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Rotate(s, 2) = %v, want %v", got, want)
			break
		}
	}
	got = revp.Rotate(got, 4)
	want = []int{0, 1, 2, 3, 4, 5}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Rotate(s, 4) = %v, want %v", got, want)
			break
		}
	}
}

func TestRotateInPlace(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		n        int
		expected []int
	}{
		{
			name:     "rotate by 0",
			input:    []int{0, 1, 2, 3, 4, 5},
			n:        0,
			expected: []int{0, 1, 2, 3, 4, 5},
		},
		{
			name:     "rotate by 2",
			input:    []int{0, 1, 2, 3, 4, 5},
			n:        2,
			expected: []int{2, 3, 4, 5, 0, 1},
		},
		{
			name:     "rotate by len(s)",
			input:    []int{0, 1, 2, 3, 4, 5},
			n:        6,
			expected: []int{0, 1, 2, 3, 4, 5},
		},
		{
			name:     "rotate by negative",
			input:    []int{0, 1, 2, 3, 4, 5},
			n:        -2,
			expected: []int{4, 5, 0, 1, 2, 3},
		},
		{
			name:     "rotate by more than len(s)",
			input:    []int{0, 1, 2, 3, 4, 5},
			n:        8,
			expected: []int{2, 3, 4, 5, 0, 1},
		},
		{
			name:     "empty slice",
			input:    []int{},
			n:        3,
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{42},
			n:        1,
			expected: []int{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make([]int, len(tt.input))
			copy(s, tt.input)
			revp.RotateInPlace(s, tt.n)
			if !reflect.DeepEqual(s, tt.expected) {
				t.Errorf("RotateInPlace(%v, %d) = %v, want %v", tt.input, tt.n, s, tt.expected)
			}
		})
	}
}
