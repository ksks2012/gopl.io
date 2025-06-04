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
	"strings"
)

type fileSet map[string]struct{}

type lineCount struct {
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
			// NOTE: for limited resources, you might want to use
			defer f.Close()
			countLines(f, counts)
			// NOTE: avoid stacking too many open files
			// f.Close()
		}
	}
	for line, n := range counts {
		count := n.count
		if count > 1 {
			fmt.Printf("%d\t%s\t", count, line)

			var fileNames []string
			for file := range n.files {
				fileNames = append(fileNames, file)
			}
			fmt.Println(strings.Join(fileNames, ", "))
		}
	}
}

func countLines(f *os.File, counts map[string]lineCount) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		entry := counts[line]
		if entry.files == nil {
			entry.files = make(fileSet)
		}
		entry.count++
		entry.files[f.Name()] = struct{}{}
		// Store the updated entry back in the map
		counts[line] = entry
	}
	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "dup2: reading %s: %v\n", f.Name(), err)
	}
}

//!-
