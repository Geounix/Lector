# Lector Comics - Documento Ingeniería de Implementación

Versión 0.3

------------------------------------------------------------------------

# 1. Docker Producción Completo

Objetivo:

Despliegue Linux en 1 comando.

Servicios:

nginx frontend backend worker postgres redis

Volúmenes:

postgres_data cache_data metadata_data logs_data backup_data
library_data

Healthchecks:

backend: GET /health

frontend: GET /

postgres: pg_isready

redis: redis-cli ping

Restart Policy:

always

Networking:

internal_network

external_network

TLS:

NGINX reverse proxy

HTTP2

gzip

brotli futuro

------------------------------------------------------------------------

# 2. Redis Diseño

Prefijos:

auth:session:{id}

reader:progress:{user}

metadata:{series}

dashboard:{user}

search:{query}

thumbnail:{chapter}

TTL:

metadata:

24 horas

dashboard:

5 minutos

thumbnail:

30 dias

Invalidación:

actualizacion metadata

escaneo nuevo

cambio biblioteca

Límite RAM:

512MB inicial

Eviction:

allkeys-lru

------------------------------------------------------------------------

# 3. PostgreSQL SQL Real

USERS

id BIGSERIAL

username VARCHAR(50)

email VARCHAR(255)

password_hash TEXT

avatar TEXT

theme VARCHAR(20)

language VARCHAR(10)

created_at TIMESTAMP

updated_at TIMESTAMP

ROLES

id BIGSERIAL

name VARCHAR(30)

permission JSONB

LIBRARY

id BIGSERIAL

path TEXT

watch_enabled BOOLEAN

scan_interval INTEGER

SERIES

id BIGSERIAL

title TEXT

sort_title TEXT

publisher TEXT

rating NUMERIC

language VARCHAR(10)

CHAPTER

id BIGSERIAL

file_path TEXT

size BIGINT

pages INTEGER

hash TEXT

READ_PROGRESS

page INTEGER

percentage NUMERIC

time_spent BIGINT

Índices:

GIN metadata

GIN search

BTREE email

BTREE hash

Particiones:

activity_logs

future metrics

------------------------------------------------------------------------

# 4. Swagger OpenAPI

Auth:

POST login

POST refresh

POST logout

Library:

GET library

POST library

PATCH library

DELETE library

Reader:

GET chapter

POST progress

Search:

GET search

Metadata:

POST refresh

GET metadata

Stats:

GET stats

Admin:

POST reindex

POST scan

Responses:

200

201

400

401

403

500

JWT:

Bearer obligatorio

Versionado:

v1

v2 futuro

------------------------------------------------------------------------

# 5. Scanner Engine Interno

Pipeline:

Filesystem Event

↓

Worker Queue

↓

SHA256

↓

Validar extensión

↓

Extraer nombre limpio

↓

Regex:

Volume

Chapter

Special

↓

Metadata local

↓

Metadata externa

↓

Generar thumbnail

↓

OCR opcional

↓

Guardar BD

↓

Actualizar Redis

↓

Actualizar búsqueda

↓

Notificar frontend

Workers:

ScannerWorker

MetadataWorker

ThumbnailWorker

CleanupWorker

IndexWorker

SchedulerWorker

Throttle:

CPU 2:

4 workers

CPU 4:

8 workers

CPU 8:

16 workers

CPU 16:

32 workers

Regex ejemplo:

One Piece Vol 10 Ch 100.cbz

Resultado:

Serie:

One Piece

Vol:

10

Cap:

100

------------------------------------------------------------------------

# 6. Thumbnail Engine

Formatos:

CBZ

PDF

EPUB

Proceso:

Extraer primera página

↓

Resize

↓

WEBP

↓

Guardar SSD

Resoluciones:

256

512

1024

Cache:

Redis

Filesystem

------------------------------------------------------------------------

# 7. Metadata Providers

Sistema Provider:

Provider Interface

Refresh()

Search()

Normalize()

Fuentes:

AniList

ComicVine

OpenLibrary

Google Books

Prioridad:

Metadata local

↓

AniList

↓

ComicVine

↓

OpenLibrary

Fallback:

nombre archivo

Cache:

Redis

TTL:

24h

Rate Limit externo:

provider queue

retry automático

------------------------------------------------------------------------

# 8. Sync Dispositivos

Objetivo:

PC

Tablet

Celular

Proceso:

Progreso local

↓

Redis

↓

PostgreSQL

↓

WebSocket

↓

Dispositivo conectado

Conflictos:

última lectura gana

Sync:

favoritos

progreso

colecciones

config usuario

------------------------------------------------------------------------

# 9. Frontend Vue Arquitectura

Views:

Home

Library

Reader

Settings

Admin

Components:

Sidebar

ReaderCanvas

BookCard

MetadataPanel

SearchBar

ProgressBar

ThemeSelector

Stores:

AuthStore

LibraryStore

ReaderStore

MetadataStore

StatsStore

Middleware:

JWT

Roles

Theme

Offline Cache

PWA:

Instalable

Android

IOS

Desktop

------------------------------------------------------------------------

# 10. Actualizaciones Automáticas

Objetivo:

docker pull

↓

descargar

↓

validar

↓

backup

↓

deploy

↓

healthcheck

Rollback:

automático

Versionado:

SemVer

1.0.0

1.0.1

1.1.0

2.0.0

------------------------------------------------------------------------

# 11. HA Futuro

Load Balancer

↓

Backend 1

Backend 2

↓

Redis

↓

PostgreSQL Replica

Cluster:

futuro

Storage:

NFS

S3 futuro

Replicación:

BD

metadata

config

------------------------------------------------------------------------

# 12. Disaster Recovery

RPO:

1 hora

RTO:

30 minutos

Escenarios:

corrupción BD

disco falla

metadata dañada

rollback actualización

Backups:

diario

semanal

mensual

Verificación:

checksum

restore prueba mensual

------------------------------------------------------------------------

# 13. Seguridad Avanzada

Headers:

CSP

XSS Protection

HSTS

Rate Limit:

100 minuto

IP Ban:

bruteforce

Logs:

JSON

SIEM futuro

2FA:

TOTP

JWT:

rotación refresh

Argon2:

password hashing

Secrets:

Docker Secrets

------------------------------------------------------------------------

# 14. Benchmarks

Objetivos:

100000 comics

Escaneo:

15 minutos

API:

150ms

Reader:

1 segundo

Redis:

90% cache hit

RAM:

2GB

CPU:

2 cores

------------------------------------------------------------------------

# 15. Milestone Desarrollo

Sprint 1:

infraestructura

Sprint 2:

scanner

Sprint 3:

reader

Sprint 4:

metadata

Sprint 5:

dashboard

Sprint 6:

sync

Sprint 7:

api publica

Sprint 8:

release

------------------------------------------------------------------------

# Estado Proyecto

Arquitectura lista.

Especificación lista.

Listo para iniciar implementación.
