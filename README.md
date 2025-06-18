# Sistema de Reservas en la Nube

<!-- Badges de Estado del Proyecto -->
[![CI Pipeline](https://github.com/enrianc/sistema_reservas/actions/workflows/ci.yml/badge.svg)](https://github.com/enrianc/sistema_reservas/actions/workflows/ci.yml)
[![Validation Pipeline](https://github.com/enrianc/sistema_reservas/actions/workflows/validate.yml/badge.svg)](https://github.com/enrianc/sistema_reservas/actions/workflows/validate.yml)
[![Security Analysis](https://github.com/enrianc/sistema_reservas/actions/workflows/security.yml/badge.svg)](https://github.com/enrianc/sistema_reservas/actions/workflows/security.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5?style=flat&logo=kubernetes)](https://kubernetes.io)

## Descripción

Sistema de reservas escalable basado en microservicios desarrollado en Go, diseñado para manejar reservas de recursos con alta disponibilidad y escalabilidad automática.

## Arquitectura

- **Patrón:** Microservicios
- **Lenguaje:** Go/Golang
- **Contenerización:** Docker
- **Orquestación:** Kubernetes
- **API Gateway:** KrakenD
- **Base de Datos:** PostgreSQL
- **Mensajería:** RabbitMQ

## Microservicios

1. **User Service** - Gestión de usuarios y autenticación
2. **Resource Service** - Gestión de recursos disponibles
3. **Booking Service** - Lógica principal de reservas
4. **Notification Service** - Sistema de notificaciones

## Estructura del Proyecto

```Directory
SISTEMA_RESERVAS/
├── services/               # Microservicios
│   ├── user-service/
│   ├── resource-service/
│   ├── booking-service/
│   └── notification-service/
├── infrastructure/         # Configuración de infraestructura
├── api-gateway/           # Configuración del API Gateway
├── docs/                  # Documentación técnica
├── scripts/               # Scripts de deployment
└── kubernetes/            # Manifiestos de Kubernetes
```

## Inicio Rápido

```bash
# Clonar el repositorio
git clone <repository-url>
cd SISTEMA_RESERVAS

# Copiar configuración de ejemplo
cp .env.example .env

# Levantar servicios localmente
docker-compose up -d

# Verificar estado de los servicios
curl http://localhost:8080/users/health
```

📖 **[Ver Guía Completa de Inicio Rápido](QUICKSTART.md)**

## Documentación

- [Arquitectura Técnica](docs/ARCHITECTURE.md)
- [API Documentation](docs/API.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [Development Guide](docs/DEVELOPMENT.md)

## Estado del Proyecto

✅ **Completado** - Sistema completo con arquitectura de microservicios funcional

### Características Implementadas

- ✅ 4 Microservicios completamente funcionales (User, Resource, Booking, Notification)
- ✅ API Gateway con KrakenD configurado
- ✅ Base de datos PostgreSQL con esquema completo
- ✅ Configuración Docker Compose para desarrollo local
- ✅ Manifiestos Kubernetes para producción
- ✅ Monitoreo con Prometheus y Grafana
- ✅ Scripts de automatización de despliegue
- ✅ Documentación técnica completa
- ✅ Guía de desarrollo para contribuidores
