# Deployment Guide - Sistema de Reservas

## Overview

This guide covers deployment options for the Sistema de Reservas using Docker Compose
for development/testing and Kubernetes for production environments.

## Prerequisites

### System Requirements

- **CPU**: 2+ cores
- **Memory**: 4GB+ RAM
- **Storage**: 10GB+ available space
- **Operating System**: Linux, macOS, or Windows with WSL2

### Required Software

#### For Docker Deployment

- Docker Engine 20.10+
- Docker Compose 2.0+

#### For Kubernetes Deployment

- Kubernetes cluster 1.24+
- kubectl CLI tool
- Helm 3.0+ (optional)

#### Development Tools

- Git
- Go 1.24.4+ (for local development)
- curl or Postman (for API testing)

---

## Docker Compose Deployment

### Quick Start

1. **Clone the repository**:

   ```bash
   git clone <repository-url>
   cd SISTEMA_RESERVAS
   ```

2. **Run the deployment script**:

   ```bash
   chmod +x scripts/deploy-docker.sh
   ./scripts/deploy-docker.sh
   ```

3. **Verify deployment**:

   ```bash
   curl http://localhost:8080/health
   ```

### Manual Docker Compose Deployment

1. **Build all images**:

   ```bash
   chmod +x scripts/build-images.sh
   ./scripts/build-images.sh
   ```

2. **Start services**:

   ```bash
   cd infrastructure
   docker-compose up -d
   ```

3. **Check service status**:

   ```bash
   docker-compose ps
   docker-compose logs -f
   ```

### Service URLs (Docker Compose)

- **API Gateway**: <http://localhost:8080>
- **User Service**: <http://localhost:8081>
- **Resource Service**: <http://localhost:8082>
- **Booking Service**: <http://localhost:8083>
- **Notification Service**: <http://localhost:8084>
- **Grafana**: <http://localhost:3000> (admin/admin)
- **Prometheus**: <http://localhost:9090>
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

### Environment Configuration

Create `.env` file in the `infrastructure` directory:

```bash
# Database Configuration
POSTGRES_DB=reservations_db
POSTGRES_USER=admin
POSTGRES_PASSWORD=your_secure_password

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_change_in_production

# External Services (Production)
EMAIL_API_KEY=your_sendgrid_api_key
SMS_ACCOUNT_SID=your_twilio_sid
SMS_AUTH_TOKEN=your_twilio_token
```

---

## Kubernetes Deployment

### Cluster Setup

#### Option 1: Local Development (minikube)

```bash
# Install minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube && sudo mv minikube /usr/local/bin/

# Start cluster
minikube start --memory=4096 --cpus=2
minikube addons enable ingress
```

#### Option 2: Cloud Providers

**AWS EKS**:

```bash
eksctl create cluster --name reservation-system --version 1.24 --region us-west-2 --nodegroup-name workers --node-type t3.medium --nodes 3
```

**Google GKE**:

```bash
gcloud container clusters create reservation-system --zone us-central1-a --machine-type e2-medium --num-nodes 3
```

**Azure AKS**:

```bash
az aks create --resource-group myResourceGroup --name reservation-system --node-count 3 --node-vm-size Standard_B2s
```

### Kubernetes Deployment Steps

1. **Prepare images**:

   ```bash
   # Build images
   ./scripts/build-images.sh
   
   # Tag for registry (if using external registry)
   docker tag user-service:latest your-registry/user-service:latest
   docker tag resource-service:latest your-registry/resource-service:latest
   docker tag booking-service:latest your-registry/booking-service:latest
   docker tag notification-service:latest your-registry/notification-service:latest
   docker tag api-gateway:latest your-registry/api-gateway:latest
   
   # Push to registry
   docker push your-registry/user-service:latest
   # ... repeat for all services
   ```

2. **Deploy to Kubernetes**:

   ```bash
   chmod +x scripts/deploy-k8s.sh
   ./scripts/deploy-k8s.sh
   ```

3. **Verify deployment**:

   ```bash
   kubectl get all -n reservation-system
   kubectl get ingress -n reservation-system
   ```

### Manual Kubernetes Deployment

1. **Create namespace and configs**:

   ```bash
   kubectl apply -f kubernetes/configmaps.yaml
   ```

2. **Deploy databases**:

   ```bash
   kubectl apply -f kubernetes/postgres.yaml
   kubectl apply -f kubernetes/redis.yaml
   
   # Wait for databases
   kubectl wait --for=condition=available --timeout=300s deployment/postgres -n reservation-system
   kubectl wait --for=condition=available --timeout=300s deployment/redis -n reservation-system
   ```

3. **Deploy microservices**:

   ```bash
   kubectl apply -f kubernetes/user-service.yaml
   kubectl apply -f kubernetes/resource-service.yaml
   kubectl apply -f kubernetes/booking-service.yaml
   kubectl apply -f kubernetes/notification-service.yaml
   ```

4. **Deploy API Gateway**:

   ```bash
   kubectl apply -f kubernetes/api-gateway.yaml
   ```

### Accessing Services in Kubernetes

#### Option 1: LoadBalancer (Cloud)

```bash
kubectl get service api-gateway-service -n reservation-system
# Use EXTERNAL-IP to access the API
```

#### Option 2: NodePort (Local)

```bash
kubectl get service api-gateway-service -n reservation-system
# Use node IP + NodePort
```

#### Option 3: Port Forward (Development)

```bash
kubectl port-forward service/api-gateway-service 8080:8080 -n reservation-system
# Access via http://localhost:8080
```

