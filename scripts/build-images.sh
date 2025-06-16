#!/bin/bash

# Build all Docker images for the reservation system

set -e

echo "🔨 Building all Docker images for Sistema de Reservas..."

# Navigate to project root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$SCRIPT_DIR/.."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Build User Service
echo "🏗️  Building User Service..."
cd "$PROJECT_ROOT/services/user-service"
docker build -t user-service:latest .
echo "✅ User Service built successfully"

# Build Resource Service
echo "🏗️  Building Resource Service..."
cd "$PROJECT_ROOT/services/resource-service"
docker build -t resource-service:latest .
echo "✅ Resource Service built successfully"

# Build Booking Service
echo "🏗️  Building Booking Service..."
cd "$PROJECT_ROOT/services/booking-service"
docker build -t booking-service:latest .
echo "✅ Booking Service built successfully"

# Build Notification Service
echo "🏗️  Building Notification Service..."
cd "$PROJECT_ROOT/services/notification-service"
docker build -t notification-service:latest .
echo "✅ Notification Service built successfully"

# Build API Gateway
echo "🏗️  Building API Gateway..."
cd "$PROJECT_ROOT/api-gateway"
docker build -t api-gateway:latest .
echo "✅ API Gateway built successfully"

# List built images
echo ""
echo "📋 Built images:"
docker images | grep -E "(user-service|resource-service|booking-service|notification-service|api-gateway)" | head -5

echo ""
echo "🎉 All images built successfully!"
echo ""
echo "🏷️  Image tags:"
echo "   user-service:latest"
echo "   resource-service:latest"
echo "   booking-service:latest"
echo "   notification-service:latest"
echo "   api-gateway:latest"
echo ""
echo "🚀 Ready for deployment!"
