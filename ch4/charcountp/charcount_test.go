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
