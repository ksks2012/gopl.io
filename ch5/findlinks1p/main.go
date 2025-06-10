// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 122.
//!+main

// Findlinks1 prints the links in an HTML document read from standard input.
// Practice 5.1: Recursively visit the children of a node, using NextSibling and FirstChild.
// Practice 5.2: Write a function that counts the number of times each element appears in an HTML tree.
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	filepath := os.Args[1]
	if filepath == "" {
		fmt.Fprintln(os.Stderr, "Usage: findlinks1 <file>")
		os.Exit(1)
	}
	htmlFile, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	defer htmlFile.Close()

	doc, err := html.Parse(htmlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}

	links := visit(nil, doc)
	fmt.Println("Links found:")
	for _, link := range links {
		fmt.Println(link)
	}

	elementCounts := make(map[string]int)
	countElements(elementCounts, doc)

	fmt.Println("\nElement counts:")
	for name, count := range elementCounts {
		fmt.Printf("%s: %d\n", name, count)
	}
}

// countElements recursively traverses the HTML tree and records the occurrence count of each element name.
// It directly updates the provided map.
func countElements(counts map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		counts[n.Data]++
	}

	if n.FirstChild != nil {
		countElements(counts, n.FirstChild)
	}

	if n.NextSibling != nil {
		countElements(counts, n.NextSibling)
	}
}

// visit appends to links each link found in n and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	if n.FirstChild != nil {
		links = visit(links, n.FirstChild)
	}

	if n.NextSibling != nil {
		links = visit(links, n.NextSibling)
	}

	return links
}

/*
//!+html
package html

type Node struct {
	Type                    NodeType
	Data                    string
	Attr                    []Attribute
	FirstChild, NextSibling *Node
}

type NodeType int32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

type Attribute struct {
	Key, Val string
}

func Parse(r io.Reader) (*Node, error)
//!-html
*/
