.PHONY: help build clean test deps proto docker-build docker-up docker-down docker-logs run-services stop-services service-status service-logs

BINARY_DIR=bin
DOCKER_COMPOSE_DIR=deployments/docker-compose
SERVICES=web_server user_and_post newsfeed

help:
	@echo "Available commands:"
	@echo "  build          - Build all services"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run all tests"
	@echo "  deps           - Download dependencies"
	@echo "  proto          - Generate protobuf code"
	@echo "  run-services   - Start all services"
	@echo "  stop-services  - Stop all services"
	@echo "  service-status - Show service status"
	@echo "  service-logs   - Show service logs (usage: make service-logs SERVICE=web_server)"
	@echo "  docker-build  - Build all Docker images"
	@echo "  docker-up     - Start all services with Docker Compose"
	@echo "  docker-down   - Stop all services"
	@echo "  docker-logs   - Show Docker logs"

build: clean
	@mkdir -p $(BINARY_DIR)
	@for service in $(SERVICES); do \
		go build -o $(BINARY_DIR)/$$service ./cmd/$$service; \
	done

clean:
	@rm -rf $(BINARY_DIR)
	@go clean -cache

test:
	@go test -v ./...

deps:
	@go mod download
	@go mod tidy

proto:
	@protoc --experimental_allow_proto3_optional --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/interfaces/proto/protobuf/*/*.proto

docker-build:
	@cd $(DOCKER_COMPOSE_DIR) && docker-compose build

docker-up:
	@cd $(DOCKER_COMPOSE_DIR) && docker-compose up -d

docker-down:
	@cd $(DOCKER_COMPOSE_DIR) && docker-compose down

docker-logs:
	@cd $(DOCKER_COMPOSE_DIR) && docker-compose logs -f

run-services:
	@./scripts/run-services.sh start

stop-services:
	@./scripts/run-services.sh stop

service-status:
	@./scripts/run-services.sh status

service-logs:
	@./scripts/run-services.sh logs $(SERVICE)
