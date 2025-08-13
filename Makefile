.PHONY: build run test clean deps fmt lint vet staticcheck coverage migrate-up migrate-down

# Variables
APP_NAME=product-service
BINARY_NAME=bin/app
CMD_PATH=cmd/main.go

# Build commands
build:
	go build -o $(BINARY_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

clean:
	rm -rf bin/

# Development commands
deps:
	go mod tidy
	go mod download

fmt:
	gofmt -l -w .

lint:
	golangci-lint run

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

# Testing commands
test:
	go test ./... -v

test-race:
	go test ./... -v -race -count=1

test-coverage:
	go test -cover ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Database commands
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

# Full validation pipeline
validate: fmt deps lint vet staticcheck test-race

# Docker commands
docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -p 8080:8080 --env-file .env $(APP_NAME)

# Docker Compose commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Development environment (PostgreSQL only)
dev-up:
	docker-compose -f docker-compose.dev.yaml up -d

dev-down:
	docker-compose -f docker-compose.dev.yaml down

dev-logs:
	docker-compose -f docker-compose.dev.yaml logs -f

# Database operations
db-migrate-up:
	migrate -path migrations -database "postgres://app_user:app_password@localhost:5432/product_db?sslmode=disable" up

db-migrate-down:
	migrate -path migrations -database "postgres://app_user:app_password@localhost:5432/product_db?sslmode=disable" down

db-migrate-force:
	migrate -path migrations -database "postgres://app_user:app_password@localhost:5432/product_db?sslmode=disable" force $(VERSION)

# Quick start for development
dev-start: dev-up
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	@echo "Running migrations..."
	@make db-migrate-up || true
	@echo "Starting application..."
	@go run cmd/main.go