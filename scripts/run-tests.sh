#!/bin/bash

# Test script for Sistema de Reservas
# This script runs unit tests for all microservices

set -e

echo "🧪 Running Sistema de Reservas Test Suite"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to run tests for a service
run_service_tests() {
    local service_name=$1
    local service_path=$2
    
    echo -e "\n${BLUE}Testing $service_name...${NC}"
    echo "----------------------------------------"
    
    cd "$service_path"
    
    # Check if go.mod exists
    if [ ! -f "go.mod" ]; then
        echo -e "${YELLOW}⚠️  No go.mod found in $service_name, skipping...${NC}"
        cd - > /dev/null
        return
    fi
    
    # Install dependencies
    echo "Installing dependencies..."
    go mod tidy
    
    # Run tests
    if go test -v ./...; then
        echo -e "${GREEN}✅ $service_name tests passed${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ $service_name tests failed${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    cd - > /dev/null
}

# Function to run linting
run_service_lint() {
    local service_name=$1
    local service_path=$2
    
    echo -e "\n${BLUE}Linting $service_name...${NC}"
    echo "----------------------------------------"
    
    cd "$service_path"
    
    # Check if go.mod exists
    if [ ! -f "go.mod" ]; then
        echo -e "${YELLOW}⚠️  No go.mod found in $service_name, skipping lint...${NC}"
        cd - > /dev/null
        return
    fi
    
    # Run go vet
    if go vet ./...; then
        echo -e "${GREEN}✅ $service_name vet passed${NC}"
    else
        echo -e "${RED}❌ $service_name vet failed${NC}"
    fi
    
    # Run go fmt check
    if [ -z "$(gofmt -l .)" ]; then
        echo -e "${GREEN}✅ $service_name formatting is correct${NC}"
    else
        echo -e "${RED}❌ $service_name has formatting issues:${NC}"
        gofmt -l .
    fi
    
    cd - > /dev/null
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed or not in PATH${NC}"
    exit 1
fi

echo -e "${BLUE}Go version:${NC} $(go version)"

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}Project root:${NC} $PROJECT_ROOT"

# Change to project root
cd "$PROJECT_ROOT"

# Services to test
SERVICES=(
    "User Service:services/user-service"
    "Resource Service:services/resource-service"
    "Booking Service:services/booking-service"
    "Notification Service:services/notification-service"
)

# Run tests for each service
echo -e "\n${YELLOW}🔧 Running Unit Tests${NC}"
echo "====================="

for service in "${SERVICES[@]}"; do
    IFS=':' read -r service_name service_path <<< "$service"
    run_service_tests "$service_name" "$service_path"
done

# Run linting
echo -e "\n${YELLOW}🔍 Running Code Quality Checks${NC}"
echo "==============================="

for service in "${SERVICES[@]}"; do
    IFS=':' read -r service_name service_path <<< "$service"
    run_service_lint "$service_name" "$service_path"
done

# Run security checks if gosec is available
echo -e "\n${YELLOW}🔒 Running Security Checks${NC}"
echo "=========================="

if command -v gosec &> /dev/null; then
    for service in "${SERVICES[@]}"; do
        IFS=':' read -r service_name service_path <<< "$service"
        echo -e "\n${BLUE}Security scan for $service_name...${NC}"
        cd "$service_path"
        if [ -f "go.mod" ]; then
            gosec ./... || echo -e "${YELLOW}⚠️  Security warnings found in $service_name${NC}"
        fi
        cd - > /dev/null
    done
else
    echo -e "${YELLOW}⚠️  gosec not installed, skipping security checks${NC}"
    echo "Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
fi

# Test results summary
echo -e "\n${YELLOW}📊 Test Summary${NC}"
echo "==============="
echo -e "Total services tested: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}💥 Some tests failed!${NC}"
    exit 1
fi
