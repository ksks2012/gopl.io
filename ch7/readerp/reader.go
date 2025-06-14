// Practice 7.4: Write a StringReader type that implements io.Reader.
package reader

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

type StringReader struct {
	s string
	i int64
	// Additions for full io.Reader:
	// prevRune int
	// // last rune read for unread (not strictly needed for this exercise)
}

// Read reads data from the StringReader into p.
// It implements the io.Reader interface.
func (r *StringReader) Read(p []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	bytesRemaining := []byte(r.s[r.i:])

	// min(len(p), len(bytesRemaining))
	numToCopy := len(p)
	if len(bytesRemaining) < numToCopy {
		numToCopy = len(bytesRemaining)
	}

	copy(p, bytesRemaining[:numToCopy])

	r.i += int64(numToCopy)

	if r.i >= int64(len(r.s)) {
		err = io.EOF
	}

	return numToCopy, err
}

func NewReader(s string) *StringReader {
	return &StringReader{s: s, i: 0}
}

func outline(r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return fmt.Errorf("parsing HTML: %v", err)
	}

	depth := 0
	startElement := func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
			depth++
		}
	}

	endElement := func(n *html.Node) {
		if n.Type == html.ElementNode {
			depth--
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
	forEachNode(doc, startElement, endElement)

	return nil
}

// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
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
