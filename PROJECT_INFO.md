# Lector Comics - Project Info

Version: 0.2.15
Last Updated: 2026-05-22

---

## 1. ARQUITECTURA

### Stack Tecnológico
- **Backend**: Go 1.22 + Fiber v2.52
- **Frontend**: Nuxt 3 + Vue 3 + Pinia
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Proxy**: nginx Alpine
- **Runtime**: Docker Compose

### Diagrama de Servicios
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│   nginx    │────▶│  frontend   │
│  (Browser) │     │  (port 80)  │     │ (Nuxt :3000) │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │ /api/
                    ┌──────▼──────┐
                    │   backend   │
                    │ (Fiber :3000)│
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              ▼                         ▼
       ┌─────────────┐          ┌─────────────┐
       │  postgres   │          │    redis    │
       │  (port 5432)│          │ (port 6379) │
       └─────────────┘          └─────────────┘
              │
       ┌──────▼──────┐
       │    worker   │
       │ (background)│
       └─────────────┘
```

### Puertos de Acceso
| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| nginx | 80 | Proxy reverso, punto de acceso principal |
| backend | 3000 | API REST |
| postgres | 5432 | Base de datos |
| redis | 6379 | Cache y pub/sub |

---

## 2. ESTRUCTURA DEL PROYECTO

```
lector-comics/
├── backend/
│   ├── cmd/
│   │   ├── api/main.go          # API server entry point
│   │   └── worker/main.go      # Background worker entry point
│   ├── internal/
│   │   ├── cache/thumbnail.go   # Thumbnail generation
│   │   ├── config/config.go     # Environment config
│   │   ├── handlers/handlers.go # HTTP handlers
│   │   ├── middleware/jwt.go    # JWT authentication
│   │   ├── models/models.go     # Data models
│   │   ├── reader/reader.go      # CBZ/ZIP reader
│   │   ├── repository/repository.go # DB & Redis access
│   │   ├── scanner/scanner.go   # Comic file parser
│   │   ├── services/services.go # Business logic
│   │   └── workers/
│   │       ├── scanner.go       # Library scanner worker
│   │       ├── thumbnail.go     # Thumbnail worker
│   │       └── cleanup.go       # Cleanup worker
│   ├── Dockerfile
│   ├── Dockerfile.worker
│   └── go.mod
├── frontend/
│   ├── pages/
│   │   ├── index.vue           # Login
│   │   ├── register.vue        # Register
│   │   ├── library.vue         # Main library view
│   │   └── reader.vue          # Comic reader
│   ├── components/
│   │   ├── AddLibraryModal.vue
│   │   ├── EditLibraryModal.vue
│   │   └── DeleteLibraryModal.vue
│   ├── layouts/default.vue
│   ├── assets/css/main.css
│   ├── nuxt.config.ts
│   ├── package.json
│   └── Dockerfile
├── infra/
│   ├── docker/
│   ├── nginx/nginx.conf        # Reverse proxy config
│   └── postgres/init.sql       # Database schema
├── library/                     # Comics storage
├── metadata/                    # Local metadata
├── cache/                      # Thumbnail cache
├── logs/                       # Application logs
├── backups/                    # Backups
├── scripts/
│   └── init.sh                 # Setup script
├── docker-compose.yml
├── CHANGELOG.md
├── instalacion.txt             # Installation guide (Spanish)
├── PROJECT_INFO.md             # This file
└── docs/                       # Specification documents
    ├── Lector_Comics_Master_Detallado.md
    ├── Lector_Comics_Especificacion_Tecnica_Empresarial.md
    └── Lector_Comics_Implementacion_Ingenieria.md
