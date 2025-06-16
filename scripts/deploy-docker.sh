#!/bin/bash

# Deploy Reservation System with Docker Compose
# This script builds and deploys all services using Docker Compose

set -e

echo "🚀 Starting Sistema de Reservas deployment..."

# Check if Docker is installed and running
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! docker info &> /dev/null; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Navigate to the infrastructure directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "$SCRIPT_DIR/../infrastructure"

echo "📁 Current directory: $(pwd)"

# Stop and remove existing containers
echo "🛑 Stopping existing containers..."
docker-compose down --volumes --remove-orphans

# Build all services
echo "🔨 Building all services..."
docker-compose build --no-cache

# Start all services
echo "🚀 Starting all services..."
docker-compose up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 30

# Check service health
echo "🔍 Checking service status..."
docker-compose ps

# Test API Gateway health
echo "🏥 Testing API Gateway health..."
for i in {1..30}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ API Gateway is healthy!"
        break
    fi
    echo "⏳ Waiting for API Gateway... (attempt $i/30)"
    sleep 5
    if [ $i -eq 30 ]; then
        echo "❌ API Gateway health check failed"
        exit 1
    fi
done

# Display service URLs
echo ""
echo "🎉 Deployment completed successfully!"
echo ""
echo "📍 Service URLs:"
echo "   API Gateway:    http://localhost:8080"
echo "   User Service:   http://localhost:8081"
echo "   Resource Service: http://localhost:8082"
echo "   Booking Service: http://localhost:8083"
echo "   Notification Service: http://localhost:8084"
echo "   Grafana:        http://localhost:3000 (admin/admin)"
echo "   Prometheus:     http://localhost:9090"
echo ""
echo "📋 API Documentation:"
echo "   GET /health     - Health check"
echo "   GET /api/v1/*   - API endpoints"
echo ""
echo "🔧 Management Commands:"
echo "   View logs:      docker-compose logs -f [service-name]"
echo "   Stop services:  docker-compose down"
echo "   Restart:        docker-compose restart [service-name]"
echo ""
echo "✨ Sistema de Reservas is now running!"
