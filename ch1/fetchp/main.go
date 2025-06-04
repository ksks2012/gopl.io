// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 16.
//!+

// Fetch prints the content found at each specified URL.
// Practice 1.7: Replace ioutil.ReadAll() by using io.Copy(dst, src)
// Practice 1.8: Ensure the URL starts with "http://" or "https://"
// Practice 1.9: Modify fetch to print out the status code of the HTTP protocol, which can be obtained from the resp.Status variable.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		// Ensure the URL starts with "http://" or "https://"
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				fmt.Fprintf(os.Stderr, "fetch: error closing body for %s: %v\n", url, closeErr)
			}
		}()

		fmt.Printf("Status Code for %s: %s\n", url, resp.Status)

		out := os.Stdout // or any other io.Writer
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

//!-
