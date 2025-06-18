// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 214.
//!+

// Xmlselect prints the text of selected elements of an XML document.
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

type selector struct {
	Name    string   // Element name (e.g., "div", "p")
	ID      string   // ID selector value (e.g., "page")
	Classes []string // Class selector values (e.g., {"wide", "left"})
}

// matches method determines whether an XML element's attributes match the selector's rules
func (s selector) matches(name string, attrs []xml.Attr) bool {
	// If the selector specifies a name, it must match
	if s.Name != "" && s.Name != name {
		return false
	}

	// If the selector specifies an ID, it must match
	if s.ID != "" {
		foundID := false
		for _, attr := range attrs {
			if attr.Name.Local == "id" && attr.Value == s.ID {
				foundID = true
				break
			}
		}
		if !foundID {
			return false
		}
	}

	// If the selector specifies classes, all must be present in the element's class attribute
	if len(s.Classes) > 0 {
		elementClasses := make(map[string]bool)
		for _, attr := range attrs {
			if attr.Name.Local == "class" {
				// The class attribute may contain multiple classes, separated by spaces
				for _, cls := range strings.Fields(attr.Value) {
					elementClasses[cls] = true
				}
				break
			}
		}

		for _, requiredClass := range s.Classes {
			if !elementClasses[requiredClass] {
				return false // Missing required class
			}
		}
	}

	// All conditions matched
	return true
}

// matchesAnySelector checks if the given element matches any of the provided selectors
func matchesAnySelector(selectors []selector, name string, attrs []xml.Attr) bool {
	for _, sel := range selectors {
		if sel.matches(name, attrs) {
			return true
		}
	}
	return false
}

// Helper function: converts attributes to a printable string
func getAttrsString(attrs []xml.Attr) string {
	var sb strings.Builder
	for i, attr := range attrs {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%s=\"%s\"", attr.Name.Local, attr.Value))
	}
	return sb.String()
}

// containsAll reports whether x contains the elements of y, in order.
func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

//!-

func main() {
	// 1. Parse command-line arguments into our selector type
	var selectors []selector
	for _, arg := range os.Args[1:] { // Skip program name
		s := selector{}
		if strings.HasPrefix(arg, "#") {
			// ID selector
			s.ID = arg[1:]
		} else if strings.HasPrefix(arg, ".") {
			// Class selector
			s.Classes = []string{arg[1:]} // Assume one class per argument
		} else {
			// Element name
			s.Name = arg
		}
		selectors = append(selectors, s)
	}

	// If no selectors are provided, default to usage message (can adjust as needed)
	if len(selectors) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: xmlselect <selector> ...")
		fmt.Fprintln(os.Stderr, "Selectors can be element names, #id, or .class")
		os.Exit(1)
	}

	decoder := xml.NewDecoder(os.Stdin)
	var stack []string // Stack of element names
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break // End of input
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}

		switch tok := token.(type) {
		case xml.StartElement:
			stack = append(stack, tok.Name.Local) // Push onto stack
			// 2. Update matching logic: check if current element matches any selector
			if matchesAnySelector(selectors, tok.Name.Local, tok.Attr) {
				fmt.Printf("%s %s\n", strings.Join(stack, " "), getAttrsString(tok.Attr))
			}
		case xml.EndElement:
			stack = stack[:len(stack)-1] // Pop from stack
		}
	}
}
