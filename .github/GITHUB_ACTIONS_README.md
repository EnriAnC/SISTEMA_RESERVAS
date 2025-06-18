# Configuración de GitHub Actions

Este proyecto utiliza GitHub Actions para automatizar la validación, construcción y análisis de código.

## 🚀 Workflows Configurados

### 1. CI Pipeline (`ci.yml`)

**Propósito:** Validación continua de código Go

- ✅ Compilación de todos los microservicios
- ✅ Verificación de formato (`gofmt`)
- ✅ Análisis estático (`go vet`)
- ✅ Linting avanzado (`golangci-lint`)
- ✅ Verificación de dependencias (`go mod tidy`)
- ✅ Validación de estructura del proyecto

**Triggers:** Push y Pull Request a `main` y `develop`

### 2. Validation Pipeline (`validate.yml`)

**Propósito:** Validación de archivos de configuración

- ✅ Sintaxis de Dockerfiles (`hadolint`)
- ✅ Archivos YAML (Kubernetes, Docker Compose)
- ✅ Archivos Markdown (`markdownlint`)
- ✅ Archivos JSON
- ✅ Scripts Bash (`shellcheck`)

**Triggers:** Push y Pull Request a `main` y `develop`

### 3. Security Analysis (`security.yml`)

**Propósito:** Análisis de seguridad del código

- ✅ Vulnerabilidades en código Go (`gosec`)
- ✅ Dependencias vulnerables (`govulncheck`)
- ✅ Detección de secretos (`truffleHog`)
- ✅ Configuraciones de seguridad
- ✅ Reporte de seguridad

**Triggers:** Push, Pull Request y schedule semanal

## 📊 Badges de Estado

Los badges en el README.md muestran el estado actual de los workflows:

```markdown
[![CI Pipeline](https://github.com/tu-usuario/sistema-reservas/actions/workflows/ci.yml/badge.svg)](https://github.com/tu-usuario/sistema-reservas/actions/workflows/ci.yml)
```

## ⚙️ Configuración Local

### Prerequisitos para desarrollo local

```bash
# Instalar herramientas de calidad de código
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

# Instalar herramientas de validación
npm install -g markdownlint-cli
sudo apt-get install yamllint shellcheck
```

### Ejecutar validaciones localmente

```bash
# Lint de código Go
golangci-lint run ./...

# Verificar formato
gofmt -l .

# Análisis de seguridad
gosec ./...

# Verificar vulnerabilidades
govulncheck ./...
```

## 🔧 Configuración de Herramientas

### golangci-lint

Configuración en `.golangci.yml`:

- Linters habilitados para calidad de código
- Reglas específicas para microservicios
- Configuración de complejidad y estilo

### Validación de archivos

- **Dockerfiles:** Validados con `hadolint`
- **YAML:** Validados con `yamllint`
- **Markdown:** Validados con `markdownlint`
- **Scripts:** Validados con `shellcheck`

## 🚨 Resolución de Problemas

### CI Pipeline falla

1. Verificar que el código compila localmente
2. Ejecutar `gofmt -w .` para formatear código
3. Ejecutar `go mod tidy` para limpiar dependencias
4. Revisar errores de `golangci-lint`

### Validation Pipeline falla

1. Verificar sintaxis de Dockerfiles
2. Validar archivos YAML con `yamllint`
3. Corregir formato de Markdown
4. Revisar scripts con `shellcheck`

### Security Pipeline falla

1. Revisar alertas de `gosec`
2. Actualizar dependencias vulnerables
3. Remover secretos hardcodeados
4. Revisar configuraciones de seguridad

## 📈 Mejoras Futuras

Cuando el proyecto madure, se pueden agregar:

- Tests unitarios e integración
- Construcción y push de imágenes Docker
- Despliegue automático a staging
- Análisis de cobertura de código
- Notificaciones de Slack/Teams
- Integración con herramientas de monitoreo

---

*Esta configuración de GitHub Actions está diseñada para evolucionar con el proyecto, manteniendo alta calidad de código desde el inicio.*
