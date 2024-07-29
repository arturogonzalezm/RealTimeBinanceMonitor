package models

import (
	"strconv"
	"time"

	"github.com/arturogonzalezm/RealTimeBinanceMonitor/pkg/flexint"
)

// TickerData represents the structure of the incoming WebSocket message
type TickerData struct {
	EventType    string        `json:"e"`
	EventTime    flexint.Int64 `json:"E"`
	Symbol       string        `json:"s"`
	LastPrice    string        `json:"c"`
	PriceChange  string        `json:"p"`
	HighPrice    string        `json:"h"`
	LowPrice     string        `json:"l"`
	Volume       string        `json:"v"`
	QuoteVolume  string        `json:"q"`
	OpenTime     flexint.Int64 `json:"O"`
	CloseTime    flexint.Int64 `json:"C"`
	FirstTradeID int           `json:"F"`
	LastTradeID  int           `json:"L"`
	TradeCount   int           `json:"n"`
}

// FormattedData represents the structure of the processed data
type FormattedData struct {
	EventTime   int64
	Symbol      string
	LastPrice   float64
	PriceChange float64
	HighPrice   float64
	LowPrice    float64
	Volume      float64
	QuoteVolume float64
	OpenTime    int64
	CloseTime   int64
	TradeCount  int
	Latency     int64
}

// FormatTickerData converts TickerData to FormattedData
func FormatTickerData(td TickerData) FormattedData {
	return FormattedData{
		EventTime:   int64(td.EventTime),
		Symbol:      td.Symbol,
		LastPrice:   parseFloat(td.LastPrice),
		PriceChange: parseFloat(td.PriceChange),
		HighPrice:   parseFloat(td.HighPrice),
		LowPrice:    parseFloat(td.LowPrice),
		Volume:      parseFloat(td.Volume),
		QuoteVolume: parseFloat(td.QuoteVolume),
		OpenTime:    int64(td.OpenTime),
		CloseTime:   int64(td.CloseTime),
		TradeCount:  td.TradeCount,
		Latency:     time.Now().UnixMilli() - int64(td.EventTime),
	}
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
