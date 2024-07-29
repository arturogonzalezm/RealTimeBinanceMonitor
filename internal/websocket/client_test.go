package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/models"
	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/processor"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockProcessor implements the processor.DataProcessor interface for testing
type MockProcessor struct {
	ProcessedData []models.FormattedData
	bufferSize    int
}

func (m *MockProcessor) Process(data models.FormattedData) {
	m.ProcessedData = append(m.ProcessedData, data)
}

func (m *MockProcessor) GetProcessedCount() int {
	return len(m.ProcessedData)
}

func (m *MockProcessor) GetBufferSize() int {
	return m.bufferSize
}

// Implement any other methods required by the processor.DataProcessor interface
func (m *MockProcessor) Start() error {
	return nil
}

func (m *MockProcessor) Stop() error {
	return nil
}

// Ensure MockProcessor implements processor.DataProcessor
var _ processor.DataProcessor = (*MockProcessor)(nil)

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client, "NewClient() should not return nil")
	assert.Empty(t, client.processors, "New client should have no processors")
}

func TestAddProcessor(t *testing.T) {
	client := NewClient()
	mockProcessor := &MockProcessor{bufferSize: 100}
	client.AddProcessor(mockProcessor)

	assert.Len(t, client.processors, 1, "Client should have 1 processor after adding")
}

func TestConnect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		_, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err, "Failed to upgrade connection")
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient()
	err := client.Connect(wsURL)
	assert.NoError(t, err, "Connect should not return an error")
	defer client.Close()
}

func TestProcessMessage(t *testing.T) {
	client := NewClient()
	mockProcessor := &MockProcessor{bufferSize: 100}
	client.AddProcessor(mockProcessor)

	tickerData := models.TickerData{
		Symbol:    "BTCUSDT",
		LastPrice: "50000.00",
	}
	message, err := json.Marshal(tickerData)
	require.NoError(t, err, "Failed to marshal ticker data")

	client.processMessage(message)

	assert.Len(t, mockProcessor.ProcessedData, 1, "Should have processed 1 message")

	processedData := mockProcessor.ProcessedData[0]
	assert.Equal(t, "BTCUSDT", processedData.Symbol, "Processed symbol should match")
	assert.Equal(t, 50000.00, processedData.LastPrice, "Processed last price should match")
}

func TestListen(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err, "Failed to upgrade connection")
		defer conn.Close()

		tickerData := models.TickerData{
			Symbol:    "ETHUSDT",
			LastPrice: "3000.00",
		}
		message, err := json.Marshal(tickerData)
		require.NoError(t, err, "Failed to marshal ticker data")

		err = conn.WriteMessage(websocket.TextMessage, message)
		require.NoError(t, err, "Failed to write message")
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient()
	mockProcessor := &MockProcessor{bufferSize: 100}
	client.AddProcessor(mockProcessor)

	err := client.Connect(wsURL)
	require.NoError(t, err, "Failed to connect")
	defer client.Close()

	stop := make(chan struct{})
	go func() {
		client.Listen(stop)
	}()

	// Give some time for the message to be processed
	time.Sleep(100 * time.Millisecond)

	close(stop)

	assert.Len(t, mockProcessor.ProcessedData, 1, "Should have processed 1 message")

	processedData := mockProcessor.ProcessedData[0]
	assert.Equal(t, "ETHUSDT", processedData.Symbol, "Processed symbol should match")
	assert.Equal(t, 3000.00, processedData.LastPrice, "Processed last price should match")
}
