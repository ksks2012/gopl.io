package revp_test

import (
	"bytes"
	"reflect"
	"testing"
	"unicode"

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

func TestRemoveAdjacentDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "no duplicates",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "all duplicates",
			input:    []int{7, 7, 7, 7},
			expected: []int{7},
		},
		{
			name:     "adjacent duplicates",
			input:    []int{1, 1, 2, 2, 2, 3, 3, 4, 5, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "non-adjacent duplicates",
			input:    []int{1, 2, 1, 2, 1},
			expected: []int{1, 2, 1, 2, 1},
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{42},
			expected: []int{42},
		},
		{
			name:     "long run at end",
			input:    []int{1, 2, 3, 4, 4, 4, 4},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "long run at start",
			input:    []int{5, 5, 5, 2, 3},
			expected: []int{5, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := make([]int, len(tt.input))
			copy(inputCopy, tt.input)
			got := revp.RemoveAdjacentDuplicatesInt(inputCopy)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RemoveAdjacentDuplicates(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRemoveAdjacentDuplicatesString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "all duplicates",
			input:    []string{"x", "x", "x", "x"},
			expected: []string{"x"},
		},
		{
			name:     "adjacent duplicates",
			input:    []string{"a", "a", "b", "b", "b", "c", "c", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "non-adjacent duplicates",
			input:    []string{"a", "b", "a", "b", "a"},
			expected: []string{"a", "b", "a", "b", "a"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []string{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "long run at end",
			input:    []string{"a", "b", "c", "c", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "long run at start",
			input:    []string{"z", "z", "z", "y", "x"},
			expected: []string{"z", "y", "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := make([]string, len(tt.input))
			copy(inputCopy, tt.input)
			got := revp.RemoveAdjacentDuplicatesString(inputCopy)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RemoveAdjacentDuplicatesString(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRemoveAdjacentSpace(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "no spaces",
			input:    []byte("abc"),
			expected: []byte("abc"),
		},
		{
			name:     "single space",
			input:    []byte("a b c"),
			expected: []byte("a b c"),
		},
		{
			name:     "multiple adjacent spaces",
			input:    []byte("a  b   c"),
			expected: []byte("a b c"),
		},
		{
			name:     "spaces at start and end",
			input:    []byte("   abc   "),
			expected: []byte(" abc "),
		},
		{
			name:     "all spaces",
			input:    []byte("     "),
			expected: []byte(" "),
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: []byte(""),
		},
		{
			name:     "single character",
			input:    []byte("x"),
			expected: []byte("x"),
		},
		{
			name:     "tab and space",
			input:    []byte("a \t  b"),
			expected: []byte("a b"),
		},
		{
			name:     "mixed unicode spaces",
			input:    []byte("a\u00A0\u00A0b"),
			expected: []byte("a b"),
		},
		{
			name:     "newline and spaces",
			input:    []byte("a\n\n b"),
			expected: []byte("a b"),
		},
		{
			name:     "complex mixed spaces",
			input:    []byte(" Hello \t\tWorld\n\r  Go! "),
			expected: []byte(" Hello World Go! "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := make([]byte, len(tt.input))
			copy(inputCopy, tt.input)
			got := revp.RemoveAdjacentSpace(inputCopy)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RemoveAdjacentSpace(%q) = %q (len %d), want %q (len %d)",
					tt.input, got, len(got), tt.expected, len(tt.expected))
			}

			// Check for adjacent spaces in the result
			if len(got) > 0 {
				runesGot := bytes.Runes(got)
				for i := 1; i < len(runesGot); i++ {
					if unicode.IsSpace(runesGot[i]) && unicode.IsSpace(runesGot[i-1]) {
						t.Errorf("RemoveAdjacentSpace(%q) result has adjacent spaces at rune index %d: %q",
							tt.input, i, string(got))
						break
					}
				}
			}
		})
	}
}

func TestReverseUTF8(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "ASCII",
			input:    "abcdef",
			expected: "fedcba",
		},
		{
			name:     "single rune",
			input:    "世",
			expected: "世",
		},
		{
			name:     "multi-rune",
			input:    "hello 世界",
			expected: "界世 olleh",
		},
		{
			name:     "emoji",
			input:    "🙂🙃",
			expected: "🙃🙂",
		},
		{
			name:     "mixed ASCII and emoji",
			input:    "A🙂B🙃C",
			expected: "C🙃B🙂A",
		},
		{
			name:     "combining characters",
			input:    "e\u0301cole", // "école" with combining acute
			expected: "eloće",
		},
		{
			name:     "multi-byte runes",
			input:    "¡Hola, 世界!",
			expected: "!界世 ,aloH¡",
		},
		{
			name:     "surrogate pairs",
			input:    "𐍈𐍉", // Gothic letters, 4-byte UTF-8
			expected: "𐍉𐍈",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputBytes := []byte(tt.input)
			got := revp.ReverseUTF8(inputBytes)
			if string(got) != tt.expected {
				t.Errorf("ReverseUTF8(%q) = %q, want %q", tt.input, string(got), tt.expected)
			}
		})
	}
}
