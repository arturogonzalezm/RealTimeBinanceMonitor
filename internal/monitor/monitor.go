package monitor

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/processor"
	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/websocket"
)

// MonitorSymbol starts monitoring for a specific symbol
func MonitorSymbol(symbol string, db *sql.DB, stop chan struct{}) {
	log.Printf("Starting monitoring for symbol: %s", symbol)
	uri := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@ticker", symbol)

	client := websocket.NewClient()
	if err := client.Connect(uri); err != nil {
		log.Fatalf("WebSocket connection error for symbol %s: %v", symbol, err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing WebSocket client for symbol %s: %v", symbol, err)
		}
	}()

	// Create PGWriter
	pgWriter, err := processor.NewPGWriter(db)
	if err != nil {
		log.Fatalf("Error creating PostgreSQL writer for symbol %s: %v", symbol, err)
	}
	defer func() {
		if err := pgWriter.Close(); err != nil {
			log.Printf("Error closing PostgreSQL writer for symbol %s: %v", symbol, err)
		}
	}()

	client.AddProcessor(pgWriter)

	go client.Listen(stop)

	log.Printf("WebSocket connection opened for %s", symbol)

	// Periodic summary
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			log.Printf("Stopping monitoring for symbol: %s", symbol)
			return
		case <-ticker.C:
			// Print a summary every 5 seconds
			log.Printf("Symbol: %s - Last 5 seconds: Processed %d messages", symbol, pgWriter.GetProcessedCount())
			log.Printf("Symbol: %s - Current buffer size: %d", symbol, pgWriter.GetBufferSize())
		}
	}
}
