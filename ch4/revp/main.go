// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 86.

// Rev reverses an array using a pointer to the array.
// Practice 4.3: Rev reverses a slice using a pointer to the slice.
// Practice 4.4: Implement a rotate function that rotates a slice left by n positions in one loop.
// Practice 4.5: Implement functions to remove adjacent duplicates for int and string slices
// Practice 4.6: Implement a function to remove adjacent spaces from a byte slice
// Practice 4.7: Implement a function to reverse the bytes with UTF-8 encoding
package revp

import (
	"bytes"
	"unicode"
)

func ReverseArrayPtr(ptr *[6]int) {
	for i, j := 0, len(*ptr)-1; i < j; i, j = i+1, j-1 {
		(*ptr)[i], (*ptr)[j] = (*ptr)[j], (*ptr)[i]
	}
}

func ReverseSlice(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Rotate(s []int, n int) []int {
	if len(s) == 0 || n == 0 {
		return s
	}

	n %= len(s)
	if n < 0 {
		n += len(s)
	}

	if n == 0 {
		return s
	}

	rotated := make([]int, len(s))

	copy(rotated, s[n:])

	copy(rotated[len(s)-n:], s[:n])

	return rotated
}

func RotateInPlace(s []int, n int) {
	if len(s) == 0 || n == 0 {
		return
	}

	n %= len(s)
	if n < 0 {
		n += len(s)
	}

	if n == 0 {
		return
	}

	ReverseSlice(s[:n])
	ReverseSlice(s[n:])
	ReverseSlice(s)
}

func RemoveAdjacentDuplicatesInt(s []int) []int {
	if len(s) == 0 {
		return s
	}
	j := 0
	for i := 1; i < len(s); i++ {
		if s[i] != s[j] {
			j++
			s[j] = s[i]
		}
	}
	return s[:j+1]
}

func RemoveAdjacentDuplicatesString(s []string) []string {
	if len(s) == 0 {
		return s
	}
	j := 0
	for i := 1; i < len(s); i++ {
		if s[i] != s[j] {
			j++
			s[j] = s[i]
		}
	}
	return s[:j+1]
}

func RemoveAdjacentSpace(b []byte) []byte {
	if len(b) == 0 {
		return b
	}

	runes := bytes.Runes(b)
	if len(runes) == 0 {
		return b[:0]
	}

	wasPrevSpace := unicode.IsSpace(runes[0])
	if wasPrevSpace {
		// Normalize to a single ASCII space
		runes[0] = ' '
	}

	j := 0
	for i := 1; i < len(runes); i++ {
		isCurrentSpace := unicode.IsSpace(runes[i])

		if isCurrentSpace {
			if !wasPrevSpace {
				j++
				// Always normalize to a single ASCII space
				runes[j] = ' '
				wasPrevSpace = true
			}
		} else {
			j++
			runes[j] = runes[i]
			// Reset flag as the last written char is not a space
			wasPrevSpace = false
		}
	}

	return []byte(string(runes[:j+1]))
}

func ReverseBytesInPlace(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func ReverseUTF8(b []byte) []byte {
	if len(b) == 0 {
		return b
	}

	runes := bytes.Runes(b)
	if len(runes) == 0 {
		return b[:0]
	}

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	copy(b, string(runes))
	b = b[:len(string(runes))]

	return b
}
