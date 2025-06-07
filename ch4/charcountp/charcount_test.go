package charcountp_test

import (
	"bytes"
	"os"
	"testing"

	"gopl.io/ch4/charcountp"
)

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

func TestCharCount(t *testing.T) {
	tests := []struct {
		name     string
		input    map[rune]int
		expected string
	}{
		{
			name:     "empty map",
			input:    map[rune]int{},
			expected: "Letters: 0\nDigits: 0\nNumbers: 0\n",
		},
		{
			name: "letters only",
			input: map[rune]int{
				'a': 2,
				'B': 3,
			},
			expected: "Letters: 5\nDigits: 0\nNumbers: 0\n",
		},
		{
			name: "digits only",
			input: map[rune]int{
				'1': 4,
				'9': 1,
			},
			expected: "Letters: 0\nDigits: 5\nNumbers: 5\n",
		},
		{
			name: "letters, digits, and numbers",
			input: map[rune]int{
				'a': 1,
				'2': 2,
				'Ⅷ': 3, // Roman numeral eight, unicode.IsNumber true, unicode.IsDigit false
				'β': 1, // Greek letter
				'٣': 2, // Arabic-Indic digit three, unicode.IsDigit true, unicode.IsNumber true
				'@': 1, // symbol, not letter/digit/number
			},
			expected: "Letters: 2\nDigits: 4\nNumbers: 7\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				charcountp.CharCount(tt.input)
			})
			if output != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, output)
			}
		})
	}
}

func TestWordFreq(t *testing.T) {
	// Save original os.Stdin
	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()

	tests := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{
			name:     "empty input",
			input:    "",
			expected: map[string]int{},
		},
		{
			name:  "single word",
			input: "hello",
			expected: map[string]int{
				"hello": 1,
			},
		},
		{
			name:  "multiple words",
			input: "foo bar foo",
			expected: map[string]int{
				"foo": 2,
				"bar": 1,
			},
		},
		{
			name:  "words with punctuation",
			input: "hi! hi, hi.",
			expected: map[string]int{
				"hi!": 1,
				"hi,": 1,
				"hi.": 1,
			},
		},
		{
			name:  "unicode words",
			input: "你好 你好 hello",
			expected: map[string]int{
				"你好":    2,
				"hello": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := charcountp.WordFreq(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("expected map of length %d, got %d", len(tt.expected), len(got))
			}
			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("for word %q: expected %d, got %d", k, v, got[k])
				}
			}
		})
	}
}
