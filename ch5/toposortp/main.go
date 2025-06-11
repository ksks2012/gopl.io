// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

// The toposort program prints the nodes of a DAG in topological order.
// Practice 5.10: Modify the topoSort function so that it does not sort the
// keys of the map before visiting them. Instead, it should visit the keys
// in the order they are encountered in the map. This will change the order
// of the output, but it will still be a valid topological sort.
package main

import (
	"fmt"
)

// !+table
// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

//!-table

// !+main
func main() {
	fmt.Println("--- Topological Sort Results (Order may vary) ---")
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

// topoSort performs a topological sort on the graph represented by m.
// The order of the output may vary on different runs due to map iteration order.
func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)

	// This is a recursive DFS (Depth-First Search) traversal function.
	// It visits each node in items and all its unvisited prerequisites.
	var visitAll func(items []string)

	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	for course := range m {
		visitAll([]string{course})
	}

	return order
}

//!-main