#### Option 4: Ingress (Recommended)

```bash
# Add to /etc/hosts (local development)
echo "$(minikube ip) api.reservation-system.local" | sudo tee -a /etc/hosts

# Access via http://api.reservation-system.local
```

---

## Production Configuration

### Security Hardening

1. **Secrets Management**:

   ```bash
   # Create production secrets
   kubectl create secret generic app-secrets \
     --from-literal=db-password='your_secure_password' \
     --from-literal=jwt-secret='your_production_jwt_secret' \
     --from-literal=email-api-key='your_email_key' \
     -n reservation-system
   ```

2. **TLS Configuration**:

   ```yaml
   # Add to ingress configuration
   spec:
     tls:
     - hosts:
       - api.yourdomain.com
       secretName: tls-secret
   ```

3. **Resource Limits**:

   ```yaml
   resources:
     requests:
       memory: "128Mi"
       cpu: "100m"
     limits:
       memory: "512Mi"
       cpu: "500m"
   ```

4. **Network Policies**:

   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: reservation-system-policy
   spec:
     podSelector: {}
     policyTypes:
     - Ingress
     - Egress
     ingress:
     - from:
       - namespaceSelector:
           matchLabels:
             name: reservation-system
   ```

### Database Configuration

#### Production PostgreSQL Setup

```bash
# Using Helm for production PostgreSQL
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install postgres bitnami/postgresql \
  --set auth.postgresPassword=your_secure_password \
  --set auth.database=reservations_db \
  --set primary.persistence.size=100Gi \
  --namespace reservation-system
```

#### Redis Configuration

```bash
# Using Helm for production Redis
helm install redis bitnami/redis \
  --set auth.password=your_redis_password \
  --set master.persistence.size=10Gi \
  --namespace reservation-system
```

### Monitoring Setup

#### Prometheus and Grafana

```bash
# Add Prometheus Helm repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

# Install Prometheus
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set grafana.adminPassword=admin

# Access Grafana
kubectl port-forward service/prometheus-grafana 3000:80 -n monitoring
```

### Backup Strategy

#### Database Backup

```bash
# PostgreSQL backup
kubectl exec -it postgres-pod -n reservation-system -- pg_dump -U admin reservations_db > backup.sql

# Automated backup with CronJob
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: postgres-backup
            image: postgres:15-alpine
            command:
            - /bin/bash
            - -c
            - pg_dump -h postgres-service -U admin reservations_db > /backup/backup-$(date +%Y%m%d).sql
```

---

## Scaling and Performance

### Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: user-service-hpa
  namespace: reservation-system
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Vertical Pod Autoscaling

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: user-service-vpa
  namespace: reservation-system
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  updatePolicy:
    updateMode: "Auto"
```

---

## Troubleshooting

### Common Issues

#### Services Not Starting

```bash
# Check pod status
kubectl get pods -n reservation-system

# Check pod logs
kubectl logs -f deployment/user-service -n reservation-system

# Describe pod for events
kubectl describe pod <pod-name> -n reservation-system
```

#### Database Connection Issues

```bash
# Test database connectivity
kubectl exec -it deployment/user-service -n reservation-system -- /bin/sh
# Inside container: test connection to postgres-service:5432
```

#### API Gateway Not Accessible

```bash
# Check service status
kubectl get service api-gateway-service -n reservation-system

# Check ingress
kubectl get ingress -n reservation-system

# Test internal connectivity
kubectl exec -it deployment/api-gateway -n reservation-system -- curl localhost:8080/health
```

### Log Analysis

#### Centralized Logging with ELK Stack

```bash
# Deploy Elasticsearch
helm install elasticsearch elastic/elasticsearch --namespace logging --create-namespace

# Deploy Kibana
helm install kibana elastic/kibana --namespace logging

# Deploy Filebeat for log collection
helm install filebeat elastic/filebeat --namespace logging
```

#### Application Logs

```bash
# Follow logs for all services
kubectl logs -f -l app=user-service -n reservation-system

# Get logs from all pods
kubectl logs -f --selector="app in (user-service,resource-service,booking-service,notification-service)" -n reservation-system
```

### Performance Monitoring

#### Resource Usage

```bash
# Check resource usage
kubectl top pods -n reservation-system
kubectl top nodes

# Detailed metrics
kubectl get --raw /metrics
```

#### Database Performance

```bash
# PostgreSQL performance
kubectl exec -it deployment/postgres -n reservation-system -- psql -U admin -d reservations_db -c "
SELECT query, calls, total_time, mean_time 
FROM pg_stat_statements 
ORDER BY total_time DESC 
LIMIT 10;"
```

---

## Maintenance and Updates

### Rolling Updates

```bash
# Update image version
kubectl set image deployment/user-service user-service=user-service:v2.0.0 -n reservation-system

# Check rollout status
kubectl rollout status deployment/user-service -n reservation-system

# Rollback if needed
kubectl rollout undo deployment/user-service -n reservation-system
```

### Backup and Restore

```bash
# Create backup
kubectl exec postgres-pod -- pg_dump -U admin reservations_db > backup.sql

# Restore from backup
kubectl exec -i postgres-pod -- psql -U admin -d reservations_db < backup.sql
```

### Health Checks

```bash
# System health check
curl http://api.reservation-system.local/health

# Individual service health
kubectl exec -it deployment/user-service -n reservation-system -- curl localhost:8081/health
```

This guide provides comprehensive instructions for deploying the Sistema de Reservas in various environments,
from local development to production Kubernetes clusters.
