package main

import (
	"testing"
	"time"
)

func TestSortTracks_Title(t *testing.T) {
	// Prepare test data
	testTracks := []*Track{
		{"B", "Artist1", "Album1", 2000, length("3m00s")},
		{"A", "Artist2", "Album2", 2001, length("2m00s")},
		{"C", "Artist3", "Album3", 1999, length("4m00s")},
	}
	tracks = testTracks

	sortTracks("Title")

	got := []string{tracks[0].Title, tracks[1].Title, tracks[2].Title}
	want := []string{"A", "B", "C"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Title sort failed: got %v, want %v", got, want)
			break
		}
	}
}

func TestSortTracks_Artist(t *testing.T) {
	testTracks := []*Track{
		{"Song", "Charlie", "Album", 2000, length("3m00s")},
		{"Song", "Bravo", "Album", 2000, length("3m00s")},
		{"Song", "Alpha", "Album", 2000, length("3m00s")},
	}
	tracks = testTracks

	sortTracks("Artist")

	got := []string{tracks[0].Artist, tracks[1].Artist, tracks[2].Artist}
	want := []string{"Alpha", "Bravo", "Charlie"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Artist sort failed: got %v, want %v", got, want)
			break
		}
	}
}

func TestSortTracks_Year(t *testing.T) {
	testTracks := []*Track{
		{"Song", "Artist", "Album", 2002, length("3m00s")},
		{"Song", "Artist", "Album", 2000, length("3m00s")},
		{"Song", "Artist", "Album", 2001, length("3m00s")},
	}
	tracks = testTracks

	sortTracks("Year")

	got := []int{tracks[0].Year, tracks[1].Year, tracks[2].Year}
	want := []int{2000, 2001, 2002}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Year sort failed: got %v, want %v", got, want)
			break
		}
	}
}

func TestSortTracks_Length(t *testing.T) {
	testTracks := []*Track{
		{"Song", "Artist", "Album", 2000, length("5m00s")},
		{"Song", "Artist", "Album", 2000, length("2m00s")},
		{"Song", "Artist", "Album", 2000, length("3m00s")},
	}
	tracks = testTracks

	sortTracks("Length")

	got := []time.Duration{tracks[0].Length, tracks[1].Length, tracks[2].Length}
	want := []time.Duration{length("2m00s"), length("3m00s"), length("5m00s")}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Length sort failed: got %v, want %v", got, want)
			break
		}
	}
}

func TestSortTracks_MultiField(t *testing.T) {
	testTracks := []*Track{
		{"Song", "Artist", "Album", 2000, length("3m00s")},
		{"Song", "Artist", "Album", 2000, length("2m00s")},
		{"Song", "Artist", "Album", 1999, length("4m00s")},
	}
	tracks = testTracks

	sortTracks("Year", "Length")

	got := []struct {
		Year   int
		Length time.Duration
	}{
		{tracks[0].Year, tracks[0].Length},
		{tracks[1].Year, tracks[1].Length},
		{tracks[2].Year, tracks[2].Length},
	}
	want := []struct {
		Year   int
		Length time.Duration
	}{
		{1999, length("4m00s")},
		{2000, length("2m00s")},
		{2000, length("3m00s")},
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Multi-field sort failed: got %v, want %v", got, want)
			break
		}
	}
}
