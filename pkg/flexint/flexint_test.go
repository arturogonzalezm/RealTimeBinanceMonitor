package flexint

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Int64
		wantErr  bool
	}{
		{"Null value", "null", 0, false},
		{"Integer", "42", 42, false},
		{"Negative integer", "-42", -42, false},
		{"Float", "42.5", 42, false},
		{"Negative float", "-42.5", -42, false},
		{"String integer", `"42"`, 42, false},
		{"String float", `"42.5"`, 42, false},
		{"String negative integer", `"-42"`, -42, false},
		{"String negative float", `"-42.5"`, -42, false},
		{"Large integer", "9223372036854775807", 9223372036854775807, false},            // max int64
		{"Large negative integer", "-9223372036854775808", -9223372036854775808, false}, // min int64
		{"Invalid string", `"not a number"`, 0, true},
		{"Boolean true", "true", 0, true},
		{"Boolean false", "false", 0, true},
		{"Empty object", "{}", 0, true},
		{"Empty array", "[]", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi Int64
			err := json.Unmarshal([]byte(tt.input), &fi)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, fi)
			}
		})
	}
}

func TestInt64_UnmarshalJSON_LargeFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Int64
	}{
		{"Large float within int64", "9223372036854775800.0", 9223372036854775807},            // Rounds up to max int64
		{"Large negative float within int64", "-9223372036854775800.0", -9223372036854775808}, // Rounds down to min int64
		{"Float slightly above max int64", "9223372036854775808.0", 9223372036854775807},      // Max int64
		{"Float slightly below min int64", "-9223372036854775809.0", -9223372036854775808},    // Min int64
		{"Float at max int64", "9223372036854775807.0", 9223372036854775807},
		{"Float at min int64", "-9223372036854775808.0", -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi Int64
			err := json.Unmarshal([]byte(tt.input), &fi)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, fi, "For input: %s", tt.input)

			// Debugging information
			if tt.expected != fi {
				t.Logf("Input: %s, Expected: %d, Got: %d", tt.input, tt.expected, fi)
			}
		})
	}
}

func TestInt64_UnmarshalJSON_FloatPrecision(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Int64
	}{
		{"Float at max int64", "9223372036854775807.0", 9223372036854775807},
		{"Float slightly above max int64", "9223372036854775808.0", 9223372036854775807},
		{"Float at min int64", "-9223372036854775808.0", -9223372036854775808},
		{"Float slightly below min int64", "-9223372036854775809.0", -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi Int64
			err := json.Unmarshal([]byte(tt.input), &fi)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, fi, "For input: %s", tt.input)

			// Debugging information
			if tt.expected != fi {
				t.Logf("Input: %s, Expected: %d, Got: %d", tt.input, tt.expected, fi)
			}
		})
	}
}

func TestInt64_UnmarshalJSON_FloatRounding(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Int64
	}{
		{"Round down", "1.4", 1},
		{"Round up", "1.6", 2},
		{"Round half even (down)", "2.5", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi Int64
			err := json.Unmarshal([]byte(tt.input), &fi)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, fi, "For input: %s", tt.input)

			// Debugging information
			if tt.expected != fi {
				t.Logf("Input: %s, Expected: %d, Got: %d", tt.input, tt.expected, fi)
			}
		})
	}
}

func TestInt64_UnmarshalJSON_ScientificNotation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Int64
	}{
		{"Positive exponent", "1.23e2", 123},
		{"Negative exponent", "1.23e-2", 0},
		{"Large positive exponent", "1.23e10", 12300000000},
		{"Large negative exponent", "1.23e-10", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi Int64
			err := json.Unmarshal([]byte(tt.input), &fi)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, fi)

			// Debugging information
			if tt.expected != fi {
				t.Logf("Input: %s, Expected: %d, Got: %d", tt.input, tt.expected, fi)
			}
		})
	}
}
