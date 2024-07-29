[![codecov](https://codecov.io/gh/arturogonzalezm/RealTimeBinanceMonitor/graph/badge.svg?token=I4cOxsac0y)](https://codecov.io/gh/arturogonzalezm/RealTimeBinanceMonitor)
![Go Test and Coverage](https://github.com/arturogonzalezm/RealTimeBinanceMonitor/actions/workflows/workflow.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-purple.svg)](https://opensource.org/licenses/MIT)

# RealTime Binance Monitor

RealTime Binance Monitor is a Go application that connects to the Binance WebSocket API to receive real-time cryptocurrency market data. It processes this data and saves it to a CSV file for further analysis.

## Features

- Real-time connection to Binance WebSocket API
- Processing of ticker data for specified cryptocurrency pairs
- Buffered writing of data to CSV files
- Configurable buffer size and flush interval
- Graceful shutdown handling

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/arturogonzalezm/RealTimeBinanceMonitor.git
   ```

2. Change to the project directory:
   ```bash
   cd RealTimeBinanceMonitor
   ```

3. Build the application:
   ```bash
   go build ./cmd/monitor
   ```

## Usage

Run the application with:

```bash
./monitor
```

By default, it will monitor the BTC/USDT pair. To monitor a different pair, modify the `symbol` variable in `cmd/monitor/main.go`.

## Configuration

You can adjust the following parameters in `cmd/monitor/main.go`:

- `symbol`: The cryptocurrency pair to monitor (e.g., "btcusdt", "ethusdt")
- Buffer size and flush interval in the `NewBufferedCSVWriter` function call

## TODO: 

- Add configuration file support
- Add support for multiple pairs
- Insert data into a database
- Add support for more Binance WebSocket API endpoints
- Add support for more data processing and analysis
- Add support for Kafka
- Add support for Vue JS front end
- Add support for Docker
- Add support for Kubernetes or doc
- Add support for CI/CD
- Add support for more tests
- Add support for more documentation


## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
