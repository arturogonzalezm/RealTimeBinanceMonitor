package processor

import "github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/models"

// DataProcessor interface for the Observer pattern
type DataProcessor interface {
	Process(data models.FormattedData)
	GetProcessedCount() int
	GetBufferSize() int
}
