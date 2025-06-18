# Guía de Desarrollo

Esta guía proporciona instrucciones para configurar el entorno de desarrollo y contribuir al proyecto Sistema de Reservas.

## Prerequisitos

### Software Requerido

- **Go 1.21+**: [Instalar Go](https://golang.org/doc/install)
- **Docker**: [Instalar Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Instalar Docker Compose](https://docs.docker.com/compose/install/)
- **kubectl**: [Instalar kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- **Git**: [Instalar Git](https://git-scm.com/downloads)

### Herramientas Opcionales

- **Postman/Insomnia**: Para pruebas de API
- **pgAdmin**: Para gestión de base de datos PostgreSQL
- **Lens**: IDE de Kubernetes para gestión de clústeres

## Configuración del Entorno de Desarrollo

### 1. Clonar el Repositorio

```bash
git clone <url-repositorio>
cd SISTEMA_RESERVAS
```

### 2. Configurar Dependencias Locales

```bash
# Instalar dependencias Go para todos los servicios
cd services/user-service && go mod tidy && cd ../..
cd services/resource-service && go mod tidy && cd ../..
cd services/booking-service && go mod tidy && cd ../..
cd services/notification-service && go mod tidy && cd ../..
```

### 3. Iniciar Servicios de Infraestructura

```bash
# Iniciar PostgreSQL, Redis, y stack de monitorización
docker-compose up -d postgres redis prometheus grafana
```

### 4. Inicializar Base de Datos

```bash
# Aplicar esquema de base de datos
docker-compose exec postgres psql -U reservas_user -d reservas_db -f /docker-entrypoint-initdb.d/init.sql
```

## Ejecutar Servicios

### Opción 1: Ejecutar Servicios Nativamente

```bash
# Terminal 1 - Servicio de Usuarios
cd services/user-service
export PORT=8081
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=reservas_user
export DB_PASSWORD=reservas_pass
export DB_NAME=reservas_db
go run .

# Terminal 2 - Servicio de Recursos
cd services/resource-service
export PORT=8082
# ... (misma configuración DB)
go run .

# Terminal 3 - Servicio de Reservas
cd services/booking-service
export PORT=8083
# ... (misma configuración DB)
go run .

# Terminal 4 - Servicio de Notificaciones
cd services/notification-service
export PORT=8084
go run .

# Terminal 5 - API Gateway
cd api-gateway
docker run -p 8080:8080 -v $PWD:/etc/krakend/ devopsfaith/krakend run --config /etc/krakend/krakend.json
```

### Opción 2: Ejecutar con Docker Compose

```bash
# Construir e iniciar todos los servicios
docker-compose up --build
```

## Pruebas

### Pruebas Unitarias

```bash
# Ejecutar pruebas para todos los servicios
./scripts/run-tests.sh

# O ejecutar individualmente
cd services/user-service && go test ./...
cd services/resource-service && go test ./...
cd services/booking-service && go test ./...
cd services/notification-service && go test ./...
```

### Pruebas de Integración

```bash
# Iniciar entorno de pruebas
docker-compose -f docker-compose.test.yml up -d

# Ejecutar pruebas de integración
go test -tags=integration ./tests/...

# Limpiar
docker-compose -f docker-compose.test.yml down
```

### Pruebas de API

Usar la colección de Postman proporcionada o comandos curl:

```bash
# Verificaciones de salud
curl http://localhost:8080/users/health
curl http://localhost:8080/resources/health
curl http://localhost:8080/bookings/health
curl http://localhost:8080/notifications/health

# Crear un usuario
curl -X POST http://localhost:8080/users/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan Pérez","email":"juan@ejemplo.com","password":"password123"}'
```

## Estilo de Código y Estándares

### Estilo de Código Go

- Seguir las guías de [Go Efectivo](https://golang.org/doc/effective_go.html)
- Usar `gofmt` para formateo
- Usar `golint` para linting
- Usar `go vet` para análisis de código

### Convención de Commits

Seguir [Commits Convencionales](https://www.conventionalcommits.org/):

```Commit
feat(user-service): agregar endpoint de registro de usuario
fix(booking-service): resolver problema de doble reserva
docs(api): actualizar documentación de endpoints de reserva
```

### Lista de Verificación de Revisión de Código

- [ ] El código sigue las mejores prácticas de Go
- [ ] Las pruebas están incluidas y pasan
- [ ] La documentación está actualizada
- [ ] Los cambios de API están documentados
- [ ] Docker construye exitosamente
- [ ] No hay datos sensibles en el código

## Depuración

### Logs de Servicios

```bash
# Logs de Docker Compose
docker-compose logs -f user-service
docker-compose logs -f resource-service
docker-compose logs -f booking-service
docker-compose logs -f notification-service

# Logs de Kubernetes
kubectl logs -f deployment/user-service
kubectl logs -f deployment/resource-service
kubectl logs -f deployment/booking-service
kubectl logs -f deployment/notification-service
```

### Acceso a Base de Datos

```bash
# Conectar a PostgreSQL
docker-compose exec postgres psql -U reservas_user -d reservas_db

# Consultas comunes
SELECT * FROM users LIMIT 10;
SELECT * FROM resources WHERE available = true;
SELECT * FROM bookings WHERE status = 'confirmed';
```

### Monitorización

- **Prometheus**: <http://localhost:9090>
- **Grafana**: <http://localhost:3000> (admin/admin)
- **Métricas de Aplicación**: Cada servicio expone endpoint `/metrics`

## Agregar Nuevas Funcionalidades

### 1. Modificaciones de Servicio

1. Actualizar modelos en `models.go`
2. Agregar lógica de negocio en `service.go`
3. Actualizar repositorio en `repository.go`
4. Agregar manejadores HTTP en `handlers.go`
5. Actualizar rutas en `main.go`

### 2. Cambios de Base de Datos

1. Crear script de migración en `infrastructure/database/migrations/`
2. Actualizar `init.sql` para instalaciones frescas
3. Probar migración en base de datos de desarrollo

### 3. Actualizaciones de API Gateway

1. Actualizar `api-gateway/krakend.json`
2. Agregar nuevas rutas y backends
3. Actualizar limitación de velocidad si es necesario

### 4. Actualizaciones de Documentación

1. Actualizar documentación de API en `docs/API.md`
2. Actualizar diagramas de arquitectura si es necesario
3. Actualizar guías de despliegue

## Solución de Problemas

### Problemas Comunes

#### Puerto Ya en Uso

```bash
# Encontrar proceso usando puerto
lsof -i :8081
# Matar proceso
kill -9 <PID>
```

#### Problemas de Conexión a Base de Datos

```bash
# Verificar estado de base de datos
docker-compose ps postgres
# Reiniciar base de datos
docker-compose restart postgres
```

#### Problemas de Construcción Docker

```bash
# Limpiar caché de Docker
docker system prune -f
# Reconstruir sin caché
docker-compose build --no-cache
```

#### Problemas de Módulos Go

```bash
# Limpiar caché de módulos
go clean -modcache
# Re-descargar dependencias
go mod download
```

### Problemas de Rendimiento

1. Verificar uso de recursos de servicios: `docker stats`
2. Monitorizar conexiones de base de datos
3. Verificar métricas de API Gateway
4. Revisar logs de aplicación para errores

## Integración Continua

### GitHub Actions

El proyecto incluye flujos de trabajo CI/CD:

- **Build**: Compila todos los servicios
- **Test**: Ejecuta pruebas unitarias y de integración
- **Security**: Escanea vulnerabilidades
- **Deploy**: Despliega a staging/producción

### Prueba Local de CI

```bash
# Ejecutar pipeline CI localmente con act (ejecutor local de GitHub Actions)
act push
```

## Contribuir

### Proceso de Pull Request

1. Crear rama de funcionalidad: `git checkout -b feature/nueva-funcionalidad`
2. Hacer cambios y agregar pruebas
3. Hacer commit usando formato de commit convencional
4. Push de rama y crear pull request
5. Abordar comentarios de revisión
6. Hacer merge después de aprobación

### Reporte de Problemas

Al reportar problemas, incluir:

- Servicio afectado
- Pasos para reproducir
- Comportamiento esperado vs real
- Logs y mensajes de error
- Detalles del entorno

## Herramientas de Desarrollo

### Extensiones Recomendadas de VS Code

- Extensión Go
- Extensión Docker
- Extensión Kubernetes
- REST Client
- GitLens
- Thunder Client (pruebas de API)

### Plantilla de Variables de Entorno

Crear archivo `.env` en cada directorio de servicio:

```bash
# Base de Datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=reservas_user
DB_PASSWORD=reservas_pass
DB_NAME=reservas_db

# Servicio
PORT=8081
ENV=development
LOG_LEVEL=debug

# Servicios externos (opcional)
REDIS_URL=localhost:6379
JWT_SECRET=tu-clave-secreta
```

## Recursos Adicionales

- [Documentación de Go](https://golang.org/doc/)
- [Documentación de Docker](https://docs.docker.com/)
- [Documentación de Kubernetes](https://kubernetes.io/docs/)
- [Documentación de KrakenD](https://www.krakend.io/docs/)
- [Documentación de PostgreSQL](https://www.postgresql.org/docs/)
- [Patrones de Microservicios](https://microservices.io/)

---

Para preguntas o soporte, por favor crear un issue en el repositorio o contactar al equipo de desarrollo.
