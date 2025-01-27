# Social Network

A microservices-based social network system built with Go, featuring user management, posts, comments, and real-time newsfeed

## Quick Start

```bash
# Build all services
make build

# Run all services locally
make run-services

# Check service status
make service-status

# Stop all services
make stop-services
```

## Architecture

- **Web Server** (Port 8080): API Gateway with REST endpoints
- **User & Post Service** (Port 8001): User management and post operations
- **Newsfeed Service** (Port 8002): Real-time newsfeed generation
- **MySQL**: Primary database
- **Redis**: Caching layer

## Tech Stack

- **Backend**: Go 1.21+, Gin, GORM
- **Database**: MySQL 8.0
- **Cache**: Redis
- **Communication**: gRPC
- **Containerization**: Docker & Docker Compose

## Development

```bash
# Generate protobuf files
make proto

# Build specific service
make build-user-and-post
make build-newsfeed
make build-web-server

# View logs
make service-logs <service_name>
```

## Configuration

Edit `configs/config.yml` to modify service settings, database connections, and ports.
