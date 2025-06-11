// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// The toposort program prints the nodes of a DAG in topological order.
// Practice 5.10: Modify the topoSort function so that it does not sort the
// keys of the map before visiting them. Instead, it should visit the keys
// in the order they are encountered in the map. This will change the order
// of the output, but it will still be a valid topological sort.
// Practice 5.11: Now the linear algebra teacher has made calculus a prerequisite.
// Improve topSort so that it can detect cycles in the directed graph.
package toposortp

import (
	"fmt"
)

// TopoSort performs a topological sort on the graph represented by m.
// It now returns an error if a cycle is detected.
func TopoSort(m map[string][]string) ([]string, error) {
	var order []string                // final topological sort result
	seen := make(map[string]bool)     // tracks nodes that have been fully visited (black)
	visiting := make(map[string]bool) // tracks nodes currently being visited in the recursion path (gray)

	var visitAll func(item string) error

	visitAll = func(item string) error {
		// 1. Check if the item is already being visited in the current recursion path (gray)
		if visiting[item] {
			return fmt.Errorf("cycle detected: %s is part of a cycle", item)
		}
		// 2. Check if the item has already been fully visited (black)
		if seen[item] {
			return nil
		}

		// Mark as currently visiting (gray)
		visiting[item] = true

		// Recursively visit all prerequisites of the current item
		// m[item] is a slice containing all prerequisites of the current item
		for _, prereq := range m[item] {
			if err := visitAll(prereq); err != nil {
				return err // If a cycle is found in recursion, return the error immediately
			}
		}

		// Mark as fully visited (black), and remove from currently visiting (no longer gray)
		visiting[item] = false
		seen[item] = true

		// Add the current item to the order, ensuring all its prerequisites are processed
		order = append(order, item)
		return nil
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	for _, key := range keys {
		if err := visitAll(key); err != nil {
			// Cycle detected, return error; order may be incomplete at this point
			return nil, err
		}
	}

	return order, nil
}
