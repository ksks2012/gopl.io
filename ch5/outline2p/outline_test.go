// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"golang.org/x/net/html"
)

// helper to capture stdout
func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = stdout
	buf.ReadFrom(r)
	return buf.String()
}

func TestPrintNodePre_ElementNode(t *testing.T) {
	depth = 0
	node := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			{Key: "class", Val: "container"},
			{Key: "id", Val: "main"},
		},
	}
	got := captureOutput(func() { printNodePre(node) })
	want := `<div class="container" id="main"/>
`
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestPrintNodePre_SelfClosingElement(t *testing.T) {
	depth = 0
	node := &html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			{Key: "src", Val: "logo.png"},
			{Key: "alt", Val: "Logo"},
		},
	}
	got := captureOutput(func() { printNodePre(node) })
	want := `<img src="logo.png" alt="Logo"/>
`
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestPrintNodePre_TextNode(t *testing.T) {
	depth = 1
	node := &html.Node{
		Type: html.TextNode,
		Data: "  Hello, world!  ",
	}
	got := captureOutput(func() { printNodePre(node) })
	want := fmt.Sprintf("%*s%s\n", depth*2, "", "Hello, world!")
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestPrintNodePre_CommentNode(t *testing.T) {
	depth = 2
	node := &html.Node{
		Type: html.CommentNode,
		Data: " This is a comment ",
	}
	got := captureOutput(func() { printNodePre(node) })
	want := fmt.Sprintf("%*s%s\n", depth*2, "", " This is a comment ")
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestPrintNodePre_DocumentNode(t *testing.T) {
	node := &html.Node{
		Type: html.DocumentNode,
	}
	got := captureOutput(func() { printNodePre(node) })
	want := "<!DOCTYPE html>\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestPrintNodePre_ScriptAndStyleNotSelfClosing(t *testing.T) {
	for _, tag := range []string{"script", "style"} {
		depth = 0
		node := &html.Node{
			Type: html.ElementNode,
			Data: tag,
			Attr: []html.Attribute{{Key: "type", Val: "text/javascript"}},
		}
		got := captureOutput(func() { printNodePre(node) })
		want := fmt.Sprintf("<%s type=\"text/javascript\">\n", tag)
		if got != want {
			t.Errorf("tag %q: got:\n%q\nwant:\n%q", tag, got, want)
		}
	}
}

func TestPrintNodePre_Indented(t *testing.T) {
	depth = 2
	node := &html.Node{
		Type: html.ElementNode,
		Data: "p",
	}
	got := captureOutput(func() { printNodePre(node) })
	want := fmt.Sprintf("%*s<%s/>\n", depth*2, "", "p")
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestPrintNodePre_EmptyTextNode(t *testing.T) {
	depth = 0
	node := &html.Node{
		Type: html.TextNode,
		Data: "   ",
	}
	got := captureOutput(func() { printNodePre(node) })
	if got != "" {
		t.Errorf("expected no output for empty text node, got: %q", got)
	}
}

func TestPrintNodePre_AttributesEscaping(t *testing.T) {
	depth = 0
	node := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{Key: "href", Val: "https://example.com/?q=go&lang=en"},
			{Key: "title", Val: `A "quote"`},
		},
	}
	got := captureOutput(func() { printNodePre(node) })
	// The function does not escape quotes, so expect raw output
	want := `<a href="https://example.com/?q=go&lang=en" title="A "quote""/>
`
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestPrintNodePre_NestedElements(t *testing.T) {
	// This test checks that depth is handled correctly for nested elements.
	dep := 0
	parent := &html.Node{
		Type: html.ElementNode,
		Data: "ul",
	}
	child := &html.Node{
		Type: html.ElementNode,
		Data: "li",
	}
	parent.FirstChild = child
	gotParent := captureOutput(func() { printNodePre(parent) })
	wantParent := fmt.Sprintf("%*s<ul>\n", dep*2, "")
	if gotParent != wantParent {
		t.Errorf("parent got:\n%q\nwant:\n%q", gotParent, wantParent)
	}
	gotChild := captureOutput(func() { printNodePre(child) })
	wantChild := fmt.Sprintf("%*s<li/>\n", (dep+1)*2, "")
	if gotChild != wantChild {
		t.Errorf("child got:\n%q\nwant:\n%q", gotChild, wantChild)
	}
}

func TestPrintNodePre_SelfClosingWithNoAttrs(t *testing.T) {
	depth = 0
	node := &html.Node{
		Type: html.ElementNode,
		Data: "br",
	}
	got := captureOutput(func() { printNodePre(node) })
	want := "<br/>\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}
