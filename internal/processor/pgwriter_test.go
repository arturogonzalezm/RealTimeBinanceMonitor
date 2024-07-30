package processor

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arturogonzalezm/RealTimeBinanceMonitor/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewPGWriter(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		err := db.Close()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}()

	writer, err := NewPGWriter(db)
	assert.NoError(t, err)
	assert.NotNil(t, writer)
}

func TestPGWriter_Process(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		err := db.Close()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}()

	writer, err := NewPGWriter(db)
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	data := models.FormattedData{
		EventTime:   1625097600000,
		Symbol:      "btcusdt",
		LastPrice:   34000.0,
		PriceChange: 100.0,
		HighPrice:   34500.0,
		LowPrice:    33500.0,
		Volume:      100.0,
		QuoteVolume: 3400000.0,
		OpenTime:    1625094000000,
		CloseTime:   1625097600000,
		TradeCount:  1000,
		Latency:     100,
	}

	mock.ExpectExec(`INSERT INTO ticker_data`).
		WithArgs(
			data.EventTime, data.Symbol, data.LastPrice, data.PriceChange, data.HighPrice, data.LowPrice,
			data.Volume, data.QuoteVolume, data.OpenTime, data.CloseTime, data.TradeCount, data.Latency,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	writer.Process(data)

	assert.Equal(t, 1, writer.GetProcessedCount())
}

func TestPGWriter_GetProcessedCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		err := db.Close()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}()

	writer, err := NewPGWriter(db)
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	assert.Equal(t, 0, writer.GetProcessedCount())
}

func TestPGWriter_GetBufferSize(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		err := db.Close()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}()

	writer, err := NewPGWriter(db)
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	assert.Equal(t, 0, writer.GetBufferSize())
}

func TestPGWriter_Close(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, mock.ExpectationsWereMet())
	}()

	writer, err := NewPGWriter(db)
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	mock.ExpectClose()

	err = writer.Close()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
