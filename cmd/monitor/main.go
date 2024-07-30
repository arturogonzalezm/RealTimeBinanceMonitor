package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/processor"
	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/websocket"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting RealTimeCryptoMonitor...")

	symbols := []string{"btcusdt", "ethusdt", "ltcusdt"} // Add more symbols as needed

	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	fmt.Printf("DB_HOST: %s, DB_USER: %s, DB_PASSWORD: %s, DB_NAME: %s\n", dbHost, dbUser, dbPassword, dbName)

	// Database connection setup
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	var db *sql.DB
	var err error

	// Retry connecting to the database until successful
	for {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error connecting to the database: %v", err)
		} else if err = db.Ping(); err == nil {
			break
		}
		log.Println("Waiting for the database to be ready...")
		time.Sleep(2 * time.Second)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing the database: %v", err)
		}
	}()

	// Channel to handle graceful shutdown
	stop := make(chan struct{})
	// Channel to listen for OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup to manage goroutines
	var wg sync.WaitGroup

	for _, symbol := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			monitorSymbol(symbol, db, stop)
		}(symbol)
	}

	// Wait for an interrupt signal
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s. Shutting down gracefully...", sig)
		close(stop)
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("All symbol monitoring stopped. Exiting program.")
}

func monitorSymbol(symbol string, db *sql.DB, stop chan struct{}) {
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
