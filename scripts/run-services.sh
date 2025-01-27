#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

case "${1:-start}" in
    start)
        echo "🚀 Starting Social Network Services..."
        
        if [ ! -f "bin/web_server" ] || [ ! -f "bin/user_and_post" ] || [ ! -f "bin/newsfeed" ]; then
            print_error "Service binaries not found. Please run 'make build' first."
            exit 1
        fi
        print_success "All service binaries found"
        
        if [ ! -f "configs/config.yml" ]; then
            print_error "config.yml not found in configs/ directory."
            exit 1
        fi
        print_success "Configuration file found"
        
        mkdir -p logs
        
        print_info "Starting user_and_post on port 8001..."
        nohup bin/user_and_post -conf configs/config.yml > logs/user_and_post.log 2>&1 &
        user_pid=$!
        sleep 1
        if kill -0 $user_pid 2>/dev/null; then
            echo $user_pid > logs/user_and_post.pid
            print_success "user_and_post started successfully"
        else
            print_error "user_and_post failed to start"
            exit 1
        fi
        
        print_info "Starting newsfeed on port 8002..."
        nohup bin/newsfeed -conf configs/config.yml > logs/newsfeed.log 2>&1 &
        news_pid=$!
        sleep 1
        if kill -0 $news_pid 2>/dev/null; then
            echo $news_pid > logs/newsfeed.pid
            print_success "newsfeed started successfully"
        else
            print_error "newsfeed failed to start"
            exit 1
        fi
        
        print_info "Starting web_server on port 8080..."
        nohup bin/web_server -config configs/config.yml > logs/web_server.log 2>&1 &
        web_pid=$!
        sleep 1
        if kill -0 $web_pid 2>/dev/null; then
            echo $web_pid > logs/web_server.pid
            print_success "web_server started successfully"
        else
            print_error "web_server failed to start"
            exit 1
        fi
        
        sleep 2
        print_success "All services started successfully!"
        echo ""
        echo "Service endpoints:"
        echo "  Web Server (API Gateway): http://localhost:8080"
        echo "  User & Post Service:  localhost:8001"
        echo "  Newsfeed Service:     localhost:8002"
        echo ""
        echo "Use 'make run-services stop' to stop all services"
        echo "Use 'make run-services status' to check service status"
        echo "Use 'make run-services logs <service>' to view logs"
        ;;
    stop)
        print_info "Stopping all services..."
        for service in web_server user_and_post newsfeed; do
            if [ -f "logs/${service}.pid" ]; then
                pid=$(cat logs/${service}.pid)
                if kill -0 $pid 2>/dev/null; then
                    kill $pid
                    print_success "Stopped $service (PID: $pid)"
                fi
                rm -f logs/${service}.pid
            fi
        done
        print_success "All services stopped"
        ;;
    status)
        print_info "Service Status:"
        for service in web_server user_and_post newsfeed; do
            if [ -f "logs/${service}.pid" ]; then
                pid=$(cat logs/${service}.pid)
                if kill -0 $pid 2>/dev/null; then
                    print_success "$service: Running (PID: $pid)"
                else
                    print_warning "$service: Not running (stale PID file)"
                    rm -f logs/${service}.pid
                fi
            else
                print_warning "$service: Not running"
            fi
        done
        ;;
    logs)
        if [ -z "$2" ]; then
            print_error "Please specify a service name (web_server, user_and_post, or newsfeed)"
            exit 1
        fi
        if [ -f "logs/${2}.log" ]; then
            print_info "Showing logs for $2:"
            tail -f logs/${2}.log
        else
            print_error "Log file for $2 not found"
        fi
        ;;
    restart)
        $0 stop
        sleep 2
        $0 start
        ;;
    *)
        echo "Usage: $0 {start|stop|status|logs|restart}"
        exit 1
        ;;
esac
