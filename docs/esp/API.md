# Documentación de API - Sistema de Reservas

## URL Base

```URL
http://localhost:8080/api/v1
```

## Autenticación

La mayoría de endpoints requieren autenticación usando tokens JWT. Incluye el token en el header de Autorización:

```Authorization
Authorization: Bearer <tu-token-jwt>
```

## Respuestas de Error

Todos los errores siguen un formato consistente:

```json
{
  "error": "Mensaje de error",
  "code": "CODIGO_ERROR",
  "timestamp": "2024-06-09T10:30:00Z"
}
```

Códigos de estado HTTP comunes:

- `200` - Éxito
- `201` - Creado
- `400` - Solicitud Incorrecta
- `401` - No Autorizado
- `403` - Prohibido
- `404` - No Encontrado
- `409` - Conflicto
- `500` - Error Interno del Servidor

---

## Gestión de Usuarios

### Autenticarse

#### Iniciar Sesión

```http
POST /auth/login
Content-Type: application/json

{
  "email": "usuario@ejemplo.com",
  "password": "password123"
}
```

**Respuesta:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "usuario@ejemplo.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "role": "user"
  }
}
```

### Operaciones de Usuario

#### Crear Usuario

```http
POST /users
Content-Type: application/json

{
  "email": "nuevousuario@ejemplo.com",
  "password": "password123",
  "first_name": "Ana",
  "last_name": "García",
  "phone": "+34123456789",
  "role": "user"
}
```

#### Obtener Todos los Usuarios

```http
GET /users?limit=10&offset=0&role=user
Authorization: Bearer <token>
```

#### Obtener Usuario por ID

```http
GET /users/{id}
Authorization: Bearer <token>
```

#### Actualizar Usuario

```http
PUT /users/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Nombre Actualizado",
  "phone": "+34987654321"
}
```

#### Eliminar Usuario

```http
DELETE /users/{id}
Authorization: Bearer <token>
```

---

## Gestión de Recursos

### Operaciones de Recursos

#### Obtener Todos los Recursos

```http
GET /resources?type=meeting_room&location=Planta1&available=true&limit=10&offset=0
```

**Parámetros de Consulta:**

- `type` - Filtrar por tipo de recurso
- `location` - Filtrar por ubicación
- `available` - Filtrar por disponibilidad
- `capacity_min` - Capacidad mínima
- `capacity_max` - Capacidad máxima
- `limit` - Número de resultados (por defecto: 10)
- `offset` - Desplazamiento de paginación (por defecto: 0)

**Respuesta:**

```json
{
  "resources": [
    {
      "id": 1,
      "name": "Sala de Conferencias A",
      "description": "Sala de conferencias grande con proyector",
      "type": "meeting_room",
      "location": "Planta 1",
      "capacity": 10,
      "price_per_hour": 50.00,
      "amenities": {
        "projector": true,
        "whiteboard": true,
        "video_conference": true
      },
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 25,
  "limit": 10,
  "offset": 0
}
```

#### Crear Recurso

```http
POST /resources
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Nueva Sala de Reuniones",
  "description": "Espacio pequeño para reuniones",
  "type": "meeting_room",
  "location": "Planta 2",
  "capacity": 6,
  "price_per_hour": 30.00,
  "amenities": {
    "whiteboard": true,
    "phone": true
  }
}
```

#### Obtener Recurso por ID

```http
GET /resources/{id}
```

#### Actualizar Recurso

```http
PUT /resources/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Nombre de Sala Actualizado",
  "capacity": 8,
  "price_per_hour": 35.00
}
```

#### Eliminar Recurso

```http
DELETE /resources/{id}
Authorization: Bearer <token>
```

#### Verificar Disponibilidad del Recurso

```http
GET /resources/{id}/availability?date=2024-06-10&start_time=09:00&end_time=17:00
```

**Respuesta:**

```json
{
  "resource_id": 1,
  "date": "2024-06-10",
  "availability": [
    {
      "start_time": "09:00",
      "end_time": "10:00",
      "available": true
    },
    {
      "start_time": "10:00",
      "end_time": "12:00",
      "available": false,
      "booking_id": 123
    }
  ]
}
```

---

## Gestión de Reservas

### Operaciones de Reservas

#### Obtener Todas las Reservas

```http
GET /bookings?user_id=123&resource_id=456&status=confirmed&start_date=2024-06-01&end_date=2024-06-30&limit=10&offset=0
Authorization: Bearer <token>
```

**Parámetros de Consulta:**

- `user_id` - Filtrar por usuario
- `resource_id` - Filtrar por recurso
- `status` - Filtrar por estado (pending, confirmed, cancelled, completed)
- `start_date` - Filtrar reservas desde fecha
- `end_date` - Filtrar reservas hasta fecha
- `limit` - Número de resultados
- `offset` - Desplazamiento de paginación

#### Crear Reserva

```http
POST /bookings
Authorization: Bearer <token>
Content-Type: application/json

