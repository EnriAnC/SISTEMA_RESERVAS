# Guía de Despliegue - Sistema de Reservas

## Resumen

Esta guía cubre las opciones de despliegue para el Sistema de Reservas usando Docker Compose para desarrollo/pruebas y Kubernetes para entornos de producción.

## Prerequisitos

### Requisitos del Sistema

- **CPU**: 2+ núcleos
- **Memoria**: 4GB+ RAM
- **Almacenamiento**: 10GB+ espacio disponible
- **Sistema Operativo**: Linux, macOS, o Windows con WSL2

### Software Requerido

#### Para Despliegue con Docker

- Docker Engine 20.10+
- Docker Compose 2.0+

#### Para Despliegue en Kubernetes

- Clúster Kubernetes 1.24+
- Herramienta CLI kubectl
- Helm 3.0+ (opcional)

#### Herramientas de Desarrollo

- Git
- Go 1.21+ (para desarrollo local)
- curl o Postman (para pruebas de API)

---

## Despliegue con Docker Compose

### Inicio Rápido

1. **Clonar el repositorio**:

   ```bash
   git clone <url-repositorio>
   cd SISTEMA_RESERVAS
   ```

2. **Ejecutar el script de despliegue**:

   ```bash
   chmod +x scripts/deploy-docker.sh
   ./scripts/deploy-docker.sh
   ```

3. **Verificar despliegue**:

   ```bash
   curl http://localhost:8080/health
   ```

### Despliegue Manual con Docker Compose

1. **Construir todas las imágenes**:

   ```bash
   chmod +x scripts/build-images.sh
   ./scripts/build-images.sh
   ```

2. **Iniciar servicios**:

   ```bash
   cd infrastructure
   docker-compose up -d
   ```

3. **Verificar estado de servicios**:

   ```bash
   docker-compose ps
   docker-compose logs -f
   ```

### URLs de Servicios (Docker Compose)

- **API Gateway**: <http://localhost:8080>
- **Servicio de Usuarios**: <http://localhost:8081>
- **Servicio de Recursos**: <http://localhost:8082>
- **Servicio de Reservas**: <http://localhost:8083>
- **Servicio de Notificaciones**: <http://localhost:8084>
- **Grafana**: <http://localhost:3000> (admin/admin)
- **Prometheus**: <http://localhost:9090>
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

### Configuración de Entorno

Crear archivo `.env` en el directorio `infrastructure`:

```bash
# Configuración de Base de Datos
POSTGRES_DB=reservations_db
POSTGRES_USER=admin
POSTGRES_PASSWORD=tu_contraseña_segura

# Configuración JWT
JWT_SECRET=tu_clave_jwt_super_secreta_cambiar_en_produccion

# Servicios Externos (Producción)
EMAIL_API_KEY=tu_clave_api_sendgrid
SMS_ACCOUNT_SID=tu_sid_twilio
SMS_AUTH_TOKEN=tu_token_twilio
```

---

## Despliegue en Kubernetes

### Configuración del Clúster

#### Opción 1: Desarrollo Local (minikube)

```bash
# Instalar minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube && sudo mv minikube /usr/local/bin/

# Iniciar clúster
minikube start --memory=4096 --cpus=2
minikube addons enable ingress
```

#### Opción 2: Proveedores de Nube

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

### Pasos de Despliegue en Kubernetes

1. **Preparar imágenes**:

   ```bash
   # Construir imágenes
   ./scripts/build-images.sh
   
   # Etiquetar para registro (si se usa registro externo)
   docker tag user-service:latest tu-registro/user-service:latest
   docker tag resource-service:latest tu-registro/resource-service:latest
   docker tag booking-service:latest tu-registro/booking-service:latest
   docker tag notification-service:latest tu-registro/notification-service:latest
   docker tag api-gateway:latest tu-registro/api-gateway:latest
   
   # Subir a registro
   docker push tu-registro/user-service:latest
   # ... repetir para todos los servicios
   ```

2. **Desplegar en Kubernetes**:

   ```bash
   chmod +x scripts/deploy-k8s.sh
   ./scripts/deploy-k8s.sh
   ```

3. **Verificar despliegue**:

   ```bash
   kubectl get all -n reservation-system
   kubectl get ingress -n reservation-system
   ```

### Despliegue Manual en Kubernetes

1. **Crear namespace y configuraciones**:

   ```bash
   kubectl apply -f kubernetes/configmaps.yaml
   ```

2. **Desplegar bases de datos**:

   ```bash
   kubectl apply -f kubernetes/postgres.yaml
   kubectl apply -f kubernetes/redis.yaml
   
   # Esperar por las bases de datos
   kubectl wait --for=condition=available --timeout=300s deployment/postgres -n reservation-system
   kubectl wait --for=condition=available --timeout=300s deployment/redis -n reservation-system
   ```

3. **Desplegar microservicios**:

   ```bash
   kubectl apply -f kubernetes/user-service.yaml
   kubectl apply -f kubernetes/resource-service.yaml
   kubectl apply -f kubernetes/booking-service.yaml
   kubectl apply -f kubernetes/notification-service.yaml
   ```

4. **Desplegar API Gateway**:

   ```bash
   kubectl apply -f kubernetes/api-gateway.yaml
   ```

### Acceso a Servicios en Kubernetes

#### Opción 1: LoadBalancer (Nube)

```bash
kubectl get service api-gateway-service -n reservation-system
# Usar EXTERNAL-IP para acceder a la API
```

#### Opción 2: NodePort (Local)

```bash
kubectl get service api-gateway-service -n reservation-system
# Usar IP del nodo + NodePort
```

#### Opción 3: Port Forward (Desarrollo)

