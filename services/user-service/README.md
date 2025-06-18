# User Service

Microservicio encargado de la gestión de usuarios y autenticación del sistema de reservas.

## Funcionalidades

- **Gestión de Usuarios**: CRUD completo de usuarios
- **Autenticación**: Login con JWT tokens
- **Autorización**: Gestión de roles (admin, user)
- **Seguridad**: Hash de contraseñas con bcrypt

## API Endpoints

### Usuarios

- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users/{id}` - Obtener usuario por ID
- `PUT /api/v1/users/{id}` - Actualizar usuario
- `DELETE /api/v1/users/{id}` - Eliminar usuario

### Autenticación

- `POST /api/v1/auth/login` - Iniciar sesión
- `POST /api/v1/auth/refresh` - Refrescar token

## Estructura del Proyecto

```Directory
user-service/
├── main.go          # Punto de entrada y configuración del servidor
├── handlers.go      # Manejadores HTTP
├── models.go        # Estructuras de datos y DTOs
├── service.go       # Lógica de negocio
├── repository.go    # Acceso a datos
├── Dockerfile       # Imagen Docker
├── go.mod          # Dependencias Go
└── README.md       # Documentación
```

## Configuración

### Variables de Entorno

```bash
PORT=8001
JWT_SECRET=your-secret-key
DB_HOST=localhost
DB_PORT=5432
DB_NAME=reservas_users
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
docker build -t user-service .

# Ejecutar con Docker
docker run -p 8001:8001 user-service
```

## Ejemplos de Uso

### Crear Usuario

```bash
curl -X POST http://localhost:8001/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "password": "password123",
    "role": "user"
  }'
```

### Login

```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan@example.com",
    "password": "password123"
  }'
```

## Próximos Pasos

- [ ] Implementar base de datos PostgreSQL
- [ ] Agregar validación de datos
- [ ] Implementar JWT real
- [ ] Agregar middleware de autenticación
- [ ] Implementar tests unitarios
- [ ] Agregar logs estructurados
- [ ] Implementar métricas
