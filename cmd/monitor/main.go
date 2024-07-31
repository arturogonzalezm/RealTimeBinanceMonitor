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

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/monitor"
)

type Config struct {
	DB struct {
		Host     string
		User     string
		Password string
		Name     string
	}
	Symbols []string
}

func main() {
	log.Println("Starting RealTimeCryptoMonitor...")

	// Load configuration
	var config Config
	if err := loadConfig(&config); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Print loaded configuration
	fmt.Printf("DB_HOST: %s, DB_USER: %s, DB_PASSWORD: %s, DB_NAME: %s\n", config.DB.Host, config.DB.User, config.DB.Password, config.DB.Name)

	// Database connection setup
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name)
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

	for _, symbol := range config.Symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			monitor.MonitorSymbol(symbol, db, stop)
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

func loadConfig(config *Config) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&config)
}
