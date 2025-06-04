// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
// Practice 1.10 URL cache: modify fetchall to fetch each URL twice,
// Practice 1.11 fetchall with goroutines: modify fetchall to start a goroutine for each URL,
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // start a goroutine
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()

	// Create a unique filename to save the response content
	fileName := fmt.Sprintf("response-%d-%s.html", time.Now().UnixNano(), urlToFilename(url))
	outFile, err := os.Create(fileName)
	if err != nil {
		ch <- fmt.Sprintf("error creating file %s: %v", fileName, err)
		return
	}
	defer outFile.Close()

	nbytes, err := io.Copy(outFile, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s and writing to %s: %v", url, fileName, err)
		return
	}

	secs := time.Since(start).Seconds()
	// Send the file path, time, and size back to the channel
	ch <- fmt.Sprintf("%.2fs  %7d  %s (saved to %s)", secs, nbytes, url, fileName)
}

// Helper function: converts a URL to a safe filename (very simplified, may not handle all URLs)
func urlToFilename(url string) string {
	// Remove protocol part
	s := strings.ReplaceAll(url, "http://", "")
	s = strings.ReplaceAll(s, "https://", "")
	// Replace characters not allowed in filenames
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, "?", "_")
	s = strings.ReplaceAll(s, "=", "_")
	s = strings.ReplaceAll(s, "&", "_")
	// Optionally truncate long filenames
	if len(s) > 50 {
		s = s[:50]
	}
	return s
}

//!-
