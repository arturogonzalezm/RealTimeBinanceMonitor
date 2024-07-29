package processor

import (
	"encoding/csv"
	"os"
	"testing"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewBufferedCSVWriter(t *testing.T) {
	filename := "test.csv"
	bufferSize := 10
	flushInterval := 5 * time.Second

	writer, err := NewBufferedCSVWriter(filename, bufferSize, flushInterval)
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	// Clean up
	writer.Close()
	os.Remove(filename)
}

func TestBufferedCSVWriter_Process(t *testing.T) {
	filename := "test_process.csv"
	bufferSize := 2
	flushInterval := 5 * time.Second

	writer, err := NewBufferedCSVWriter(filename, bufferSize, flushInterval)
	assert.NoError(t, err)

	data := models.FormattedData{
		EventTime:   1625097600000,
		Symbol:      "BTCUSDT",
		LastPrice:   35000.00,
		PriceChange: 1000.00,
		HighPrice:   36000.00,
		LowPrice:    34000.00,
		Volume:      1000.5,
		QuoteVolume: 35000000.00,
		OpenTime:    1625011200000,
		CloseTime:   1625097600000,
		TradeCount:  100,
		Latency:     50,
	}

	writer.Process(data)
	assert.Equal(t, 1, writer.GetProcessedCount())
	assert.Equal(t, 1, writer.GetBufferSize())

	writer.Process(data)
	assert.Equal(t, 2, writer.GetProcessedCount())
	assert.Equal(t, 0, writer.GetBufferSize()) // Buffer should be flushed

	// Clean up
	writer.Close()
	os.Remove(filename)
}

func TestBufferedCSVWriter_Flush(t *testing.T) {
	filename := "test_flush.csv"
	bufferSize := 10
	flushInterval := 100 * time.Millisecond

	writer, err := NewBufferedCSVWriter(filename, bufferSize, flushInterval)
	assert.NoError(t, err)

	data := models.FormattedData{
		EventTime:   1625097600000,
		Symbol:      "BTCUSDT",
		LastPrice:   35000.00,
		PriceChange: 1000.00,
		HighPrice:   36000.00,
		LowPrice:    34000.00,
		Volume:      1000.5,
		QuoteVolume: 35000000.00,
		OpenTime:    1625011200000,
		CloseTime:   1625097600000,
		TradeCount:  100,
		Latency:     50,
	}

	writer.Process(data)

	// Force a flush
	err = writer.Close()
	assert.NoError(t, err)

	// Check if data was written to file
	file, err := os.Open(filename)
	assert.NoError(t, err)
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(records), "Expected 2 records (header + data row), but got %d", len(records))

	// Clean up
	os.Remove(filename)
}

func TestBufferedCSVWriter_Close(t *testing.T) {
	filename := "test_close.csv"
	bufferSize := 10
	flushInterval := 5 * time.Second

	writer, err := NewBufferedCSVWriter(filename, bufferSize, flushInterval)
	assert.NoError(t, err)

	data := models.FormattedData{
		EventTime:   1625097600000,
		Symbol:      "BTCUSDT",
		LastPrice:   35000.00,
		PriceChange: 1000.00,
		HighPrice:   36000.00,
		LowPrice:    34000.00,
		Volume:      1000.5,
		QuoteVolume: 35000000.00,
		OpenTime:    1625011200000,
		CloseTime:   1625097600000,
		TradeCount:  100,
		Latency:     50,
	}

	writer.Process(data)
	assert.Equal(t, 1, writer.GetBufferSize())

	err = writer.Close()
	assert.NoError(t, err)

	// Check if data was written to file
	file, err := os.Open(filename)
	assert.NoError(t, err)
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(records)) // Header + 1 data row

	// Clean up
	os.Remove(filename)
}
