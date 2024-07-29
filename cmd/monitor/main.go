package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/processor"
	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/websocket"
)

func main() {
	log.Println("Starting RealTimeCryptoMonitor...")

	symbol := "btcusdt"
	uri := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@ticker", symbol)

	client := websocket.NewClient()
	if err := client.Connect(uri); err != nil {
		log.Fatal("WebSocket connection error:", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing WebSocket client: %v", err)
		}
	}()

	// Create data directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	dataDir := filepath.Join(currentDir, "data")
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		log.Fatal("Error creating directory:", err)
	}

	// Create BufferedCSVWriter with a 5-second flush interval
	csvWriter, err := processor.NewBufferedCSVWriter(filepath.Join(dataDir, fmt.Sprintf("%s_data.csv", symbol)), 100, 5*time.Second)
	if err != nil {
		log.Fatal("Error creating CSV writer:", err)
	}
	defer func() {
		if err := csvWriter.Close(); err != nil {
			log.Printf("Error closing CSV writer: %v", err)
		}
	}()

	client.AddProcessor(csvWriter)

	stop := make(chan struct{})
	go client.Listen(stop)

	log.Printf("WebSocket connection opened for %s", symbol)

	// Set up signal catching
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Periodic summary
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigs:
			close(stop)
			log.Println("Shutting down gracefully...")
			return
		case <-ticker.C:
			// Print a summary every 5 seconds
			log.Printf("Last 5 seconds: Processed %d messages", csvWriter.GetProcessedCount())
			log.Printf("Current buffer size: %d", csvWriter.GetBufferSize())
		}
	}
}
