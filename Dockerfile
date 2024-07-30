# Start with the official Go image as a build stage
FROM golang:1.22 AS build

# Set the Current Working Directory inside the container
WORKDIR /RealTimeBinanceMonitor

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Set the working directory to where the main.go file is located
WORKDIR /RealTimeBinanceMonitor/cmd/monitor

# Build the Go app
RUN go build -o /RealTimeBinanceMonitor/main .

# Start a new stage from a smaller image
FROM debian:bullseye-slim

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /RealTimeBinanceMonitor

# Copy the Pre-built binary file from the previous stage
COPY --from=build /RealTimeBinanceMonitor/main .

# Command to run the executable
CMD ["./main"]
