package monitor

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/processor"
)

func TestMonitorSymbol(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock dependencies
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockWebSocketClient := websocket_mocks.NewMockWebSocketClient(ctrl)
	mockPGWriter := processor_mocks.NewMockDataProcessor(ctrl)

	// Expectations
	mockWebSocketClient.EXPECT().Connect(gomock.Any()).Return(nil)
	mockWebSocketClient.EXPECT().Close().Return(nil)
	mockWebSocketClient.EXPECT().AddProcessor(mockPGWriter)
	mockWebSocketClient.EXPECT().Listen(gomock.Any()).Do(func(stop chan struct{}) {
		time.Sleep(1 * time.Second)
		close(stop)
	})

	mockPGWriter.EXPECT().Close().Return(nil)
	mockPGWriter.EXPECT().GetProcessedCount().AnyTimes().Return(0)
	mockPGWriter.EXPECT().GetBufferSize().AnyTimes().Return(0)

	stop := make(chan struct{})
	go MonitorSymbol("btcusdt", db, stop, mockWebSocketClient, mockPGWriter)

	select {
	case <-stop:
		// Test passed
	case <-time.After(2 * time.Second):
		t.Fatalf("Test timed out")
	}
}