```

---

## 3. API ENDPOINTS

### Autenticación (públicos)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | /api/v1/auth/register | Registro de usuario |
| POST | /api/v1/auth/login | Inicio de sesión |
| POST | /api/v1/auth/logout | Cierre de sesión |
| POST | /api/v1/auth/refresh | Refrescar token |

### Bibliotecas (requiere auth)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/v1/library | Listar bibliotecas |
| POST | /api/v1/library | Crear biblioteca |
| PATCH | /api/v1/library/:id | Actualizar biblioteca |
| DELETE | /api/v1/library/:id | Eliminar biblioteca |
| POST | /api/v1/library/:id/scan | Escanear biblioteca (async) |

### Series (requiere auth)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/v1/series | Listar series |
| GET | /api/v1/series/:id | Obtener serie + volumes + chapters |

### Scan Jobs (requiere auth)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/v1/scans/:scanId | Consultar estado de scan |

### Reader (requiere auth)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/v1/reader/chapter/:id | Info del capítulo |
| GET | /api/v1/reader/chapter/:id/page/:page | Imagen de página |
| POST | /api/v1/reader/progress | Guardar progreso |

### Otros (requiere auth)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/v1/search?q= | Buscar series |
| GET | /api/v1/stats | Estadísticas de usuario |
| GET | /api/v1/user | Usuario actual |
| GET | /health | Health check (público) |

---

## 4. VARIABLES DE ENTORNO

### Backend/Worker
| Variable | Descripción | Valor ejemplo |
|----------|-------------|----------------|
| DATABASE_URL | Connection string PostgreSQL | `postgres://lector:lector_secret@postgres:5432/lector?sslmode=disable` |
| REDIS_URL | Connection string Redis | `redis://redis:6379` |
| JWT_SECRET | Secret para tokens JWT | `change_this_in_production` |
| PORT | Puerto del API server | `3000` |

### Frontend
| Variable | Descripción | Valor ejemplo |
|----------|-------------|----------------|
| API_URL | URL del backend API | `http://backend:3000` |
| HOST | Host para Nuxt dev server | `0.0.0.0` |
| PORT | Puerto para Nuxt dev server | `3000` |

---

## 5. SCHEMA DE BASE DE DATOS

Ver archivo: `infra/postgres/init.sql`

Tablas principales:
- `users` - Usuarios del sistema
- `library` - Bibliotecas de cómics
- `series` - Series dentro de bibliotecas
- `volume` - Volúmenes de series
- `chapter` - Capítulos (archivos CBZ/ZIP)
- `reading_progress` - Progreso de lectura por usuario
- `scan_job` - Jobs de escaneo asíncrono
- `collection` - Colecciones de usuario
- `collection_item` - Items en colecciones
- `metadata` - Metadata externa de series

---

## 6. ESTADO DEL PROYECTO

### Implementado ✅
- Sistema de autenticación (JWT con bcrypt)
- CRUD completo de bibliotecas
- Escaneo asíncrono de cómics (CBZ/ZIP) via Redis pub/sub
- Lector de cómics con navegación por teclado
- Guardado de progreso de lectura
- Búsqueda de series
- Estadísticas básicas
- Thumbnails de portadas con cache en Redis

### Pendiente ⏳
- Metadata externo (cover art, descripciones desde APIs externas)
- Favoritos/Want To Read
- Actualización de perfil de usuario (UpdateUser es stub)
- Zoom en reader
- Modo Webtoon vertical scroll
- Marcadores/Bookmarks
- PDF/EPUB reader
- SSL/TLS en producción
- Rate limiting
- Backup automatizado

---

## 7. DEPLOYMENT

### Requisitos
- Docker 20.10+
- Docker Compose 2.0+
- 2GB RAM mínimo
- 10GB disco

### Comandos de Deploy
```bash
# 1. Clonar/copiar proyecto
git clone <repo> ~/lector/Lector
cd ~/lector/Lector

# 2. Crear carpetas necesarias
mkdir -p library metadata cache backups logs

# 3. Rebuild e iniciar
docker compose build --no-cache
docker compose up -d

# 4. Verificar
docker compose ps
curl http://localhost/health

# 5. Acceder
http://<tu-servidor>:80
```

### Acceso desde red externa
- Punto de entrada: `http://<ip-servidor>:80`
- **NO usar** puertos 3000 ni 8080 directamente
- Todo el tráfico pasa por nginx (proxy reverso)

---

## 8. LIMITACIONES CONOCIDAS

1. **Frontend en modo dev**: `npm run dev` en Docker
   - NO es óptimo para producción
   - Solución futura: Build completo con `nuxt build`

2. **JWT_SECRET hardcoded**: `change_this_in_production`
   - **OBLIGATORIO** cambiar antes de producción

3. **Sin SSL/TLS**: nginx sin certificados
   - Añadir Let's Encrypt para producción

4. **Scanner solo CBZ/ZIP**: No soporta CBR (RAR)
   - Formatos soportados: .cbz, .cbr, .zip

---

