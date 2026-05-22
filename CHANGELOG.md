# Changelog - Lector Comics

## v0.1.0 - 2026-05-21

### Added
- Estructura inicial del proyecto (monorepo)
- Documentos de especificaciones:
  - Lector_Comics_Master_Detallado.md (v0.1)
  - Lector_Comics_Especificacion_Tecnica_Empresarial.md (v0.2)
  - Lector_Comics_Implementacion_Ingenieria.md (v0.3)
- Plan de Implementación: Backend Go + Fiber, Frontend Vue + Nuxt, PostgreSQL, Redis, CBZ/CBR/ZIP

### Estructura Monorepo Creada
```
lector-comics/
├── backend/          # Go API
├── frontend/         # Vue/Nuxt
├── infra/            # Docker, nginx, configs
├── docs/             # Documentación
├── scripts/          # Scripts útiles
├── library/          # Archivos de comics
├── metadata/         # Metadata local
├── cache/            # Cache thumbnails
├── backups/          # Backups
└── logs/             # Logs
```

### Infraestructura Docker
- **docker-compose.yml**: 7 servicios (postgres, redis, backend, worker, frontend, nginx)
- Healthchecks configurados para todos los servicios
- Volúmenes: postgres_data, redis_data
- Redes internas aisladas
- **infra/docker/**: Dockerfiles para backend y worker
- **infra/nginx/nginx.conf**: Configuración reverse proxy con HTTP2 y gzip
- **infra/postgres/init.sql**: Schema completo con 12 tablas e índices

### Backend Go (Fiber)
- **go.mod**: Dependencias core (fiber, jwt, postgres, redis, bcrypt)
- **cmd/api/main.go**: Router principal con todos los endpoints v1
- **cmd/worker/main.go**: Worker service para procesamiento background
- **internal/config/**: Configuración via environment variables
- **internal/models/**: Modelos User, Library, Series, Volume, Chapter, ReadingProgress, Metadata, Collection
- **internal/repository/**: Conexión PostgreSQL y Redis client
- **internal/middleware/**: JWT auth middleware con claims
- **internal/services/**: Service layer completo con métodos de negocio
- **internal/handlers/**: Handlers completos para auth, library, series, reader, search, stats, user

### Scanner Engine
- **internal/scanner/scanner.go**: Parser para filenames CBZ/CBR/ZIP
  - Regex para detectar volumen, capítulo, special
  - Extracción de nombre de serie
  - Cálculo SHA256 hash
  - Conteo de páginas en archivos ZIP
- **internal/workers/scanner.go**: Worker de escaneo con auto-scan periódico
- **internal/services/scanner_service.go**: Servicio de escaneo con walker recursivo

### Reader Engine
- **internal/reader/reader.go**: Lector de comics CBZ/ZIP
  - Extracción de páginas desde archivos ZIP
  - Cache de lista de páginas
  - Extracción de datos de imagen

### Frontend Vue/Nuxt
- **package.json**: Nuxt 3 + Pinia
- **nuxt.config.ts**: Configuración app con runtime API
- **assets/css/main.css**: Tema oscuro base completo
- **layouts/default.vue**: Layout principal
- **pages/index.vue**: Login page funcional
- **pages/register.vue**: Registro de usuarios
- **pages/library.vue**: Biblioteca con grid de cómics y búsqueda
- **pages/reader.vue**: Lector con navegación teclado (A/D/F) y barra de progreso
- **components/AddLibraryModal.vue**: Modal para crear bibliotecas

### Scripts
- **scripts/init.sh**: Script de setup que:
  - Verifica prerequisitos (Docker)
  - Crea directorios necesarios
  - Build de imágenes Docker
  - Espera a que PostgreSQL esté listo
  - Inicia todos los servicios
  - Verifica health del backend

### Autenticación Implementada
- Register endpoint (POST /api/v1/auth/register)
- Login endpoint (POST /api/v1/auth/login) con bcrypt
- Logout endpoint (POST /api/v1/auth/logout)
- Refresh token endpoint (POST /api/v1/auth/refresh)
- JWT con expiration de 24h, refresh token de 7 días
- Middleware de protección para rutas autenticadas

### API Endpoints Implementados
- GET /health - Health check
- GET/POST /api/v1/library - Gestión de bibliotecas
- GET /api/v1/series - Lista de series
- GET /api/v1/reader/chapter/:id - Info del capítulo con progreso
- POST /api/v1/reader/progress - Guardar progreso de lectura
- GET /api/v1/search?q= - Búsqueda de series
- GET /api/v1/stats - Estadísticas de usuario
- GET /api/v1/user - Usuario actual

## v0.1.1 - 2026-05-21

### Added
- Handlers de autenticación completamente funcionales
  - Login con validación de password bcrypt
  - Register con hash de password automático
  - Generación de JWT tokens y refresh tokens
  - JWTAuth middleware mejorado con paso de jwt_secret por contexto

### Fixed
- Errores de compilación en services.go (tipos duplicados)
- Import de strings en reader.go (mal ubicado al final del archivo)
- Middleware JWT no pasaba jwt_secret a los handlers

### Modified
- Handler ahora recibe jwtSecret en el constructor
- Main.go pasa cfg.JWTSecret al handler
- Ruta protegida incluye middleware que pasa jwt_secret por contexto

## v0.1.2 - 2026-05-21

### Added
- Scanner real implementado:
  - Servicio de escaneo que recorre directorios recursivamente
  - Walker que procesa archivos CBZ/ZIP
  - Extracción de volumenes desde nombre de archivo (regex vol)
  - Hash de archivos para detección de duplicados
  - Conteo de páginas en archivos ZIP
  - Inserción en BD con ON CONFLICT para evitar duplicados

- Worker de scanner mejorado:
  - Escaneo automático cada 5 minutos
  - Solo bibliotecas con watch_enabled=true
  - Processamiento en background con WaitGroup
  - Stop graceful con contexto

### Files Added
- internal/workers/scanner.go - Worker de escaneo periódico
- internal/services/scanner_service.go - Servicio con lógica de walker

## v0.1.3 - 2026-05-21

### Added
- Script init.sh para setup rápido del proyecto

### Files Added
- scripts/init.sh - Script bash de inicialización

## v0.1.4 - 2026-05-21

### Added
- CHANGELOG.md registrado con todos los cambios del proyecto

### Modified
- CHANGELOG.md creado como historial completo del proyecto
- Incluye todas las versiones y cambios realizados

## v0.2.0 - 2026-05-21

### Added
- GetPage endpoint funcional para servir imágenes de comics
  - GET /api/v1/reader/chapter/:id/page/:page
  - Retorna imagen binaria con Content-Type: image/jpeg
  - Extrae imágenes de archivos CBZ/ZIP

- ComicReader integrado en handlers
  - Nuevo campo comicReader en struct Handler
  - Inicializado en NewHandler
  - GetPage usa comicReader.GetPage() para extraer imágenes

- ThumbnailService para generar portadas
  - internal/cache/thumbnail.go
  - Genera thumbnails de series desde primera página del comic
  - Cache en Redis con TTL de 30 días
  - Guarda en disco en /app/cache/

### Modified
- GetPage handler ahora:
  - Obtiene chapter de BD para obtener file_path
  - Usa ComicReader para extraer imagen
  - Retorna imagen binaria directamente
  - Mejor manejo de errores (404 para páginas no encontradas)

- Reader Vue component actualizado:
  - totalPages ahora es ref en lugar de computed
  - Carga páginas desde API (/page/:page)
  - Limpia URLs de objetos anteriores para evitar memory leaks
  - Campos de API alineados con backend (page_count, current_page)
  - chapter_id en saveProgress en lugar de chapterId

### Fixed
- Import de strconv faltante en handlers.go
- chapterId como string simple en lugar de computed (evita problemas de reactivity)

## v0.2.1 - 2026-05-21

### Added
- Frontend library.vue mejorado
  - Carga tanto bibliotecas como series en paralelo (Promise.all)
  - Asocia series con sus bibliotecas对应的library
  - openSeries ahora abre el primer capítulo disponible de la serie

### Modified
- fetchLibrary ahora hace dos llamadas:
  1. GET /api/v1/library - lista de bibliotecas
  2. GET /api/v1/series - lista de series
- Asocia cada serie con su biblioteca对应的library.id

## v0.2.2 - 2026-05-22

### Added
- Archivos go.mod y go.sum faltantes para el backend Go
  - Dependencias completas incluyendo fiber, jwt, postgres, redis, bcrypt
  -go.sum con checksums de todas las dependencias

- Archivos .dockerignore para backend y frontend
  - backend/.dockerignore: excluye .git, .md, docs, scripts, tests
  - frontend/.dockerignore: excluye .git, .md, docs, scripts, tests, node_modules, .nuxt

- Guia de instalacion completa (instalacion.txt)
  - Requisitos previos (Docker, Docker Compose)
  - Estructura del proyecto
  - Permisos y carpetas
  - Inicio de servicios (rapido y manual)
  - Verificacion de instalacion
  - Acceso a la aplicacion
  - Primeros pasos
  - Comandos utiles
  - Solucion de problemas
  - Actualizaciones
  - Desinstalacion

### Fixed
- Error de build en Docker: "go.sum not found"
  - Agregado go.sum faltante
  - Creados .dockerignore para evitar que archivos incorrectos se copien al contexto

## v0.2.3 - 2026-05-22

### Fixed
- Error "checksum mismatch" en go.mod/go.sum
  - El hash de github.com/golang-jwt/jwt/v5@v5.2.0 era incorrecto
  - go.sum eliminado - ahora se genera automaticamente con `go mod tidy`

### Modified
- **Dockerfile** y **Dockerfile.worker** actualizados:
  - `COPY go.mod go.sum ./` → `COPY go.mod ./` y `RUN go mod tidy`
  - Ya no depende de go.sum pre-generado
  - `go mod tidy` genera go.sum correctamente dentro del contenedor
  - `go mod download` eliminado ya que `go mod tidy` lo incluye

### Files Changed
- backend/Dockerfile: Removida dependencia de go.sum externo
- backend/Dockerfile.worker: Removida dependencia de go.sum externo

## v0.2.4 - 2026-05-22

### Fixed
- Error "missing go.sum entry for module providing package"
  - `go mod tidy` solo actualiza go.sum pero no descarga los módulos
  - Agregado `go mod download` después de `go mod tidy`

### Modified
- **Dockerfile** y **Dockerfile.worker** actualizados:
  - `RUN go mod tidy` → `RUN go mod tidy && go mod download`
  - Ambos comandos se ejecutan en una sola capa RUN