```bash
kubectl port-forward service/api-gateway-service 8080:8080 -n reservation-system
# Acceder vía http://localhost:8080
```

#### Opción 4: Ingress (Recomendado)

```bash
# Agregar a /etc/hosts (desarrollo local)
echo "$(minikube ip) api.reservation-system.local" | sudo tee -a /etc/hosts

# Acceder vía http://api.reservation-system.local
```

---

## Configuración de Producción

### Endurecimiento de Seguridad

1. **Gestión de Secretos**:

   ```bash
   # Crear secretos de producción
   kubectl create secret generic app-secrets \
     --from-literal=db-password='tu_contraseña_segura' \
     --from-literal=jwt-secret='tu_jwt_secreto_produccion' \
     --from-literal=email-api-key='tu_clave_email' \
     -n reservation-system
   ```

2. **Configuración TLS**:

   ```yaml
   # Agregar a configuración de ingress
   spec:
     tls:
     - hosts:
       - api.tudominio.com
       secretName: tls-secret
   ```

3. **Límites de Recursos**:

   ```yaml
   resources:
     requests:
       memory: "128Mi"
       cpu: "100m"
     limits:
       memory: "512Mi"
       cpu: "500m"
   ```

4. **Políticas de Red**:

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

### Configuración de Base de Datos

#### Configuración PostgreSQL de Producción

```bash
# Usar Helm para PostgreSQL de producción
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install postgres bitnami/postgresql \
  --set auth.postgresPassword=tu_contraseña_segura \
  --set auth.database=reservations_db \
  --set primary.persistence.size=100Gi \
  --namespace reservation-system
```

#### Configuración Redis

```bash
# Usar Helm para Redis de producción
helm install redis bitnami/redis \
  --set auth.password=tu_contraseña_redis \
  --set master.persistence.size=10Gi \
  --namespace reservation-system
```

### Configuración de Monitorización

#### Prometheus y Grafana

```bash
# Agregar repositorio Helm de Prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

# Instalar Prometheus
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set grafana.adminPassword=admin

# Acceder a Grafana
kubectl port-forward service/prometheus-grafana 3000:80 -n monitoring
```

### Estrategia de Backup

#### Backup de Base de Datos

```bash
# Backup PostgreSQL
kubectl exec -it postgres-pod -n reservation-system -- pg_dump -U admin reservations_db > backup.sql

# Backup automatizado con CronJob
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
spec:
  schedule: "0 2 * * *"  # Diario a las 2 AM
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

## Escalado y Rendimiento

### Escalado Automático Horizontal de Pods

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

### Escalado Automático Vertical de Pods

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

## Solución de Problemas

### Problemas Comunes

#### Servicios No Inician

```bash
# Verificar estado de pods
kubectl get pods -n reservation-system

# Verificar logs de pods
kubectl logs -f deployment/user-service -n reservation-system

# Describir pod para eventos
kubectl describe pod <nombre-pod> -n reservation-system
```

#### Problemas de Conexión a Base de Datos

```bash
# Probar conectividad de base de datos
kubectl exec -it deployment/user-service -n reservation-system -- /bin/sh
# Dentro del contenedor: probar conexión a postgres-service:5432
```

#### API Gateway No Accesible

```bash
# Verificar estado del servicio
kubectl get service api-gateway-service -n reservation-system

# Verificar ingress
kubectl get ingress -n reservation-system

# Probar conectividad interna
kubectl exec -it deployment/api-gateway -n reservation-system -- curl localhost:8080/health
```

### Análisis de Logs

#### Logging Centralizado con ELK Stack

```bash
# Desplegar Elasticsearch
helm install elasticsearch elastic/elasticsearch --namespace logging --create-namespace

# Desplegar Kibana
helm install kibana elastic/kibana --namespace logging

# Desplegar Filebeat para recolección de logs
helm install filebeat elastic/filebeat --namespace logging
```

#### Logs de Aplicación

```bash
# Seguir logs para todos los servicios
kubectl logs -f -l app=user-service -n reservation-system

# Obtener logs de todos los pods
kubectl logs -f --selector="app in (user-service,resource-service,booking-service,notification-service)" -n reservation-system
```

### Monitorización de Rendimiento

#### Uso de Recursos

```bash
# Verificar uso de recursos
kubectl top pods -n reservation-system
kubectl top nodes

# Métricas detalladas
kubectl get --raw /metrics
```

#### Rendimiento de Base de Datos

```bash
# Rendimiento PostgreSQL
kubectl exec -it deployment/postgres -n reservation-system -- psql -U admin -d reservations_db -c "
SELECT query, calls, total_time, mean_time 
FROM pg_stat_statements 
ORDER BY total_time DESC 
LIMIT 10;"
```

---

## Mantenimiento y Actualizaciones

### Actualizaciones Graduales

```bash
# Actualizar versión de imagen
kubectl set image deployment/user-service user-service=user-service:v2.0.0 -n reservation-system

# Verificar estado de rollout
kubectl rollout status deployment/user-service -n reservation-system

# Rollback si es necesario
kubectl rollout undo deployment/user-service -n reservation-system
```

### Backup y Restauración

```bash
# Crear backup
kubectl exec postgres-pod -- pg_dump -U admin reservations_db > backup.sql

# Restaurar desde backup
kubectl exec -i postgres-pod -- psql -U admin -d reservations_db < backup.sql
```

### Verificaciones de Salud

```bash
# Verificación de salud del sistema
curl http://api.reservation-system.local/health

# Salud de servicio individual
kubectl exec -it deployment/user-service -n reservation-system -- curl localhost:8081/health
```

Esta guía de despliegue proporciona instrucciones completas para desplegar el Sistema de Reservas en varios entornos,
desde desarrollo local hasta clústeres de Kubernetes de producción.