## 9. FLUJO DE ESCANEO ASÍNCRONO

```
1. User pulsa "Escanear" en frontend
       ↓
2. Frontend: POST /api/v1/library/{id}/scan
       ↓
3. API: Crea scan_job en BD (status=queued)
       ↓
4. API: Publica en Redis canal "lector:scan_library"
       ↓
5. API: Retorna {scan_id: "uuid-xxx"}
       ↓
6. Frontend: Inicia polling GET /api/v1/scans/{scan_id} cada 2s
       ↓
7. Worker: Subscribe a Redis channel
       ↓
8. Worker: Recibe mensaje, actualiza status="scanning"
       ↓
9. Worker: Escanea biblioteca, inserta chapters en BD
       ↓
10. Worker: Actualiza status="completed" o "failed"
       ↓
11. Frontend: Finaliza polling, muestra resultado
```

### Tabla scan_job
```sql
CREATE TABLE scan_job (
    id UUID PRIMARY KEY,
    library_id BIGINT REFERENCES library(id),
    user_id BIGINT REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'queued', -- queued, scanning, completed, failed
    chapters_found INT DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 10. ENDPOINTS STUBS (Pendientes de Implementar)

| Handler | Ruta | Estado |
|---------|------|--------|
| UpdateUser | PATCH /api/v1/user | Stub "todo" |
| GetSeriesByID | GET /api/v1/series/:id | ✅ Implementado v0.2.14 |
| UpdateLibrary | PATCH /api/v1/library/:id | ✅ Implementado v0.2.14 |
| DeleteLibrary | DELETE /api/v1/library/:id | ✅ Implementado v0.2.14 |
| ScanLibrary | POST /api/v1/library/:id/scan | ✅ Implementado v0.2.14 |
| GetScanStatus | GET /api/v1/scans/:scanId | ✅ Implementado v0.2.14 |
| GetMetadata | GET /api/v1/metadata/:seriesId | Stub "todo" |
| RefreshMetadata | POST /api/v1/metadata/refresh/:seriesId | Stub "todo" |

---

## 11. CONFIGURACIÓN DE RED (Ejemplo)

### Servidor con IP 192.168.0.12
```
┌─────────────────────────────────────────────────┐
│                   Servidor                       │
│  ┌─────────┐    ┌─────────┐    ┌─────────────┐  │
│  │   nginx │────│backend  │────│  postgres   │  │
│  │  :80    │    │  :3000  │    │   :5432     │  │
│  └─────────┘    └─────────┘    └─────────────┘  │
│       │              │              │           │
│       │         ┌─────────┐        │           │
│       │         │  redis  │        │           │
│       │         │  :6379  │        │           │
│       │         └─────────┘        │           │
│       │              │              │           │
│       │         ┌─────────┐        │           │
│       └────────▶│frontend │        │           │
│                 │  :3000  │        │           │
│                 └─────────┘        │           │
└─────────────────────────────────────────────────┘
         ▲
         │ HTTP (puerto 80)
         │
┌─────────────────────────────────────────────────┐
│              Red 192.168.1.x                     │
│         Tu PC accede a: 192.168.0.12            │
└─────────────────────────────────────────────────┘
```

---

## 12. CHANGELOG RESUMEN

| Versión | Fecha | Cambios Principales |
|---------|-------|---------------------|
| v0.1.0 | 2026-05-21 | Estructura inicial, backend Go + Fiber, frontend Vue/Nuxt |
| v0.1.1 | 2026-05-21 | Autenticación JWT completa |
| v0.1.2 | 2026-05-21 | Scanner engine implementado |
| v0.2.0 | 2026-05-21 | GetPage endpoint, ComicReader |
| v0.2.10 | 2026-05-22 | Fix npm/pinia, frontend corregido |
| v0.2.13 | 2026-05-22 | Fix SSL, modo dev frontend |
| v0.2.14 | 2026-05-22 | Handlers stubs implementados, scan async con Redis |

Ver archivo completo: `CHANGELOG.md`

---

## 13. CONTACTO/SOPORTE

Proyecto creado para uso self-hosted.
Repositorio: (pendiente de configurar git remote)

---

## 14. CREDITOS

Inspirado en Kavita, Jellyfin/Plex para lectura digital.
Stack: Go + Fiber / Vue + Nuxt / PostgreSQL / Redis / Docker