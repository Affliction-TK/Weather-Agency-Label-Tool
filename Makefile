.PHONY: help build frontend-build backend-build run dev-frontend dev-backend clean test init-db

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

FRONTEND_DIR=frontend
BACKEND_DIR=backend
BACKEND_BIN_DIR=$(BACKEND_DIR)/bin
BACKEND_BIN=$(BACKEND_BIN_DIR)/server

build: frontend-build backend-build ## Build the frontend and backend artifacts
	@echo "Build complete!"

frontend-build: ## Build the frontend bundle
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && npm run build

backend-build: ## Build the backend server binary
	@echo "Building backend..."
	@mkdir -p $(BACKEND_BIN_DIR)
	@cd $(BACKEND_DIR) && go build -o bin/server .

run: build ## Build everything and run the backend server
	@echo "Starting server..."
	@cd $(BACKEND_DIR) && ./bin/server

dev-frontend: ## Run the frontend in development mode
	@cd $(FRONTEND_DIR) && npm run dev

dev-backend: ## Run the backend in development mode
	@cd $(BACKEND_DIR) && go run .

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf $(BACKEND_BIN_DIR)
	@rm -rf $(FRONTEND_DIR)/dist
	@echo "Clean complete!"

test: ## Run backend unit tests
	@cd $(BACKEND_DIR) && go test -v ./...

init-db: ## Initialize database with schema
	@echo "Initializing database..."
	@mysql -h 127.0.0.1 -u weather_user -pweather_password weather_label_db < $(BACKEND_DIR)/schema.sql
	@echo "Database initialized!"
