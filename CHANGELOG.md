# Changelog - Lector Comics

## v0.2.15 - 2026-05-22

### Added
- **scanner/utils.go**: Nuevo archivo con funciones utilitarias compartidas
  - `CalculateHash` y `CalculateHashPartial` (hash parcial de 1024 bytes para archivos grandes)
  - `IsImageFile`, `GetFileSize`, `GenerateSortTitle`
  - `ExtractVolumeNumber`, `ExtractChapterNumber`
  - Elimina código duplicado entre scanner.go y services/scanner_service.go

- **library.vue**: Corregido `openSeries()` para cargar datos completos de serie
  - Ahora hace `GET /api/v1/series/:id` antes de navegar al reader
  - Obtiene correctamente el primer capítulo disponible de la serie

- **EditLibraryModal.vue**: Nuevo campo `scan_interval` editable
  - Input numérico para configurar intervalo de escaneo
  - Valor cargado desde la biblioteca y enviado correctamente al API

- **reader.vue**: Múltiples mejoras
  - `saveProgress()` ahora es `async` con `await` en el fetch
  - Separado `loading` (capítulo) de `pageLoading` (páginas)
  - Mejor UX durante carga de páginas

### Fixed
- **handlers.go:GetSeriesByID**: Corregido para incluir volúmenes sin capítulos
  - Ahora todos los volúmenes aparecen en la respuesta, no solo los que tienen capítulos asignados

- **scanner.go**: Removidas funciones duplicadas, ahora usa utils.go
- **services/scanner_service.go**: Refactorizado para usar funciones de `scanner` package
- **workers/scanner.go**: Refactorizado para usar funciones de `scanner` package
  - Usa `scanner.CalculateHashPartial` en lugar de `scanner.CalculateHash` para optimizar memoria
  - Usa `scanner.GenerateSortTitle`, `scanner.GetFileSize`

## v0.2.14 - 2026-05-22

### Added
- Tabla `scan_job` en PostgreSQL para tracking de scans asíncronos
  - Campos: id (UUID), library_id, user_id, status, chapters_found, error_message, started_at, completed_at
  - Índices para library_id, user_id, status

- Endpoints API nuevos:
  - `PATCH /api/v1/library/:id` - Actualizar biblioteca
  - `DELETE /api/v1/library/:id` - Eliminar biblioteca
  - `POST /api/v1/library/:id/scan` - Iniciar scan asíncrono (retorna scan_id)
  - `GET /api/v1/scans/:scanId` - Consultar estado de scan

- Sistema de scan asíncrono con Redis pub/sub:
  - Worker se suscribe al canal `lector:scan_library`
  - API publica mensaje con library_id y scan_id
  - Worker recibe, procesa y actualiza estado en BD

- Frontend: Modales para Editar/Eliminar biblioteca

### Implemented Services
- **UpdateLibrary**: Actualiza name, type, path, watch_enabled, scan_interval
- **DeleteLibrary**: Elimina biblioteca (CASCADE elimina series/chapters)
- **GetSeriesByID**: Retorna serie + volumes + chapters
- **TriggerScan**: Crea scan_job y publica en Redis
- **GetScanStatus**: Consulta estado de scan por UUID

### Changed Files
- **infra/postgres/init.sql**: Nueva tabla `scan_job` con índices
- **backend/internal/repository/repository.go**:
  - +CreateScanJob, GetScanJob, UpdateScanJob, UpdateScanJobStarted, PublishScanCommand
- **backend/internal/services/services.go**:
  - +UpdateLibrary, DeleteLibrary, GetSeriesByID, GetVolumesBySeries, TriggerScan, GetScanStatus
  - +Volume, ScanStatus types
- **backend/internal/handlers/handlers.go**: Implementados UpdateLibrary, DeleteLibrary, GetSeriesByID, ScanLibrary, GetScanStatus
- **backend/cmd/api/main.go**: Rutas actualizadas para PATCH/DELETE library/:id, POST library/:id/scan, GET scans/:scanId
- **backend/internal/workers/scanner.go**:
  - +SubscribeToScans, handleScanCommand, scanLibraryWithCount, scanSeriesDirectoryWithCount, processComicFileWithCount
  - +ScanCommand type
- **backend/cmd/worker/main.go**: Goroutine para SubscribeToScans
- **frontend/pages/library.vue**: Botones Scan/Edit/Delete con polling
- **frontend/components/EditLibraryModal.vue**: Modal para editar biblioteca (nuevo)
- **frontend/components/DeleteLibraryModal.vue**: Modal para eliminar biblioteca (nuevo)

### API Endpoints States
| Método | Ruta | Handler | Estado |
|--------|------|---------|--------|
| PATCH | /api/v1/library/:id | UpdateLibrary | ✅ Implementado |
| DELETE | /api/v1/library/:id | DeleteLibrary | ✅ Implementado |
| GET | /api/v1/series/:id | GetSeriesByID | ✅ Implementado |
| POST | /api/v1/library/:id/scan | ScanLibrary | ✅ Implementado |
| GET | /api/v1/scans/:scanId | GetScanStatus | ✅ Implementado |

### Notas Técnicas
- Scan polling: Frontend hace polling cada 2s, máximo 30 intentos (60s timeout)
- Worker SubscribeToScans usa Redis pub/sub channel
- ON DELETE CASCADE en scan_job.library_id elimina jobs al borrar library

### Fixed
- Error compilación Go: "errMsg declared and not used"
  - scanner.go:121 - Cambiado `&errMsgStr` por `errMsg` en UpdateScanJob call
  - El puntero a variable local estaba siendo usado directamente en lugar de la variable declarada

---

## DOCUMENTACIÓN COMPLETA

Ver archivo **`PROJECT_INFO.md`** para:
- Arquitectura completa con diagrama de servicios
- Estructura del proyecto
- Lista detallada de todos los API endpoints
- Variables de entorno
- Schema de base de datos
- Estado del proyecto (implementado vs pendiente)
- Guía de deployment
- Limitaciones conocidas
- Flujo de scan asíncrono
- Configuración de red de ejemplo
- Resumen del changelog