// Practice 5.8: Write a program that takes a file name and an ID as command-line arguments, fetches the HTML from the file name, and prints the first element with the given ID. If no such element exists, it should print a message indicating that the element was not found. Use the `golang.org/x/net/html` package to parse the HTML.
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file name> <id>\n", os.Args[0])
		os.Exit(1)
	}
	fileName := os.Args[1]
	targetID := os.Args[2]

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

	foundElement := ElementByID(doc, targetID)
	if foundElement != nil {
		fmt.Printf("Found element with ID '%s': <%s", targetID, foundElement.Data)
		for _, a := range foundElement.Attr {
			fmt.Printf(" %s=\"%s\"", a.Key, a.Val)
		}
		fmt.Printf(">\n")
	} else {
		fmt.Printf("Element with ID '%s' not found.\n", targetID)
	}
}

// ElementByID searches for the first HTML element with the given id.
// It returns the found node or nil if not found.
func ElementByID(doc *html.Node, id string) *html.Node {
	var foundNode *html.Node

	// pre function: check if the node is the target node before visiting its children
	pre := func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					foundNode = n
					return false
				}
			}
		}
		return true
	}

	// post function: executed after visiting child nodes; should also terminate if the node has been found
	post := func(n *html.Node) bool {
		return foundNode == nil
	}

	forEachNode(doc, pre, post)

	return foundNode
}

// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
// It returns true if the traversal should continue, false otherwise.
func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) bool {
	if pre != nil {
		if !pre(n) {
			return false
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !forEachNode(c, pre, post) {
			return false
		}
	}

	if post != nil {
		if !post(n) {
			return false
		}
	}
	return true
}
