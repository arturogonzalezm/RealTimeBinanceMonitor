# Makefile for RealTimeBinanceMonitor

DOCKER_COMPOSE = docker-compose
DB_CONTAINER = realtimebinancemonitor-db-1
APP_CONTAINER = realtimebinancemonitor-app-1
DB_VOLUME = realtimebinancemonitor_pgdata

.PHONY: all build up down clean logs db-shell

# Default target
all: build up

rebuild: clean build up

# Build the Docker images without using cache
build:
	@echo "Building Docker images..."
	docker compose build --no-cache

# Start the Docker containers in detached mode
up:
	@echo "Starting Docker containers..."
	docker compose up -d

# Stop the Docker containers
down:
	@echo "Stopping Docker containers..."
	docker compose down

# Remove Docker volumes
clean:
	@echo "Removing Docker volumes..."
	docker compose down -v
	docker volume rm realtimebinancemonitor_pgdata || true

# Show logs for all containers
logs:
	@echo "Showing logs for all containers..."
	docker compose logs -f

db-shell:
	docker exec -it realtimebinancemonitor_pgdata psql -U postgres -d postgres

test:
	go test ./...

run:
	go run cmd/monitor/main.go

# Run the Go application
#run:
#	@echo "Running the Go application..."
#	DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_NAME=postgres go run cmd/monitor/main.go