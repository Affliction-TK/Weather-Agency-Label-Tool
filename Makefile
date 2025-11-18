.PHONY: help setup build run dev clean test docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Run initial setup
	@bash setup.sh

build: ## Build the application
	@echo "Building frontend..."
	@cd frontend && npm run build
	@echo "Building backend..."
	@go build -o server main.go
	@echo "Build complete!"

run: build ## Build and run the application
	@echo "Starting server..."
	@./server

dev-frontend: ## Run frontend in development mode
	@cd frontend && npm run dev

dev-backend: ## Run backend in development mode
	@go run main.go

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f server
	@rm -rf frontend/dist
	@echo "Clean complete!"

test: ## Run tests
	@go test -v ./...

docker-up: ## Start MySQL database using docker-compose
	@docker-compose up -d
	@echo "Waiting for MySQL to be ready..."
	@sleep 10
	@echo "MySQL is ready!"

docker-down: ## Stop MySQL database
	@docker-compose down

docker-logs: ## View MySQL logs
	@docker-compose logs -f mysql

init-db: ## Initialize database with schema
	@echo "Initializing database..."
	@mysql -h 127.0.0.1 -u weather_user -pweather_password weather_label_db < schema.sql
	@echo "Database initialized!"

all: setup docker-up build ## Complete setup and build
	@echo "All tasks complete! Run 'make run' to start the server."
