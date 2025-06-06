// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 86.

// Rev reverses an array using a pointer to the array.
// Practice 4.3: Rev reverses a slice using a pointer to the slice.
// Practice 4.4: Implement a rotate function that rotates a slice left by n positions in one loop.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	//!+array
	a := [...]int{0, 1, 2, 3, 4, 5}
	reverseArrayPtr(&a)
	fmt.Println("Array: ", a) // "[5 4 3 2 1 0]"
	s := []int{0, 1, 2, 3, 4, 5}
	reverseSlice(s)
	fmt.Println("Slice: ", s) // "[5 4 3 2 1 0]"

	s = rotate(s, 2)
	fmt.Println("Rotated Slice: ", s) // "[2 3 4 5 0 1]"')

	input := bufio.NewScanner(os.Stdin)
outer:
	for input.Scan() {
		var ints []int
		for _, s := range strings.Fields(input.Text()) {
			x, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue outer
			}
			ints = append(ints, int(x))
		}
		reverseSlice(ints)
		fmt.Printf("%v\n", ints)
	}
}

func reverseArrayPtr(ptr *[6]int) {
	for i, j := 0, len(*ptr)-1; i < j; i, j = i+1, j-1 {
		(*ptr)[i], (*ptr)[j] = (*ptr)[j], (*ptr)[i]
	}
}

func reverseSlice(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func rotate(s []int, n int) []int {
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
