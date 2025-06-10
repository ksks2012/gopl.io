// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 122.
//!+main

// Findlinks1 prints the links in an HTML document read from standard input.
// Practice 5.1: Recursively visit the children of a node, using NextSibling and FirstChild.
// Practice 5.2: Write a function that counts the number of times each element appears in an HTML tree.
// Practice 5.3: Write a function that visits the text attributes of an HTML tree, skipping <script> and <style> elements.
// Practice 5.4: Write a function that collects the src attributes of <img> and <script> elements, and the href attributes of <link> elements.
package main

import (
	"fmt"
	"os"
	"strings"

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

	// Example of visiting text attributes
	texts := visitText(nil, doc)
	fmt.Println("\nText found:")
	for _, text := range texts {
		fmt.Println(text)
	}

	assetMap := map[string]string{
		"img":    "src",
		"script": "src",
		"":       "href",
	}
	otherAssets := findAssets(nil, doc, assetMap)
	fmt.Println("\nOther assets found (img src, script src, link href):")
	for _, asset := range otherAssets {
		fmt.Println(asset)
	}

}

// findAssets is a generalized function to collect specific attributes from various element types.
// The assetMap maps element names (e.g., "img", "script") to their attribute keys (e.g., "src", "href").
func findAssets(assets []string, n *html.Node, assetMap map[string]string) []string {
	if n.Type == html.ElementNode {
		attrKey, exists := assetMap[n.Data]
		if exists {
			for _, a := range n.Attr {
				if a.Key == attrKey {
					assets = append(assets, a.Val)
					break
				}
			}
		}
	}

	if n.FirstChild != nil {
		assets = findAssets(assets, n.FirstChild, assetMap)
	}

	if n.NextSibling != nil {
		assets = findAssets(assets, n.NextSibling, assetMap)
	}

	return assets
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

// It skips <script> and <style> elements.
func visitText(texts []string, n *html.Node) []string {
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		return texts
	}

	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			texts = append(texts, text)
		}
	}

	// Recursive
	if n.FirstChild != nil {
		texts = visitText(texts, n.FirstChild)
	}

	if n.NextSibling != nil {
		texts = visitText(texts, n.NextSibling)
	}
	return texts
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
