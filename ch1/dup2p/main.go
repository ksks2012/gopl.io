// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
// Practice: Print file names if a line appears in multiple files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

type fileSet map[string]struct{}

type lineCount struct {
	line  string
	files fileSet
	count int
}

func main() {
	counts := make(map[string]lineCount)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		count := n.count
		if count > 1 {
			fmt.Printf("%d\t%s\t", count, line)
			i := 0
			for file := range n.files {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", file)
				i++
			}
			fmt.Println()
		}
	}
}

func countLines(f *os.File, counts map[string]lineCount) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		lc, ok := counts[line]
		if !ok {
			lc = lineCount{
				line:  line,
				files: make(fileSet),
				count: 1,
			}
		} else {
			lc.count++
		}
		lc.files[f.Name()] = struct{}{}
		counts[line] = lc
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!-