{
  "resource_id": 1,
  "start_time": "2024-06-10T10:00:00Z",
  "end_time": "2024-06-10T12:00:00Z",
  "notes": "Reunión de equipo para planificación del proyecto"
}
```

**Respuesta:**

```json
{
  "id": 123,
  "user_id": 456,
  "resource_id": 1,
  "start_time": "2024-06-10T10:00:00Z",
  "end_time": "2024-06-10T12:00:00Z",
  "status": "pending",
  "total_price": 100.00,
  "notes": "Reunión de equipo para planificación del proyecto",
  "created_at": "2024-06-09T15:30:00Z"
}
```

#### Obtener Reserva por ID

```http
GET /bookings/{id}
Authorization: Bearer <token>
```

#### Actualizar Reserva

```http
PUT /bookings/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "start_time": "2024-06-10T11:00:00Z",
  "end_time": "2024-06-10T13:00:00Z",
  "notes": "Horario de reunión actualizado"
}
```

#### Cancelar Reserva

```http
PUT /bookings/{id}/cancel
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "Reunión pospuesta"
}
```

#### Eliminar Reserva

```http
DELETE /bookings/{id}
Authorization: Bearer <token>
```

---

## Gestión de Notificaciones

### Operaciones de Notificaciones

#### Enviar Notificación

```http
POST /notifications
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 123,
  "type": "booking",
  "title": "Reserva Confirmada",
  "message": "Su reserva para la Sala de Conferencias A ha sido confirmada",
  "channel": "email",
  "priority": "high",
  "metadata": {
    "booking_id": 456,
    "resource_name": "Sala de Conferencias A"
  }
}
```

#### Obtener Notificaciones del Usuario

```http
GET /notifications?user_id=123&type=booking&is_read=false&limit=10&offset=0
Authorization: Bearer <token>
```

**Respuesta:**

```json
{
  "notifications": [
    {
      "id": 1,
      "user_id": 123,
      "type": "booking",
      "title": "Reserva Confirmada",
      "message": "Su reserva ha sido confirmada",
      "channel": "email",
      "priority": "high",
      "is_read": false,
      "status": "sent",
      "created_at": "2024-06-09T10:30:00Z",
      "sent_at": "2024-06-09T10:30:05Z"
    }
  ],
  "total": 25,
  "limit": 10,
  "offset": 0
}
```

#### Actualizar Estado de Notificación

```http
PUT /notifications/{id}/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_read": true
}
```

#### Obtener Estadísticas de Notificaciones

```http
GET /notifications/stats?user_id=123
Authorization: Bearer <token>
```

**Respuesta:**

```json
{
  "user_id": 123,
  "total": 50,
  "read": 30,
  "unread": 20,
  "high": 10,
  "normal": 35,
  "low": 5,
  "sent": 45,
  "failed": 2,
  "pending": 3
}
```

---

## Endpoints del Sistema

### Verificación de Salud

```http
GET /health
```

**Respuesta:**

```json
{
  "status": "healthy",
  "services": {
    "user-service": "healthy",
    "resource-service": "healthy",
    "booking-service": "healthy",
    "notification-service": "healthy"
  },
  "timestamp": "2024-06-09T10:30:00Z"
}
```

---

## Modelos de Datos

### Modelo de Usuario

```json
{
  "id": 1,
  "email": "usuario@ejemplo.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "phone": "+34123456789",
  "role": "user",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "last_login": "2024-06-09T08:00:00Z"
}
```

### Modelo de Recurso

```json
{
  "id": 1,
  "name": "Sala de Conferencias A",
  "description": "Sala de conferencias grande con proyector",
  "type": "meeting_room",
  "location": "Planta 1",
  "capacity": 10,
  "price_per_hour": 50.00,
  "amenities": {
    "projector": true,
    "whiteboard": true,
    "video_conference": true
  },
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Modelo de Reserva

```json
{
  "id": 123,
  "user_id": 456,
  "resource_id": 1,
  "start_time": "2024-06-10T10:00:00Z",
  "end_time": "2024-06-10T12:00:00Z",
  "status": "confirmed",
  "total_price": 100.00,
  "notes": "Reunión de equipo",
  "metadata": {},
  "created_at": "2024-06-09T15:30:00Z",
  "updated_at": "2024-06-09T15:30:00Z",
  "cancelled_at": null,
  "cancellation_reason": null
}
```

### Modelo de Notificación

```json
{
  "id": 1,
  "user_id": 123,
  "type": "booking",
  "title": "Reserva Confirmada",
  "message": "Su reserva ha sido confirmada",
  "channel": "email",
  "priority": "high",
  "status": "sent",
  "is_read": false,
  "metadata": {
    "booking_id": 456
  },
  "created_at": "2024-06-09T10:30:00Z",
  "sent_at": "2024-06-09T10:30:05Z",
  "read_at": null
}
```

---

## Eventos de Webhook

El sistema publica eventos para integración externa:

### Eventos de Reserva

- `booking.created` - Nueva reserva creada
- `booking.confirmed` - Reserva confirmada
- `booking.cancelled` - Reserva cancelada
- `booking.updated` - Reserva modificada

### Eventos de Usuario

- `user.registered` - Nuevo usuario registrado
- `user.updated` - Información de usuario actualizada

### Ejemplo de Payload de Evento

```json
{
  "event": "booking.confirmed",
  "timestamp": "2024-06-09T10:30:00Z",
  "data": {
    "booking_id": 123,
    "user_id": 456,
    "resource_id": 1,
    "start_time": "2024-06-10T10:00:00Z",
    "end_time": "2024-06-10T12:00:00Z"
  }
}
```

---

## Limitación de Velocidad

Los endpoints de la API tienen limitación de velocidad para prevenir abuso:

- **Endpoints de autenticación**: 5 solicitudes por minuto por IP
- **Endpoints de usuario**: 100 solicitudes por hora por usuario
- **Endpoints de recurso**: 200 solicitudes por hora por usuario
- **Endpoints de reserva**: 50 solicitudes por hora por usuario
- **Endpoints de notificación**: 100 solicitudes por hora por usuario

Los headers de límite de velocidad se incluyen en las respuestas:

```Headers
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1623456789
```

---

## SDK y Ejemplos

### Ejemplo JavaScript/Node.js

```javascript
const axios = require('axios');

// Iniciar sesión
const loginResponse = await axios.post('http://localhost:8080/api/v1/auth/login', {
  email: 'usuario@ejemplo.com',
  password: 'password123'
});

const token = loginResponse.data.token;

// Crear reserva
const bookingResponse = await axios.post(
  'http://localhost:8080/api/v1/bookings',
  {
    resource_id: 1,
    start_time: '2024-06-10T10:00:00Z',
    end_time: '2024-06-10T12:00:00Z',
    notes: 'Reunión de equipo'
  },
  {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  }
);
```

### Ejemplos cURL

```bash
# Iniciar sesión
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@ejemplo.com","password":"password123"}'

# Obtener recursos
curl -X GET "http://localhost:8080/api/v1/resources?type=meeting_room" \
  -H "Authorization: Bearer TU_TOKEN"

# Crear reserva
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Authorization: Bearer TU_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": 1,
    "start_time": "2024-06-10T10:00:00Z",
    "end_time": "2024-06-10T12:00:00Z",
    "notes": "Reunión de equipo"
  }'
```
