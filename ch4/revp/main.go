// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 86.

// Rev reverses an array using a pointer to the array.
// Practice 4.3: Rev reverses a slice using a pointer to the slice.
// Practice 4.4: Implement a rotate function that rotates a slice left by n positions in one loop.
// Practice 4.5: Implement functions to remove adjacent duplicates for int and string slices
package revp

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
