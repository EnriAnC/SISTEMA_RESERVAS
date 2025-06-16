# Resource Service

Microservicio encargado de la gestión de recursos disponibles para reservar en el sistema.

## Funcionalidades

- **Gestión de Recursos**: CRUD completo de recursos (salas, equipos, vehículos, etc.)
- **Tipos de Recursos**: Soporte para diferentes tipos de recursos
- **Disponibilidad**: Gestión de horarios de disponibilidad por recurso
- **Consultas**: Filtrado y búsqueda de recursos por tipo, ubicación y capacidad

## API Endpoints

### Recursos
- `POST /api/v1/resources` - Crear recurso
- `GET /api/v1/resources` - Listar recursos (con filtros)
- `GET /api/v1/resources/{id}` - Obtener recurso por ID
- `PUT /api/v1/resources/{id}` - Actualizar recurso
- `DELETE /api/v1/resources/{id}` - Eliminar recurso (soft delete)

### Disponibilidad
- `GET /api/v1/resources/{id}/availability` - Consultar disponibilidad
- `PUT /api/v1/resources/{id}/availability` - Actualizar horarios de disponibilidad

## Estructura del Proyecto

```
resource-service/
├── main.go          # Punto de entrada y configuración del servidor
├── handlers.go      # Manejadores HTTP
├── models.go        # Estructuras de datos y DTOs
├── service.go       # Lógica de negocio
├── repository.go    # Acceso a datos
├── Dockerfile       # Imagen Docker
├── go.mod          # Dependencias Go
└── README.md       # Documentación
```

## Tipos de Recursos Soportados

- **room**: Salas de reuniones, oficinas
- **equipment**: Equipos, proyectores, herramientas
- **vehicle**: Vehículos de la empresa
- **space**: Espacios comunes, estacionamientos

## Configuración

### Variables de Entorno
```bash
PORT=8002
DB_HOST=localhost
DB_PORT=5432
DB_NAME=reservas_resources
DB_USER=postgres
DB_PASSWORD=password
```

## Desarrollo Local

```bash
# Instalar dependencias
go mod tidy

# Ejecutar el servicio
go run .

# Construir imagen Docker
docker build -t resource-service .

# Ejecutar con Docker
docker run -p 8002:8002 resource-service
```

## Ejemplos de Uso

### Crear Recurso
```bash
curl -X POST http://localhost:8002/api/v1/resources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sala de Conferencias A",
    "type": "room",
    "description": "Sala con capacidad para 20 personas",
    "capacity": 20,
    "location": "Piso 3, Edificio Principal",
    "properties": {
      "projector": true,
      "whiteboard": true,
      "video_conference": true
    }
  }'
```

### Listar Recursos con Filtros
```bash
curl "http://localhost:8002/api/v1/resources?type=room&min_capacity=10&location=Piso%203"
```

### Configurar Disponibilidad
```bash
curl -X PUT http://localhost:8002/api/v1/resources/1/availability \
  -H "Content-Type: application/json" \
  -d '[
    {
      "day_of_week": 1,
      "start_time": "09:00",
      "end_time": "17:00"
    },
    {
      "day_of_week": 2,
      "start_time": "09:00",
      "end_time": "17:00"
    }
  ]'
```

### Consultar Disponibilidad
```bash
curl "http://localhost:8002/api/v1/resources/1/availability?start_date=2025-06-10&end_date=2025-06-12"
```

## Próximos Pasos

- [ ] Implementar base de datos PostgreSQL
- [ ] Agregar validación de datos
- [ ] Implementar cache para consultas frecuentes
- [ ] Agregar métricas de uso de recursos
- [ ] Implementar tests unitarios
- [ ] Agregar logs estructurados
- [ ] Optimizar consultas de disponibilidad
