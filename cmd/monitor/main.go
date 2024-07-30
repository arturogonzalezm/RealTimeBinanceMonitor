package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/processor"
	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/websocket"
	_ "github.com/lib/pq"
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

	// Create PGWriter
	pgWriter, err := processor.NewPGWriter(db)
	if err != nil {
		log.Fatal("Error creating PostgreSQL writer:", err)
	}
	defer func() {
		if err := pgWriter.Close(); err != nil {
			log.Printf("Error closing PostgreSQL writer: %v", err)
		}
	}()

	client.AddProcessor(pgWriter)

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
			log.Printf("Last 5 seconds: Processed %d messages", pgWriter.GetProcessedCount())
			log.Printf("Current buffer size: %d", pgWriter.GetBufferSize())
		}
	}
}
