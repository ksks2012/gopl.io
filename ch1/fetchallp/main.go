// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
// Practice 1.10 URL cache: modify fetchall to fetch each URL twice,
// Practice 1.11 fetchall with goroutines: modify fetchall to start a goroutine for each URL,
package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// The result struct is used to pass more detailed results through the channel
type result struct {
	url      string
	secs     float64
	nbytes   int64
	fileName string
	hash     string
	err      error
}

func main() {
	start := time.Now()
	// Use a buffered channel to handle all goroutines sending results at the same time
	// The buffer size is twice the number of command-line arguments, since each URL is requested twice
	ch := make(chan result, len(os.Args[1:])*2)

	var wg sync.WaitGroup

	// Store the hash values of each URL's two requests for comparison
	// key: URL, value: []string (two hash values)
	urlHashes := make(map[string][]string)
	var mu sync.Mutex // Protect concurrent access to urlHashes map

	urls := os.Args[1:]
	if len(urls) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: go run fetchall.go <URL1> <URL2> ...\n")
		os.Exit(1)
	}

	for _, url := range urls {
		wg.Add(2)
		go fetch(url, ch, &wg)
		go fetch(url, ch, &wg)
	}

	wg.Wait()

	// Receive and process all results from the channel
	close(ch)
	for res := range ch {
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "fetch error: %v\n", res.err)
		} else {
			fmt.Printf("%.2fs  %7d  %s (saved to %s, hash: %s)\n", res.secs, res.nbytes, res.url, res.fileName, res.hash)

			mu.Lock()
			urlHashes[res.url] = append(urlHashes[res.url], res.hash)
			mu.Unlock()
		}
	}

	fmt.Printf("\n--- Content Comparison Results ---\n")
	for url, hashes := range urlHashes {
		// Ensure there are two hashes to compare
		if len(hashes) == 2 {
			if hashes[0] == hashes[1] {
				fmt.Printf("URL: %s - Content consistent (Hash: %s)\n", url, hashes[0])
			} else {
				fmt.Printf("URL: %s - Content INCONSISTENT (Hash1: %s, Hash2: %s)\n", url, hashes[0], hashes[1])
			}
		} else {
			// If, due to some error, there are not two downloads, skip comparison
			fmt.Printf("URL: %s - Not enough data for comparison (received %d results)\n", url, len(hashes))
		}
	}

	fmt.Printf("\n%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()

	reqStart := time.Now()

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		ch <- result{url: url, err: fmt.Errorf("getting %s: %w", url, err)}
		return
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "fetch: error closing body for %s: %v\n", url, closeErr)
		}
	}()

	fmt.Printf("Status Code for %s: %s\n", url, resp.Status)

	// Create a unique filename to save the response content
	fileName := fmt.Sprintf("response-%d-%s.html", time.Now().UnixNano(), urlToFilename(url))
	outFile, err := os.Create(fileName)
	if err != nil {
		ch <- result{url: url, err: fmt.Errorf("creating file %s: %w", fileName, err)}
		return
	}
	defer outFile.Close()

	// Create a hash calculator
	hasher := sha256.New() // or md5.New()

	// io.MultiWriter sends written data to multiple Writers at the same time
	// Here, it writes to both outFile (the file) and hasher (the hash calculator)
	writer := io.MultiWriter(outFile, hasher)

	nbytes, err := io.Copy(writer, resp.Body)
	if err != nil {
		ch <- result{url: url, err: fmt.Errorf("reading %s and writing to %s: %w", url, fileName, err)}
		return
	}

	secs := time.Since(reqStart).Seconds()
	contentHash := fmt.Sprintf("%x", hasher.Sum(nil))

	ch <- result{
		url:      url,
		secs:     secs,
		nbytes:   nbytes,
		fileName: fileName,
		hash:     contentHash,
		err:      nil,
	}
}

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
