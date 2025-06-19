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

### 3. Security Analysis Deep Scan (`security.yml`)

**Propósito:** Análisis de seguridad profundo y reportes consolidados

- ✅ Análisis de contenedores Docker (`Trivy`)
- ✅ Infraestructura como código (`Checkov`)
- ✅ Detección avanzada de secretos (`TruffleHog`)
- ✅ Análisis de cadena de suministro
- ✅ Reportes consolidados con recomendaciones
- ✅ Comentarios automáticos en PRs

**Triggers:** Push a main (cambios críticos), schedule semanal y ejecución manual

**Nota:** El análisis básico de seguridad (gosec, govulncheck) se ejecuta en CI para feedback inmediato

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

1. **Análisis básico (CI):** Revisar alertas de `gosec` y `govulncheck`
2. **Análisis profundo (Security):** Verificar reportes de Trivy, Checkov y TruffleHog
3. Actualizar dependencias vulnerables
4. Corregir configuraciones inseguras identificadas
5. Rotar secretos comprometidos identificados
6. Revisar GitHub Security tab para alertas adicionales

## 📈 Mejoras Futuras

Cuando el proyecto madure, se pueden agregar:

- Tests unitarios e integración
- Construcción y push de imágenes Docker
- Despliegue automático a staging
- Análisis de cobertura de código
- Notificaciones de Slack/Teams
- Integración con herramientas de monitoreo
- Políticas de seguridad automatizadas (OPA/Gatekeeper)
- Análisis de performance y carga

---

*Esta configuración de GitHub Actions está diseñada para evolucionar con el proyecto,  
manteniendo alta calidad de código desde el inicio.*
