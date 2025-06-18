# Booking Service

Microservicio central encargado de la lógica de negocio para crear, gestionar y validar reservas en el sistema.

## Funcionalidades

- **Gestión de Reservas**: CRUD completo de reservas
- **Validación de Disponibilidad**: Verificación de conflictos de horario
- **Estados de Reserva**: Gestión del ciclo de vida (PENDING → CONFIRMED → COMPLETED/CANCELLED)
- **Eventos**: Publicación de eventos para notificaciones
- **Consultas**: Filtrado por usuario, recurso, fecha y estado

## API Endpoints

### Reservas

- `POST /api/v1/bookings` - Crear reserva
- `GET /api/v1/bookings` - Listar reservas (con filtros)
- `GET /api/v1/bookings/{id}` - Obtener reserva por ID
- `PUT /api/v1/bookings/{id}` - Actualizar reserva
- `DELETE /api/v1/bookings/{id}` - Cancelar reserva
- `POST /api/v1/bookings/{id}/confirm` - Confirmar reserva

### Consultas Específicas

- `GET /api/v1/users/{user_id}/bookings` - Reservas de un usuario
- `POST /api/v1/bookings/check-availability` - Verificar disponibilidad

## Estructura del Proyecto

```Directory
booking-service/
├── main.go          # Punto de entrada y configuración del servidor
├── handlers.go      # Manejadores HTTP
├── models.go        # Estructuras de datos y DTOs
├── service.go       # Lógica de negocio central
├── repository.go    # Acceso a datos
├── Dockerfile       # Imagen Docker
├── go.mod          # Dependencias Go
└── README.md       # Documentación
```

## Estados de Reserva

- **PENDING**: Reserva creada, pendiente de confirmación
- **CONFIRMED**: Reserva confirmada y válida
- **CANCELLED**: Reserva cancelada
- **COMPLETED**: Reserva completada (después del tiempo de uso)

## Transiciones de Estado Válidas

```Transaction
PENDING → CONFIRMED
PENDING → CANCELLED
CONFIRMED → CANCELLED
CONFIRMED → COMPLETED
```

## Configuración

### Variables de Entorno

```bash
PORT=8003
DB_HOST=localhost
DB_PORT=5432
DB_NAME=reservas_bookings
DB_USER=postgres
DB_PASSWORD=password
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
USER_SERVICE_URL=http://user-service:8001
RESOURCE_SERVICE_URL=http://resource-service:8002
```

## Desarrollo Local

```bash
# Instalar dependencias
go mod tidy

# Ejecutar el servicio
go run .

# Construir imagen Docker
docker build -t booking-service .

# Ejecutar con Docker
docker run -p 8003:8003 booking-service
```

## Ejemplos de Uso

### Crear Reserva

```bash
curl -X POST http://localhost:8003/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "resource_id": 1,
    "start_time": "2025-06-10T14:00:00Z",
    "end_time": "2025-06-10T16:00:00Z",
    "notes": "Reunión de equipo"
  }'
```

### Listar Reservas por Usuario

```bash
curl "http://localhost:8003/api/v1/users/1/bookings?status=CONFIRMED&page=1&size=10"
```

### Verificar Disponibilidad

```bash
curl -X POST http://localhost:8003/api/v1/bookings/check-availability \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": 1,
    "start_time": "2025-06-10T14:00:00Z",
    "end_time": "2025-06-10T16:00:00Z"
  }'
```

### Confirmar Reserva

```bash
curl -X POST http://localhost:8003/api/v1/bookings/1/confirm \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Lógica de Negocio

### Validaciones

- No se pueden crear reservas en el pasado
- La hora de fin debe ser posterior a la hora de inicio
- No puede haber solapamiento de horarios para el mismo recurso
- Solo se pueden modificar reservas en estado PENDING o CONFIRMED

### Eventos Publicados

- `booking.created` - Nueva reserva creada
- `booking.updated` - Reserva modificada
- `booking.confirmed` - Reserva confirmada
- `booking.cancelled` - Reserva cancelada

## Integración con Otros Servicios

### User Service

- Validación de tokens JWT
- Obtención de información de usuario

### Resource Service  

- Verificación de existencia de recursos
- Consulta de disponibilidad de recursos

### Notification Service

- Envío de eventos para notificaciones automáticas

## Próximos Pasos

- [ ] Implementar base de datos PostgreSQL
- [ ] Integrar con User Service para autenticación
- [ ] Integrar con Resource Service para validación
- [ ] Implementar publicación de eventos en RabbitMQ
- [ ] Agregar validación de datos
- [ ] Implementar tests unitarios e integración
- [ ] Agregar logs estructurados
- [ ] Implementar métricas de negocio
