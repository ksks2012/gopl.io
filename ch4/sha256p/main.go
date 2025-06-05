// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 83.

// The sha256 command computes the SHA256 hash (an array) of a string.
// Practice 4.1: Compare two SHA256 hashes and count the number of different bits.
// Practice 4.2: Write a program that by default prints the SHA256 hash of its standard input,
// and supports a command-line flag to print SHA384 or SHA512 hash instead.

package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func countDifferentBits(c1, c2 [32]byte) int {
	count := 0
	for i := 0; i < len(c1); i++ {
		diffByte := c1[i] ^ c2[i]
		for j := 0; j < 8; j++ {
			if (diffByte>>j)&1 == 1 {
				count++
			}
		}
	}
	return count
}

func main() {
	// c1 := sha256.Sum256([]byte("x"))
	// c2 := sha256.Sum256([]byte("X"))
	// fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)

	// Output:
	// 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	// 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	// false
	// [32]uint8

	// diffBits := countDifferentBits(c1, c2)
	// fmt.Printf("Number of different bits: %d\n", diffBits)

	fmt.Printf("==================================================================\n")

	hashType := flag.String("hash", "sha256", "Specify hash algorithm: sha256, sha384, or sha512")

	flag.Parse()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading from standard input: %v", err)
	}

	switch *hashType {
	case "sha256":
		h := sha256.Sum256(data)
		fmt.Printf("%x\n", h)
	case "sha384":
		h := sha512.Sum384(data)
		fmt.Printf("%x\n", h)
	case "sha512":
		h := sha512.Sum512(data)
		fmt.Printf("%x\n", h)
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid hash type specified: %s\n", *hashType)
		fmt.Fprintln(os.Stderr, "Available types: sha256, sha384, sha512")
		os.Exit(1)
	}
}
