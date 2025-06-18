package main

import (
	"encoding/xml"
	"testing"
)

// TestSelectorMatches tests the selector.matches method for various selector types.
func TestSelectorMatches(t *testing.T) {
	tests := []struct {
		name     string
		selector selector
		elemName string
		attrs    []xml.Attr
		want     bool
	}{
		{
			name:     "element name match",
			selector: selector{Name: "div"},
			elemName: "div",
			attrs:    nil,
			want:     true,
		},
		{
			name:     "element name mismatch",
			selector: selector{Name: "div"},
			elemName: "span",
			attrs:    nil,
			want:     false,
		},
		{
			name:     "id match",
			selector: selector{ID: "main"},
			elemName: "div",
			attrs:    []xml.Attr{{Name: xml.Name{Local: "id"}, Value: "main"}},
			want:     true,
		},
		{
			name:     "id mismatch",
			selector: selector{ID: "main"},
			elemName: "div",
			attrs:    []xml.Attr{{Name: xml.Name{Local: "id"}, Value: "other"}},
			want:     false,
		},
		{
			name:     "class match single",
			selector: selector{Classes: []string{"foo"}},
			elemName: "div",
			attrs:    []xml.Attr{{Name: xml.Name{Local: "class"}, Value: "foo"}},
			want:     true,
		},
		{
			name:     "class match multiple",
			selector: selector{Classes: []string{"foo", "bar"}},
			elemName: "div",
			attrs:    []xml.Attr{{Name: xml.Name{Local: "class"}, Value: "foo bar"}},
			want:     true,
		},
		{
			name:     "class mismatch",
			selector: selector{Classes: []string{"foo", "baz"}},
			elemName: "div",
			attrs:    []xml.Attr{{Name: xml.Name{Local: "class"}, Value: "foo bar"}},
			want:     false,
		},
		{
			name:     "name and id match",
			selector: selector{Name: "div", ID: "main"},
			elemName: "div",
			attrs:    []xml.Attr{{Name: xml.Name{Local: "id"}, Value: "main"}},
			want:     true,
		},
		{
			name:     "name and id mismatch",
			selector: selector{Name: "div", ID: "main"},
			elemName: "span",
			attrs:    []xml.Attr{{Name: xml.Name{Local: "id"}, Value: "main"}},
			want:     false,
		},
		{
			name:     "name, id, and class match",
			selector: selector{Name: "div", ID: "main", Classes: []string{"foo"}},
			elemName: "div",
			attrs: []xml.Attr{
				{Name: xml.Name{Local: "id"}, Value: "main"},
				{Name: xml.Name{Local: "class"}, Value: "foo"},
			},
			want: true,
		},
		{
			name:     "name, id, and class mismatch",
			selector: selector{Name: "div", ID: "main", Classes: []string{"foo"}},
			elemName: "div",
			attrs: []xml.Attr{
				{Name: xml.Name{Local: "id"}, Value: "main"},
				{Name: xml.Name{Local: "class"}, Value: "bar"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		got := tt.selector.matches(tt.elemName, tt.attrs)
		if got != tt.want {
			t.Errorf("%s: selector.matches(%q, %v) = %v, want %v", tt.name, tt.elemName, tt.attrs, got, tt.want)
		}
	}
}

// TestMatchesAnySelector tests matchesAnySelector function.
func TestMatchesAnySelector(t *testing.T) {
	selectors := []selector{
		{Name: "div"},
		{ID: "main"},
		{Classes: []string{"foo"}},
	}
	tests := []struct {
		name     string
		elemName string
		attrs    []xml.Attr
		want     bool
	}{
		{"matches name", "div", nil, true},
		{"matches id", "span", []xml.Attr{{Name: xml.Name{Local: "id"}, Value: "main"}}, true},
		{"matches class", "p", []xml.Attr{{Name: xml.Name{Local: "class"}, Value: "foo"}}, true},
		{"no match", "span", nil, false},
	}
	for _, tt := range tests {
		got := matchesAnySelector(selectors, tt.elemName, tt.attrs)
		if got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

// TestGetAttrsString tests getAttrsString function.
func TestGetAttrsString(t *testing.T) {
	attrs := []xml.Attr{
		{Name: xml.Name{Local: "id"}, Value: "main"},
		{Name: xml.Name{Local: "class"}, Value: "foo bar"},
	}
	got := getAttrsString(attrs)
	want := `id="main" class="foo bar"`
	if got != want {
		t.Errorf("getAttrsString() = %q, want %q", got, want)
	}
}

// TestContainsAll tests containsAll function.
func TestContainsAll(t *testing.T) {
	tests := []struct {
		x, y []string
		want bool
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b"}, true},
		{[]string{"a", "b", "c"}, []string{"b", "c"}, true},
		{[]string{"a", "b", "c"}, []string{"a", "c"}, true},
		{[]string{"a", "b", "c"}, []string{"c", "a"}, false},
		{[]string{"a", "b", "c"}, []string{"d"}, false},
		{[]string{"a", "b", "c"}, []string{}, true},
	}
	for _, tt := range tests {
		got := containsAll(tt.x, tt.y)
		if got != tt.want {
			t.Errorf("containsAll(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
		}
	}
}
