// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 97.
//!+

// Charcount computes counts of Unicode characters.
// Practice 4.8: Add CharCount for counting letters, digits, and numbers.
// Practice 4.9:
package charcountp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int)    // counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

//!-

func CharCount(mp map[rune]int) {
	letterCount := 0
	digitCount := 0
	numberCount := 0
	for c, n := range mp {
		if unicode.IsLetter(c) {
			letterCount += n
		}
		if unicode.IsDigit(c) {
			digitCount += n
		}
		if unicode.IsNumber(c) {
			numberCount += n
		}
	}
	fmt.Printf("Letters: %d\n", letterCount)
	fmt.Printf("Digits: %d\n", digitCount)
	fmt.Printf("Numbers: %d\n", numberCount)
}

func WordFreq(s string) map[string]int {
	freq := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(s))

	scanner.Split(bufio.ScanWords) // set the split mode to words

	for scanner.Scan() {
		word := scanner.Text()
		freq[word]++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "WordFreq: reading input: %v\n", err)
		os.Exit(1)
	}
	return freq
}
