package toposortp

import (
	"errors"
	"strings"
	"testing"
)

// Helper to check if a slice contains all elements of another slice (order doesn't matter)
func containsAll(got, want []string) bool {
	gotMap := make(map[string]bool)
	for _, v := range got {
		gotMap[v] = true
	}
	for _, v := range want {
		if !gotMap[v] {
			return false
		}
	}
	return true
}

func TestTopoSort_Acyclic(t *testing.T) {
	acyclic := map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"d"},
		"d": {},
	}
	order, err := TopoSort(acyclic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Valid topological sorts: d before b and c, b and c before a
	pos := make(map[string]int)
	for i, v := range order {
		pos[v] = i
	}
	if !(pos["d"] < pos["b"] && pos["d"] < pos["c"] && pos["b"] < pos["a"] && pos["c"] < pos["a"]) {
		t.Errorf("invalid topological order: %v", order)
	}
}

func TestTopoSort_Cycle(t *testing.T) {
	cyclic := map[string][]string{
		"a": {"b"},
		"b": {"c"},
		"c": {"a"},
	}
	order, err := TopoSort(cyclic)
	if err == nil {
		t.Fatalf("expected error for cycle, got nil")
	}
	if order != nil {
		t.Errorf("expected nil order on cycle, got %v", order)
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("error does not mention cycle: %v", err)
	}
}

func TestTopoSort_SingleNode(t *testing.T) {
	single := map[string][]string{
		"a": {},
	}
	order, err := TopoSort(single)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 1 || order[0] != "a" {
		t.Errorf("unexpected order: %v", order)
	}
}

func TestTopoSort_Empty(t *testing.T) {
	empty := map[string][]string{}
	order, err := TopoSort(empty)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 0 {
		t.Errorf("expected empty order, got %v", order)
	}
}

func TestTopoSort_MultipleRoots(t *testing.T) {
	graph := map[string][]string{
		"a": {},
		"b": {},
		"c": {"a", "b"},
	}
	order, err := TopoSort(graph)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := make(map[string]int)
	for i, v := range order {
		pos[v] = i
	}
	if !(pos["a"] < pos["c"] && pos["b"] < pos["c"]) {
		t.Errorf("invalid topological order: %v", order)
	}
}

func TestTopoSort_PrereqsGraph(t *testing.T) {
	// Remove the cycle for this test
	prereqsNoCycle := map[string][]string{
		"algorithms":            {"data structures"},
		"calculus":              {"linear algebra"},
		"linear algebra":        {},
		"compilers":             {"data structures", "formal languages", "computer organization"},
		"data structures":       {"discrete math"},
		"databases":             {"data structures"},
		"discrete math":         {"intro to programming"},
		"formal languages":      {"discrete math"},
		"networks":              {"operating systems"},
		"operating systems":     {"data structures", "computer organization"},
		"programming languages": {"data structures", "computer organization"},
		"intro to programming":  {},
		"computer organization": {},
	}
	order, err := TopoSort(prereqsNoCycle)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Check that all courses are present
	if !containsAll(order, []string{
		"algorithms", "calculus", "linear algebra", "compilers", "data structures",
		"databases", "discrete math", "formal languages", "networks", "operating systems",
		"programming languages", "intro to programming", "computer organization",
	}) {
		t.Errorf("missing courses in order: %v", order)
	}
	// Spot check: "intro to programming" must come before "discrete math"
	pos := make(map[string]int)
	for i, v := range order {
		pos[v] = i
	}
	if !(pos["intro to programming"] < pos["discrete math"]) {
		t.Errorf("intro to programming must come before discrete math: %v", order)
	}
}

func TestTopoSort_PrereqsGraphCycle(t *testing.T) {
	// Remove the cycle for this test
	prereqsCycle := map[string][]string{
		"algorithms":            {"data structures"},
		"calculus":              {"linear algebra"},
		"linear algebra":        {"calculus"},
		"compilers":             {"data structures", "formal languages", "computer organization"},
		"data structures":       {"discrete math"},
		"databases":             {"data structures"},
		"discrete math":         {"intro to programming"},
		"formal languages":      {"discrete math"},
		"networks":              {"operating systems"},
		"operating systems":     {"data structures", "computer organization"},
		"programming languages": {"data structures", "computer organization"},
		"intro to programming":  {},
		"computer organization": {},
	}
	_, err := TopoSort(prereqsCycle)
	if err == nil {
		t.Fatal("expected error for self-cycle, got nil")
	}
	if !errors.Is(err, err) || !strings.Contains(err.Error(), "cycle") {
		t.Errorf("error does not mention cycle: %v", err)
	}
}

func TestTopoSort_SelfCycle(t *testing.T) {
	selfCycle := map[string][]string{
		"a": {"a"},
	}
	_, err := TopoSort(selfCycle)
	if err == nil {
		t.Fatal("expected error for self-cycle, got nil")
	}
	if !errors.Is(err, err) || !strings.Contains(err.Error(), "cycle") {
		t.Errorf("error does not mention cycle: %v", err)
	}
}
