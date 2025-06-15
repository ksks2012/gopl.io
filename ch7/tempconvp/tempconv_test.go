package tempconv

import (
	"testing"
)

func TestCelsiusFlagSet(t *testing.T) {
	tests := []struct {
		input   string
		want    Celsius
		wantErr bool
	}{
		{"100C", 100, false},
		{"0C", 0, false},
		{"-40C", -40, false},
		{"212F", 100, false},
		{"32F", 0, false},
		{"-40F", -40, false},
		{"273.15K", 0, false},
		{"373.15K", 100, false},
		{"0K", -273.15, false},
		{"100", 0, true},         // missing unit
		{"abcC", 0, true},        // invalid number
		{"100X", 0, true},        // invalid unit
		{"", 0, true},            // empty string
		{"100°C", 100, false},    // unicode degree
		{"212°F", 100, false},    // unicode degree
		{"373.15°K", 100, false}, // unicode degree
	}

	for _, tt := range tests {
		var f celsiusFlag
		err := f.Set(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("Set(%q) expected error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("Set(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if diff := float64(f.Celsius - tt.want); diff < -1e-6 || diff > 1e-6 {
			t.Errorf("Set(%q) = %v, want %v", tt.input, f.Celsius, tt.want)
		}
	}
}
