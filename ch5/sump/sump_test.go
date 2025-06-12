package sump

import (
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"no arguments", []int{}, 0},
		{"single positive", []int{5}, 5},
		{"single negative", []int{-7}, -7},
		{"all positives", []int{1, 2, 3, 4}, 10},
		{"all negatives", []int{-1, -2, -3, -4}, -10},
		{"mixed values", []int{3, -1, 2, -4}, 0},
		{"with zero", []int{0, 1, 2}, 3},
		{"large numbers", []int{1000000, 2000000, -500000}, 2500000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sum(tt.input...)
			if got != tt.expected {
				t.Errorf("sum(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMaxFlexible(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"no arguments", []int{}, math.MinInt},
		{"single positive", []int{7}, 7},
		{"single negative", []int{-3}, -3},
		{"all positives", []int{1, 2, 3, 4, 5}, 5},
		{"all negatives", []int{-10, -20, -3, -4}, -3},
		{"mixed values", []int{-1, 0, 5, -10, 3}, 5},
		{"duplicates", []int{2, 2, 2}, 2},
		{"max at start", []int{9, 1, 2, 3}, 9},
		{"max at end", []int{1, 2, 3, 10}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := max_flexible(tt.input...)
			if got != tt.expected {
				t.Errorf("max_flexible(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMinFlexible(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"no arguments", []int{}, math.MaxInt},
		{"single positive", []int{7}, 7},
		{"single negative", []int{-3}, -3},
		{"all positives", []int{1, 2, 3, 4, 5}, 1},
		{"all negatives", []int{-10, -20, -3, -4}, -20},
		{"mixed values", []int{-1, 0, 5, -10, 3}, -10},
		{"duplicates", []int{2, 2, 2}, 2},
		{"min at start", []int{-9, 1, 2, 3}, -9},
		{"min at end", []int{1, 2, 3, -10}, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := min_flexible(tt.input...)
			if got != tt.expected {
				t.Errorf("min_flexible(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMaxAtLeastOne(t *testing.T) {
	tests := []struct {
		name     string
		first    int
		input    []int
		expected int
	}{
		{"single positive", 7, nil, 7},
		{"single negative", -3, nil, -3},
		{"all positives", 1, []int{2, 3, 4, 5}, 5},
		{"all negatives", -10, []int{-20, -3, -4}, -3},
		{"mixed values", -1, []int{0, 5, -10, 3}, 5},
		{"duplicates", 2, []int{2, 2}, 2},
		{"max at start", 9, []int{1, 2, 3}, 9},
		{"max at end", 1, []int{2, 3, 10}, 10},
		{"max in middle", 1, []int{10, 3, 2}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := max_atLeastOne(tt.first, tt.input...)
			if got != tt.expected {
				t.Errorf("max_atLeastOne(%v, %v) = %v; want %v", tt.first, tt.input, got, tt.expected)
			}
		})
	}
}

func TestMinAtLeastOne(t *testing.T) {
	tests := []struct {
		name     string
		first    int
		input    []int
		expected int
	}{
		{"single positive", 7, nil, 7},
		{"single negative", -3, nil, -3},
		{"all positives", 5, []int{6, 7, 8, 9}, 5},
		{"all negatives", -1, []int{-2, -3, -4}, -4},
		{"mixed values", 3, []int{0, -5, 10, 2}, -5},
		{"duplicates", 2, []int{2, 2}, 2},
		{"min at start", -9, []int{1, 2, 3}, -9},
		{"min at end", 1, []int{2, 3, -10}, -10},
		{"min in middle", 10, []int{3, -7, 2}, -7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := min_atLeastOne(tt.first, tt.input...)
			if got != tt.expected {
				t.Errorf("min_atLeastOne(%v, %v) = %v; want %v", tt.first, tt.input, got, tt.expected)
			}
		})
	}
}
