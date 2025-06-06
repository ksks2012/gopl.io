package revp_test

import (
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
