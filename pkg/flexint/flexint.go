package flexint

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
)

type Int64 int64

func (i *Int64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*i = 0
		return nil
	}

	// Try to unmarshal as an int64 directly
	var intVal int64
	if err := json.Unmarshal(data, &intVal); err == nil {
		*i = Int64(intVal)
		return nil
	}

	// Try to unmarshal as a float64 and then convert
	var floatVal float64
	if err := json.Unmarshal(data, &floatVal); err == nil {
		roundedVal := roundHalfToEven(floatVal)
		if roundedVal > float64(math.MaxInt64) {
			*i = Int64(math.MaxInt64)
		} else if roundedVal < float64(math.MinInt64) {
			*i = Int64(math.MinInt64)
		} else {
			*i = Int64(roundedVal)
		}
		return nil
	}

	// Try to unmarshal as a string and then convert
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			*i = Int64(intVal)
			return nil
		} else if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
			roundedVal := roundHalfToEven(floatVal)
			if roundedVal > float64(math.MaxInt64) {
				*i = Int64(math.MaxInt64)
			} else if roundedVal < float64(math.MinInt64) {
				*i = Int64(math.MinInt64)
			} else {
				*i = Int64(roundedVal)
			}
			return nil
		}
	}

	return errors.New("invalid value for Int64")
}

func roundHalfToEven(val float64) float64 {
	if val == float64(int64(val)) {
		return val
	}
	floor := math.Floor(val)
	if val-floor > 0.5 || (val-floor == 0.5 && int64(floor)%2 != 0) {
		return floor + 1
	}
	return floor
}
