package main

import (
	"sort"
	"testing"
)

func TestIsPalindrome_IntSlice(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected bool
	}{
		{"empty", []int{}, true},
		{"single", []int{1}, true},
		{"palindrome even", []int{1, 2, 2, 1}, true},
		{"palindrome odd", []int{1, 2, 1}, true},
		{"not palindrome", []int{1, 2, 3}, false},
		{"not palindrome 2", []int{1, 2, 3, 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(sort.IntSlice(tt.values))
			if got != tt.expected {
				t.Errorf("IsPalindrome(%v) = %v, want %v", tt.values, got, tt.expected)
			}
		})
	}
}

type stringSlice []string

func (s stringSlice) Len() int           { return len(s) }
func (s stringSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s stringSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func TestIsPalindrome_StringSlice(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected bool
	}{
		{"empty", []string{}, true},
		{"single", []string{"a"}, true},
		{"palindrome", []string{"a", "b", "a"}, true},
		{"not palindrome", []string{"a", "b", "c"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(stringSlice(tt.values))
			if got != tt.expected {
				t.Errorf("IsPalindrome(%v) = %v, want %v", tt.values, got, tt.expected)
			}
		})
	}
}

func TestIsPalindrome_byArtist(t *testing.T) {
	tracks := []*Track{
		{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
		{"Go", "Moby", "Moby", 1992, length("3m37s")},
		{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	}
	if !IsPalindrome(byArtist(tracks)) {
		t.Errorf("IsPalindrome(byArtist) = false, want true")
	}
	tracks2 := []*Track{
		{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
		{"Go", "Moby", "Moby", 1992, length("3m37s")},
		{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	}
	if IsPalindrome(byArtist(tracks2)) {
		t.Errorf("IsPalindrome(byArtist) = true, want false")
	}
}
