// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 173.

// Bytecounter demonstrates an implementation of io.Writer that counts bytes.
// Practice 7.1: Write WordCounter and LineCounter types that count words and lines.
package counter

import (
	"bufio"
	"bytes"
	"io"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	// Use bufio.Scanner to scan words
	// Convert []byte to bytes.Reader (which implements io.Reader)
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanWords)

	words := 0
	for scanner.Scan() { // Each call to Scan() finds one wor
		words++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	*c += WordCounter(words)
	return len(p), nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	// Use bufio.Scanner to scan lines
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanLines) // Set the scanner to split by lines

	lines := 0
	for scanner.Scan() {
		lines++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	*c += LineCounter(lines)
	return len(p), nil
}

type countingWriter struct {
	writer io.Writer
	count  *int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.writer.Write(p)
	*cw.count += int64(n)
	return n, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var c int64
	cw := &countingWriter{writer: w, count: &c}
	return cw, &c
}
