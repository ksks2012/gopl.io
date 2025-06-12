// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// The sum program demonstrates a variadic function.
package sump

import (
	"fmt"
	"math"
)

// !+sum
func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

//!-sum

// Practice 5.15: Write variadic functions max and min similar to sum.
// Consider how max and min should handle no arguments, then write versions
// that take at least one argument.

// max_flexible finds the maximum value among the given integers.
// If no arguments are provided, it returns math.MinInt (or you might return 0, or an error).
func max_flexible(vals ...int) int {
	if len(vals) == 0 {
		fmt.Println("Warning: max_flexible called with no arguments, returning math.MinInt.")
		return math.MinInt
	}

	m := vals[0]
	for _, val := range vals {
		if val > m {
			m = val
		}
	}
	return m
}

// min_flexible finds the minimum value among the given integers.
// If no arguments are provided, it returns math.MaxInt (or you might return 0, or an error).
func min_flexible(vals ...int) int {
	if len(vals) == 0 {
		fmt.Println("Warning: min_flexible called with no arguments, returning math.MaxInt.")
		return math.MaxInt
	}

	m := vals[0]
	for _, val := range vals {
		if val < m {
			m = val
		}
	}
	return m
}

// max_atLeastOne finds the maximum value among at least one given integer.
// It requires at least one argument.
func max_atLeastOne(first int, vals ...int) int {
	m := first
	for _, val := range vals {
		if val > m {
			m = val
		}
	}
	return m
}

// min_atLeastOne finds the minimum value among at least one given integer.
// It requires at least one argument.
func min_atLeastOne(first int, vals ...int) int {
	m := first
	for _, val := range vals {
		if val < m {
			m = val
		}
	}
	return m
}
