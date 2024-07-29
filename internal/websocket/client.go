package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/models"
	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/processor"
	"github.com/gorilla/websocket"
)

// Client manages the WebSocket connection and data processing
type Client struct {
	conn       *websocket.Conn
	processors []processor.DataProcessor
	mutex      sync.RWMutex
}

// NewClient creates a new Client
func NewClient() *Client {
	return &Client{
		processors: make([]processor.DataProcessor, 0),
	}
}

// AddProcessor adds a new data processor
func (c *Client) AddProcessor(proc processor.DataProcessor) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.processors = append(c.processors, proc)
}

// Connect establishes a WebSocket connection
func (c *Client) Connect(uri string) error {
	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial(uri, nil)
	return err
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// Listen starts listening for WebSocket messages
func (c *Client) Listen(stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			c.processMessage(message)
		}
	}
}

// processMessage handles incoming WebSocket messages
func (c *Client) processMessage(message []byte) {
	var tickerData models.TickerData
	if err := json.Unmarshal(message, &tickerData); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	formattedData := models.FormatTickerData(tickerData)

	// Print some information to the console
	log.Printf("Received data for %s - Price: %.2f, Change: %.2f, Volume: %.2f",
		formattedData.Symbol,
		formattedData.LastPrice,
		formattedData.PriceChange,
		formattedData.Volume)

	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, proc := range c.processors {
		proc.Process(formattedData)
	}
}
