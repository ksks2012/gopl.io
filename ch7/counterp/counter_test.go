package counter

import (
	"testing"
)

func TestByteCounter_Write(t *testing.T) {
	var c ByteCounter
	n, err := c.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 11 {
		t.Errorf("expected 11 bytes written, got %d", n)
	}
	if c != 11 {
		t.Errorf("expected counter to be 11, got %d", c)
	}

	n, err = c.Write([]byte("!"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 byte written, got %d", n)
	}
	if c != 12 {
		t.Errorf("expected counter to be 12, got %d", c)
	}
}

func TestWordCounter_Write(t *testing.T) {
	var c WordCounter
	n, err := c.Write([]byte("hello world test"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 16 {
		t.Errorf("expected 16 bytes written, got %d", n)
	}
	if c != 3 {
		t.Errorf("expected counter to be 3, got %d", c)
	}

	n, err = c.Write([]byte(" more\nwords\t"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 12 {
		t.Errorf("expected 12 bytes written, got %d", n)
	}
	if c != 5 {
		t.Errorf("expected counter to be 5, got %d", c)
	}

	var c2 WordCounter
	c2.Write([]byte(""))
	if c2 != 0 {
		t.Errorf("expected 0 words for empty string, got %d", c2)
	}
	c2.Write([]byte("   \n\t "))
	if c2 != 0 {
		t.Errorf("expected 0 words for only whitespace, got %d", c2)
	}
	c2.Write([]byte("singleword"))
	if c2 != 1 {
		t.Errorf("expected 1 word for 'singleword', got %d", c2)
	}
}

func TestLineCounter_Write(t *testing.T) {
	var c LineCounter
	n, err := c.Write([]byte("line1\nline2\nline3"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 17 {
		t.Errorf("expected 17 bytes written, got %d", n)
	}
	if c != 3 {
		t.Errorf("expected counter to be 3, got %d", c)
	}

	n, err = c.Write([]byte("\n\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 bytes written, got %d", n)
	}
	if c != 5 {
		t.Errorf("expected counter to be 5, got %d", c)
	}

	var c2 LineCounter
	c2.Write([]byte(""))
	if c2 != 0 {
		t.Errorf("expected 0 lines for empty string, got %d", c2)
	}
	c2.Write([]byte("\n"))
	if c2 != 1 {
		t.Errorf("expected 1 line for single newline, got %d", c2)
	}
	c2.Write([]byte("last line"))
	if c2 != 2 {
		t.Errorf("expected 2 lines, got %d", c2)
	}
}
