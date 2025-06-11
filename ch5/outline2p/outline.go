// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Practice 5.7: Complete the startElement and endElement functions to be a generic HTML pretty-printer.
// It should output comments, text nodes, and attributes for each element (<a href='...'>).
// Use the compact form for elements without children (i.e., <img/> instead of <img></img>).
// Write tests to verify the program's output format is correct (see Chapter 11).

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var depth int

func main() {
	depth = 0

	fileName := os.Args[1]
	if fileName == "" {
		fmt.Println("Please provide a file name.")
		return
	}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	doc, err := html.Parse(file)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}
	forEachNode(doc, printNodePre, printNodePost)
}

func outline(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("parsing HTML: %s", err)
	}

	forEachNode(doc, printNodePre, printNodePost)

	return nil
}

// forEachNode traverses the HTML tree and calls the callback functions before (pre) and after (post) visiting each node.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

// The printNodePre function is called before visiting a node, used to print start tags, comments, or text content.
func printNodePre(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		isSelfClosingTag := (n.FirstChild == nil && n.Data != "script" && n.Data != "style")

		fmt.Printf("%*s<%s", depth*2, "", n.Data)
		for _, a := range n.Attr {
			fmt.Printf(" %s=\"%s\"", a.Key, a.Val)
		}
		// Self-closing format
		// For self-closing tags, there is no need to increase depth, nor will the post function print an end tag
		// So there is no need for depth++ here or depth-- in the post function
		if isSelfClosingTag {
			fmt.Printf("/>\n")
		} else {
			fmt.Printf(">\n")
			depth++
		}

	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Printf("%*s%s\n", depth*2, "", text)
		}
	case html.CommentNode:
		fmt.Printf("%*s%s\n", depth*2, "", n.Data)
	case html.DocumentNode:
		fmt.Printf("<!DOCTYPE html>\n")
	}
}

func printNodePost(n *html.Node) {
	if n.Type == html.ElementNode {
		// Check if this is a self-closing tag; if so, do not print the end tag
		isSelfClosingTag := (n.FirstChild == nil && n.Data != "script" && n.Data != "style")
		if isSelfClosingTag {
			// Self-closing tags are already handled in printNodePre, nothing to do here
			return
		}

		depth-- // Leaving the child node, decrease the depth
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}
