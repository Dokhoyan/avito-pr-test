TEST_FLAGS ?= -v -race -parallel 5 -shuffle=on
COVER_FLAGS ?= -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
BINARY_NAME ?= bin/http-server

.PHONY: docker-up docker-clean test lint deps build clean mock docker-build loadtest loadtest-docker
.DEFAULT_GOAL := help

docker-up:
	docker compose up -d

docker-build:
	docker compose build

docker-clean:
	docker compose down
	docker image prune -f

test:
	go test $(TEST_FLAGS) $(COVER_FLAGS) ./...

lint:
	golangci-lint run ./...

deps:
	go mod download
	go mod verify
	go mod tidy

build: deps
	go build -o $(BINARY_NAME) ./cmd/http-server

clean:
	rm -rf bin/
	rm -f cover.out

mock:
	mockery --config config/.mockery.yaml

loadtest:
	@mkdir -p loadtest
	@echo "Installing vegeta (if not installed)..."
	@go install github.com/tsenart/vegeta/v12@latest
	@echo "Running load test..."
	@BASE_URL=$(BASE_URL) go run loadtest/loadtest.go

loadtest-docker:
	@echo "Docker mode not needed - using native Go tool"
	@$(MAKE) loadtest
	
help:
	@echo "Available targets:"
	@echo "  docker-up       - Start docker containers"
	@echo "  docker-clean    - Clean docker containers and images"
	@echo "  docker-build    - Build docker images"
	@echo "  test            - Run tests with race detection and coverage"
	@echo "  lint            - Run golangci-lint"
	@echo "  deps            - Download dependencies"
	@echo "  build           - Build application"
	@echo "  clean           - Clean build artifacts"
	@echo "  mock            - Generate mocks using mockery with config"
	@echo "  loadtest        - Run load tests with k6 (local or Docker)"
	@echo "  loadtest-docker - Run load tests using Docker"