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

func BenchmarkPopCountKernighan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount.PopCountKernighan(0x1234567890ABCDEF)
	}
}

// Go 1.6, 3.7GHz i5-12600K
// $ go test -cpu=4 -bench=. gopl.io/ch2/popcountp
// BenchmarkPopCount-4             	1000000000	         0.1187 ns/op
// BenchmarkBitCount-4             	1000000000	         0.1309 ns/op
// BenchmarkPopCountByClearing-4   	80153720	        15.50 ns/op
// BenchmarkPopCountByShifting-4   	64560848	        18.26 ns/op
// BenchmarkPopCountLoop-4         	66691446	        18.58 ns/op
// BenchmarkPopCountKernighan-4    	76319629	        15.82 ns/op

// Go 1.6, 3.7GHz i5-12600K
// $ go test -cpu=8 -bench=. gopl.io/ch2/popcountp
// BenchmarkPopCount-8             	1000000000	         0.1141 ns/op
// BenchmarkBitCount-8             	1000000000	         0.1157 ns/op
// BenchmarkPopCountByClearing-8   	76763002	        16.05 ns/op
// BenchmarkPopCountByShifting-8   	64865430	        18.06 ns/op
// BenchmarkPopCountLoop-8         	64134643	        18.24 ns/op
// BenchmarkPopCountKernighan-8    	70834881	        15.25 ns/op
