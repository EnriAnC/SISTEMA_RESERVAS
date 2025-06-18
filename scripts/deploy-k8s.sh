#!/bin/bash

# Deploy Reservation System to Kubernetes
# This script deploys all services to a Kubernetes cluster

set -e

echo "🚀 Starting Kubernetes deployment for Sistema de Reservas..."

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if we can connect to the cluster
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

# Navigate to the kubernetes directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "$SCRIPT_DIR/../kubernetes"

echo "📁 Current directory: $(pwd)"

# Apply namespace and configurations
echo "📋 Creating namespace and configurations..."
kubectl apply -f configmaps.yaml

# Deploy databases
echo "💾 Deploying databases..."
kubectl apply -f postgres.yaml
kubectl apply -f redis.yaml

# Wait for databases to be ready
echo "⏳ Waiting for databases to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/postgres -n reservation-system
kubectl wait --for=condition=available --timeout=300s deployment/redis -n reservation-system

# Deploy microservices
echo "🔧 Deploying microservices..."
kubectl apply -f user-service.yaml
kubectl apply -f resource-service.yaml
kubectl apply -f booking-service.yaml
kubectl apply -f notification-service.yaml

# Wait for services to be ready
echo "⏳ Waiting for microservices to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/user-service -n reservation-system
kubectl wait --for=condition=available --timeout=300s deployment/resource-service -n reservation-system
kubectl wait --for=condition=available --timeout=300s deployment/booking-service -n reservation-system
kubectl wait --for=condition=available --timeout=300s deployment/notification-service -n reservation-system

# Deploy API Gateway
echo "🌐 Deploying API Gateway..."
kubectl apply -f api-gateway.yaml

# Wait for API Gateway to be ready
echo "⏳ Waiting for API Gateway to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/api-gateway -n reservation-system

# Get service information
echo "📊 Getting service status..."
kubectl get all -n reservation-system

# Get external IP for API Gateway
echo "🔍 Getting API Gateway external access..."
EXTERNAL_IP=$(kubectl get service api-gateway-service -n reservation-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "pending")
EXTERNAL_PORT=$(kubectl get service api-gateway-service -n reservation-system -o jsonpath='{.spec.ports[0].port}')

if [ "$EXTERNAL_IP" = "pending" ] || [ -z "$EXTERNAL_IP" ]; then
    echo "⏳ External IP is still pending. You can check later with:"
    echo "   kubectl get service api-gateway-service -n reservation-system"
    
    # Try to get NodePort information
    NODE_PORT=$(kubectl get service api-gateway-service -n reservation-system -o jsonpath='{.spec.ports[0].nodePort}' 2>/dev/null || echo "")
    if [ -z "$NODE_PORT" ]; then
        NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="ExternalIP")].address}' 2>/dev/null || kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
        echo "   Or try NodePort access: http://$NODE_IP:$NODE_PORT"
    fi
else
    echo "✅ API Gateway is accessible at: http://$EXTERNAL_IP:$EXTERNAL_PORT"
fi

# Show port-forward command for local access
echo ""
echo "🔧 For local access, use port-forwarding:"
echo "   kubectl port-forward service/api-gateway-service 8080:8080 -n reservation-system"
echo "   Then access: http://localhost:8080"

echo ""
echo "🎉 Kubernetes deployment completed successfully!"
echo ""
echo "📋 Useful commands:"
echo "   Check pods:     kubectl get pods -n reservation-system"
echo "   Check services: kubectl get services -n reservation-system"
echo "   View logs:      kubectl logs -f deployment/[service-name] -n reservation-system"
echo "   Delete all:     kubectl delete namespace reservation-system"
echo ""
echo "✨ Sistema de Reservas is now running on Kubernetes!"
