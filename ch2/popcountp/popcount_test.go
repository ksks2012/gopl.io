// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package popcountp_test

import (
	"testing"

	popcount "gopl.io/ch2/popcountp"
)

// -- Alternative implementations --

func BitCount(x uint64) int {
	// Hacker's Delight, Figure 5-2.
	x = x - ((x >> 1) & 0x5555555555555555)
	x = (x & 0x3333333333333333) + ((x >> 2) & 0x3333333333333333)
	x = (x + (x >> 4)) & 0x0f0f0f0f0f0f0f0f
	x = x + (x >> 8)
	x = x + (x >> 16)
	x = x + (x >> 32)
	return int(x & 0x7f)
}

func PopCountByClearing(x uint64) int {
	n := 0
	for x != 0 {
		x = x & (x - 1) // clear rightmost non-zero bit
		n++
	}
	return n
}

func PopCountByShifting(x uint64) int {
	n := 0
	for i := uint(0); i < 64; i++ {
		if x&(1<<i) != 0 {
			n++
		}
	}
	return n
}

// -- Benchmarks --

func BenchmarkPopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount.PopCount(0x1234567890ABCDEF)
	}
}

func BenchmarkBitCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BitCount(0x1234567890ABCDEF)
	}
}

func BenchmarkPopCountByClearing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountByClearing(0x1234567890ABCDEF)
	}
}

func BenchmarkPopCountByShifting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountByShifting(0x1234567890ABCDEF)
	}
}
func BenchmarkPopCountLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount.PopCountLoop(0x1234567890ABCDEF)
	}
}

// cpu: 3.7GHz i5-12600K
// $ go test -cpu=4 -bench=. gopl.io/ch2/popcountp
// BenchmarkPopCount-4             	1000000000	         0.1198 ns/op
// BenchmarkBitCount-4             	1000000000	         0.1248 ns/op
// BenchmarkPopCountByClearing-4   	76853898	        16.00 ns/op
// BenchmarkPopCountByShifting-4   	64625599	        18.90 ns/op
// BenchmarkPopCountLoop-4         	64130601	        19.26 ns/op

// cpu: 3.7GHz i5-12600K
// $ go test -cpu=8 -bench=. gopl.io/ch2/popcountp
// BenchmarkPopCount-8             	1000000000	         0.1161 ns/op
// BenchmarkBitCount-8             	1000000000	         0.1221 ns/op
// BenchmarkPopCountByClearing-8   	72335458	        15.87 ns/op
// BenchmarkPopCountByShifting-8   	63105985	        19.22 ns/op
// BenchmarkPopCountLoop-8         	60994216	        19.34 ns/op

// Go 1.6, 2.67GHz Xeon
// $ go test -cpu=4 -bench=. gopl.io/ch2/popcount
// BenchmarkPopCount-4                  200000000         6.30 ns/op
// BenchmarkBitCount-4                  300000000         4.15 ns/op
// BenchmarkPopCountByClearing-4        30000000         45.2 ns/op
// BenchmarkPopCountByShifting-4        10000000        153 ns/op
//
// Go 1.6, 2.5GHz Intel Core i5
// $ go test -cpu=4 -bench=. gopl.io/ch2/popcount
// BenchmarkPopCount-4                  200000000         7.52 ns/op
// BenchmarkBitCount-4                  500000000         3.36 ns/op
// BenchmarkPopCountByClearing-4        50000000         34.3 ns/op
// BenchmarkPopCountByShifting-4        20000000        108 ns/op
//
// Go 1.7, 3.5GHz Xeon
// $ go test -cpu=4 -bench=. gopl.io/ch2/popcount
// BenchmarkPopCount-12                 2000000000        0.28 ns/op
// BenchmarkBitCount-12                 2000000000        0.27 ns/op
// BenchmarkPopCountByClearing-12       100000000        18.5 ns/op
// BenchmarkPopCountByShifting-12       20000000         70.1 ns/op
