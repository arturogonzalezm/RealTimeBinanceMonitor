package utils

import (
	"encoding/json"
	"math"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  Int64
		expectErr bool
	}{
		{"null value", "null", 0, false},
		{"integer value", "12345", 12345, false},
		{"floating point value", "12345.67", 12346, false},
		{"floating point value round down", "12345.4", 12345, false},
		{"floating point value round half to even", "12345.5", 12346, false},
		{"floating point value round half to even 12344.5", "12344.5", 12344, false},
		{"max int64", "9223372036854775807", Int64(math.MaxInt64), false},
		{"min int64", "-9223372036854775808", Int64(math.MinInt64), false},
		{"float64 below min int64", "-9223372036854775809.0", Int64(math.MinInt64), false},
		{"string integer value", `"12345"`, 12345, false},
		{"string floating point value", `"12345.67"`, 12346, false},
		{"invalid string value", `"invalid"`, 0, true},
		{"invalid json", `{}`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i Int64
			err := json.Unmarshal([]byte(tt.input), &i)
			if (err != nil) != tt.expectErr {
				t.Errorf("UnmarshalJSON() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if i != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, expected %v", i, tt.expected)
			}
		})
	}
}
