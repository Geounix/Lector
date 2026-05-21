# Lector Comics - Especificación Técnica Empresarial

Versión 0.2

------------------------------------------------------------------------

# 1. Arquitectura Empresarial Completa

Objetivo:

Crear una plataforma Linux First self-hosted escalable para:

-   Manga
-   Comics
-   Webtoon
-   EPUB
-   PDF
-   Novelas ligeras

Diseño:

Frontend ↓ NGINX ↓ API Gateway ↓ Backend GO (Fiber)

Servicios internos:

Auth Service Scanner Service Metadata Service Reader Service Stats
Service Search Service

Persistencia:

PostgreSQL

Cache:

Redis

Archivos:

Filesystem Linux

Procesamiento:

Workers Paralelos

------------------------------------------------------------------------

# 2. Estructura Monorepo

lector-comics/

backend/

cmd/

api/

internal/

models/

repository/

middleware/

services/

scanner/

reader/

metadata/

security/

workers/

cache/

frontend/

components/

pages/

stores/

layouts/

plugins/

infra/

docker/

nginx/

postgres/

redis/

scripts/

tests/

integration/

benchmark/

docs/

api/

erd/

architecture/

backups/

metadata/

cache/

logs/

library/

------------------------------------------------------------------------

# 3. ERD (Diseño Base Datos)

USERS

id username email password_hash avatar theme language 2fa_enabled
created_at updated_at

ROLES

id name permissions

USER_ROLE

user_id role_id

LIBRARY

id name path type watch_enabled scan_interval

SERIES

id library_id title sort_title description publisher year language cover
rating

VOLUME

id series_id number

CHAPTER

id volume_id title file_path hash size pages created_at

READING_PROGRESS

user_id chapter_id page percentage time_spent last_read

COLLECTION

id name owner

COLLECTION_ITEMS

collection_id chapter_id

METADATA

id series_id writer artist genre tag release_date

REVIEW

id user_id series_id score

ACTIVITY_LOG

id user_id action ip timestamp

DEVICE

id user_id token platform

ACHIEVEMENT

id name condition

USER_ACHIEVEMENT

user_id achievement_id

BACKUP_JOB

id type status

SYSTEM_CONFIG

key value

------------------------------------------------------------------------

# 4. Flujo Scanner Interno

Nuevo archivo

↓

Watch Folder

↓

SHA256

↓

Detectar formato

↓

CBZ

CBR

EPUB

PDF

↓

Metadata local

↓

Metadata externa

↓

Normalizar nombres

↓

Miniatura

↓

OCR opcional

↓

Guardar PostgreSQL

↓

Actualizar Redis

↓

Actualizar búsqueda

↓

Notificar frontend websocket

Workers:

scanner-worker

thumbnail-worker

ocr-worker

cleanup-worker

metadata-worker

index-worker

scheduler-worker

Escaneo paralelo:

CPU 2 núcleos

4 workers

CPU 4 núcleos

8 workers

CPU 8 núcleos

16 workers

------------------------------------------------------------------------

# 5. Diseño API

AUTH

POST /api/auth/login

POST /api/auth/logout

POST /api/auth/register

POST /api/auth/refresh

LIBRARY

GET /api/library

POST /api/library

DELETE /api/library

PATCH /api/library

READER

GET /api/reader/chapter

POST /api/reader/progress

SEARCH

GET /api/search

METADATA

GET /api/metadata

POST /api/metadata/refresh

USER

GET /api/users

PATCH /api/users

DELETE /api/users

STATS

GET /api/stats

ADMIN

GET /api/admin/jobs

POST /api/admin/reindex

WEBSOCKET

progress_update

scan_finished

metadata_updated

dashboard_refresh

OpenAPI obligatorio

Swagger obligatorio

Rate limit:

100 requests minuto

JWT obligatorio

------------------------------------------------------------------------

# 6. RBAC

ADMIN

Todo acceso

MODERADOR

Biblioteca

Metadata

Usuarios limitados

USUARIO

Lectura

Favoritos

Progreso

NIÑO

Restricción contenido

INVITADO

Solo lectura pública

Sistema permisos:

library.create

library.delete

user.delete

reader.use

admin.jobs

metadata.refresh

------------------------------------------------------------------------

# 7. Reader Engine

Comic Engine:

Precarga 5 páginas

Cache RAM

Cache Disco

Modo:

RTL

LTR

Webtoon

Doble página

Fullscreen

Gestos:

Swipe

Zoom

Keyboard:

A

D

F

B

Cache:

Redis

Filesystem

EPUB:

Motor dedicado

Bookmarks

Fuentes custom

Notas

PDF:

Render progresivo

OCR futuro

------------------------------------------------------------------------

# 8. Docker Producción

Servicios:

nginx

frontend

backend

worker

postgres

redis

healthcheck

network interna

volúmenes:

library

cache

postgres_data

metadata

backups

logs

Restart:

always

Healthcheck:

backend

frontend

postgres

redis

------------------------------------------------------------------------

# 9. Estrategia PostgreSQL

Migraciones:

golang-migrate

versionadas

001_users

002_library

003_series

004_metadata

005_progress

Índices:

series_title

chapter_hash

user_email

reading_progress

Particionado futuro:

activity_log

benchmark

logs

------------------------------------------------------------------------

# 10. Sistema Búsqueda

Inicial:

PostgreSQL Full Text

Futuro:

OpenSearch

Campos indexados:

title

genre

writer

artist

tags

publisher

Sinopsis

Búsqueda:

aproximada

exacta

filtros

autocomplete

------------------------------------------------------------------------

# 11. Backups y DRP

Manual

Automático

Incremental

Full semanal

Retención:

7 diarios

4 semanales

3 mensuales

Restauración:

BD

metadata

cache

library

Objetivo recuperación:

RPO 1 hora

RTO 30 minutos

------------------------------------------------------------------------

# 12. CI/CD

GitHub

↓

Pull Request

↓

Tests

↓

Lint

↓

Build Docker

↓

Security Scan

↓

Deploy

↓

Smoke Test

Herramientas:

Github Actions

Trivy

GolangCI

------------------------------------------------------------------------

# 13. Plugins

Arquitectura:

plugin-reader

plugin-metadata

plugin-ai

plugin-export

plugin-backup

Interfaz:

Enable

Disable

Versionado

Sandbox futuro

------------------------------------------------------------------------

# 14. Monitoreo

Métricas:

CPU

RAM

Usuarios

Tiempo API

Redis hit ratio

Errores

Workers activos

Integraciones:

Grafana

Prometheus

Zabbix

Logs:

JSON estructurado

Rotación

Compresión

------------------------------------------------------------------------

# 15. Benchmarks Objetivo

Escaneo:

50k archivos

menos 15 minutos

API:

menos 150ms

Reader:

menos 1 segundo apertura

Búsqueda:

menos 300ms

RAM:

2GB mínimo

------------------------------------------------------------------------

# 16. Roadmap Técnico

v0.1

Core sistema

v0.3

Metadata

v0.5

Dashboard

v0.7

OCR

v1.0

API estable

v2.0

Cluster

Apps móviles

AI metadata

Federación servidores

------------------------------------------------------------------------

# Objetivo Final

Construir una plataforma Linux First profesional preparada para crecer
desde servidor personal hasta despliegues empresariales.
