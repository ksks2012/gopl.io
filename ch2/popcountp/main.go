// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 45.

// (Package doc comment intentionally malformed to demonstrate golint.)
// Pracice 2.3: Change the PopCount function to use a loop instead of a sequence of shifts and masks.
// Practice 2.4: Implement PopCount using bit shifts and masks, which counts the number of set bits by checking each bit position.
// Practice 2.5: Implement PopCount using Kernighan's algorithm, which counts the number of set bits by repeatedly clearing the least significant bit set.

// !+
package popcountp

// pc[i] is the population count of i.
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountLoop(x uint64) int {
	n := 0
	for i := 0; i < 8; i++ {
		n += int(pc[byte(x>>(i*8))])
	}
	return n
}

func PopCountBitShift(x uint64) int {
	n := 0
	for i := 0; i < 64; i++ {
		if (x>>i)&1 == 1 {
			n += 1
		}
	}
	return n
}

// PopCountKernighan returns the population count of x using Kernighan's algorithm.
func PopCountKernighan(x uint64) int {
	n := 0
	for x != 0 {
		x = x & (x - 1)
		n++
	}
	return n
}

//!-
