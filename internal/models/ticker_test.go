package models

import (
	"testing"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/pkg/flexint"
	"github.com/stretchr/testify/assert"
)

func TestFormatTickerData(t *testing.T) {
	now := time.Now().UnixMilli()
	// Create a mock TickerData
	mockTickerData := TickerData{
		EventType:    "24hrTicker",
		EventTime:    flexint.Int64(now), // Use current time
		Symbol:       "BTCUSDT",
		LastPrice:    "35000.00",
		PriceChange:  "1000.00",
		HighPrice:    "36000.00",
		LowPrice:     "34000.00",
		Volume:       "1000.5",
		QuoteVolume:  "35000000.00",
		OpenTime:     flexint.Int64(now - 86400000), // 24 hours ago
		CloseTime:    flexint.Int64(now),
		FirstTradeID: 100,
		LastTradeID:  200,
		TradeCount:   100,
	}

	// Call the function we're testing
	result := FormatTickerData(mockTickerData)

	// Assert the results
	assert.Equal(t, now, result.EventTime)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, 35000.00, result.LastPrice)
	assert.Equal(t, 1000.00, result.PriceChange)
	assert.Equal(t, 36000.00, result.HighPrice)
	assert.Equal(t, 34000.00, result.LowPrice)
	assert.Equal(t, 1000.5, result.Volume)
	assert.Equal(t, 35000000.00, result.QuoteVolume)
	assert.Equal(t, now-86400000, result.OpenTime)
	assert.Equal(t, now, result.CloseTime)
	assert.Equal(t, 100, result.TradeCount)

	// Check if Latency is reasonable (should be very small in this case)
	assert.True(t, result.Latency >= 0 && result.Latency < 100, "Latency should be non-negative and less than 100ms, got %d", result.Latency)
}

func TestParseFloat(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{"123.45", 123.45},
		{"0", 0},
		{"-67.89", -67.89},
		{"abc", 0}, // Invalid input should return 0
	}

	for _, tc := range testCases {
		result := parseFloat(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}
