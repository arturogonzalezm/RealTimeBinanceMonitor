package flexint

import (
	"encoding/json"
	"strconv"
)

// Int64 is a type that can unmarshal from string, int64, and float64 JSON values
type Int64 int64

func (fi *Int64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*fi = Int64(i)
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*fi = Int64(f)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Try parsing as int64 first
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		*fi = Int64(i)
		return nil
	}

	// If that fails, try parsing as float64 and convert to int64
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		*fi = Int64(f)
		return nil
	}

	return json.Unmarshal(data, (*int64)(fi))
}
