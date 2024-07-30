package processor

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/models"
)

// PGWriter implements DataProcessor interface for PostgreSQL
type PGWriter struct {
	db             *sql.DB
	mutex          sync.Mutex
	processedCount int
}

// NewPGWriter creates a new PGWriter
func NewPGWriter(db *sql.DB) (*PGWriter, error) {
	writer := &PGWriter{
		db: db,
	}

	return writer, nil
}

// Process implements the DataProcessor interface
func (w *PGWriter) Process(data models.FormattedData) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	_, err := w.db.Exec(`INSERT INTO ticker_data (
        event_time, symbol, last_price, price_change, high_price, low_price, volume, quote_volume, open_time, close_time, trade_count, latency
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		data.EventTime, data.Symbol, data.LastPrice, data.PriceChange, data.HighPrice, data.LowPrice,
		data.Volume, data.QuoteVolume, data.OpenTime, data.CloseTime, data.TradeCount, data.Latency,
	)
	if err != nil {
		fmt.Printf("Error inserting data: %v\n", err)
		return
	}

	w.processedCount++
}

// GetProcessedCount returns the number of processed messages
func (w *PGWriter) GetProcessedCount() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.processedCount
}

// GetBufferSize returns the current size of the buffer (always 0 for immediate insert)
func (w *PGWriter) GetBufferSize() int {
	return 0
}

// Close closes the PostgreSQL connection
func (w *PGWriter) Close() error {
	return w.db.Close()
}
