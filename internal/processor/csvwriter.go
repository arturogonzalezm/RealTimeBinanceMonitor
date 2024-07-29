package processor

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/models"
)

// BufferedCSVWriter implements DataProcessor interface
type BufferedCSVWriter struct {
	buffer         []models.FormattedData
	bufferSize     int
	csvFile        *os.File
	csvWriter      *csv.Writer
	mutex          sync.Mutex
	processedCount int
	lastFlushTime  time.Time
	flushInterval  time.Duration
}

// NewBufferedCSVWriter creates a new BufferedCSVWriter
func NewBufferedCSVWriter(filename string, bufferSize int, flushInterval time.Duration) (*BufferedCSVWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{
		"Event Time", "Symbol", "Last Price", "Price Change", "High Price", "Low Price",
		"Volume", "Quote Volume", "Open Time", "Close Time", "Trade Count", "Latency",
	}); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("error writing CSV header: %v; error closing file: %v", err, closeErr)
		}
		return nil, fmt.Errorf("error writing CSV header: %v", err)
	}

	return &BufferedCSVWriter{
		buffer:        make([]models.FormattedData, 0, bufferSize),
		bufferSize:    bufferSize,
		csvFile:       file,
		csvWriter:     writer,
		lastFlushTime: time.Now(),
		flushInterval: flushInterval,
	}, nil
}

// Process implements the DataProcessor interface
func (w *BufferedCSVWriter) Process(data models.FormattedData) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.buffer = append(w.buffer, data)
	w.processedCount++

	if len(w.buffer) >= w.bufferSize || time.Since(w.lastFlushTime) >= w.flushInterval {
		if err := w.flush(); err != nil {
			fmt.Printf("Error flushing buffer: %v\n", err)
		}
		w.lastFlushTime = time.Now()
	}
}

// GetProcessedCount returns the number of processed messages
func (w *BufferedCSVWriter) GetProcessedCount() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.processedCount
}

// GetBufferSize returns the current size of the buffer
func (w *BufferedCSVWriter) GetBufferSize() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return len(w.buffer)
}

// flush writes the buffered data to the CSV file
func (w *BufferedCSVWriter) flush() error {
	for _, data := range w.buffer {
		if err := w.csvWriter.Write([]string{
			strconv.FormatInt(data.EventTime, 10),
			data.Symbol,
			strconv.FormatFloat(data.LastPrice, 'f', -1, 64),
			strconv.FormatFloat(data.PriceChange, 'f', -1, 64),
			strconv.FormatFloat(data.HighPrice, 'f', -1, 64),
			strconv.FormatFloat(data.LowPrice, 'f', -1, 64),
			strconv.FormatFloat(data.Volume, 'f', -1, 64),
			strconv.FormatFloat(data.QuoteVolume, 'f', -1, 64),
			strconv.FormatInt(data.OpenTime, 10),
			strconv.FormatInt(data.CloseTime, 10),
			strconv.Itoa(data.TradeCount),
			strconv.FormatInt(data.Latency, 10),
		}); err != nil {
			return fmt.Errorf("error writing CSV record: %v", err)
		}
	}
	w.csvWriter.Flush()
	if err := w.csvWriter.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %v", err)
	}
	w.buffer = w.buffer[:0]
	return nil
}

// Close flushes any remaining data and closes the CSV file
func (w *BufferedCSVWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err := w.flush(); err != nil {
		return fmt.Errorf("error flushing buffer on close: %v", err)
	}
	if err := w.csvFile.Close(); err != nil {
		return fmt.Errorf("error closing CSV file: %v", err)
	}
	return nil
}
