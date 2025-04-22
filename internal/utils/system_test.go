package utils

import (
	"fmt"
	"math"
	"testing"
)

func TestParseSizeString(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
		wantErr  bool
	}{
		// Basic units
		{"1024", 1024, false},
		{"1024b", 1024, false},
		{"1k", 1024, false},
		{"1kb", 1024, false},
		{"1m", 1024 * 1024, false},
		{"1mb", 1024 * 1024, false},
		{"1g", 1024 * 1024 * 1024, false},
		{"1gb", 1024 * 1024 * 1024, false},
		{"1t", 1024 * 1024 * 1024 * 1024, false},
		{"1tb", 1024 * 1024 * 1024 * 1024, false},

		// Case insensitive
		{"2KB", 2 * 1024, false},
		{"3MB", 3 * 1024 * 1024, false},
		{"4GB", 4 * 1024 * 1024 * 1024, false},
		{"5TB", 5 * 1024 * 1024 * 1024 * 1024, false},

		// Float values
		{"1.5kb", int64(1.5 * 1024), false},
		{"0.5MB", int64(0.5 * 1024 * 1024), false},
		{"100.2b", 100, false},

		// Zero
		{"0", 0, false},
		{"0kb", 0, false},

		// Whitespace
		{" 512kb ", 512 * 1024, false},
		{"1 mb", 1024 * 1024, false},

		// Max Int
		{fmt.Sprintf("%d", math.MaxInt64), math.MaxInt64, false},

		// Errors
		{"", 0, true},          // Empty string
		{"kb", 0, true},        // No number
		{"10xx", 0, true},      // Invalid unit
		{"-1kb", 0, true},      // Negative number
		{"1.2.3kb", 0, true},   // Invalid number format
		{"invalid", 0, true},   // Completely invalid
		{"1eb", 0, true},       // Unsupported large unit
		{"1 kb test", 0, true}, // Extra text after unit
		{"1.5.mb", 0, true},    // Multiple decimals
		{"1,024kb", 0, true},   // Comma in number
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := ParseSizeString(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseSizeString(%q) expected an error, but got nil", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseSizeString(%q) unexpected error: %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("ParseSizeString(%q) = %d, want %d", tc.input, result, tc.expected)
				}
			}
		})
	}
}
