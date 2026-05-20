.PHONY: help build run test db-up db-down db-reset docker-up docker-down docker-logs

help:
	@echo "Available commands:"
	@echo "  make build           - Build the binary"
	@echo "  make run             - Run the application"
	@echo "  make test            - Run tests"
	@echo "  make db-up           - Start PostgreSQL container"
	@echo "  make db-down         - Stop PostgreSQL container"
	@echo "  make db-reset        - Reset database"
	@echo "  make db-migrate      - Run migrations"
	@echo "  make db-migrate-revert - Revert migrations"
	@echo "  make docker-up       - Start all Docker services"
	@echo "  make docker-down     - Stop all Docker services"
	@echo "  make docker-logs     - View Docker logs"

build:
	@echo "Building application..."
	go build -o bin/quoteyouros-backend ./cmd

run:
	@echo "Running application..."
	go run ./cmd

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

db-up:
	@echo "Starting PostgreSQL container..."
	docker-compose up -d postgres

db-down:
	@echo "Stopping PostgreSQL container..."
	docker-compose down

db-reset: db-migrate-revert db-migrate
	@echo "Database reset complete"

db-migrate:
	@echo "Running migrations..."
	@powershell -Command "Start-Sleep -Seconds 2"
	docker exec -i quoteyouros_db psql -U postgres -d quoteyouros < migrations/001_init_schema.sql

db-migrate-revert:
	@echo "Reverting migrations..."
	docker exec -i quoteyouros_db psql -U postgres -d quoteyouros < migrations/002_revert_schema.sql

docker-up:
	@echo "Starting all Docker services..."
	docker-compose up

docker-up-rebuild:
	@echo "Starting all Docker services..."
	docker-compose up --build 

docker-down:
	@echo "Stopping all Docker services..."
	docker-compose down

docker-logs:
	@echo "Viewing Docker logs..."
	docker-compose logs -f

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./...

vet:
	@echo "Running go vet..."
	go vet ./...
