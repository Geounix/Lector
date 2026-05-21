# Lector Comics - Documento Maestro del Proyecto

Version: 0.1 Arquitectura Linux First

------------------------------------------------------------------------

# Vision del Proyecto

Crear una plataforma self-hosted moderna inspirada en Kavita pero
diseñada desde cero para:

-   Manga
-   Comics occidentales
-   Webtoon
-   EPUB
-   PDF
-   Novelas ligeras
-   Archivos de imagen

Objetivo principal:

"Crear el Jellyfin/Plex definitivo para lectura digital"

Principios:

-   Linux First
-   Docker First
-   Open Source
-   Bajo consumo recursos
-   Escalable
-   Modular
-   API First
-   Seguridad por diseño
-   Instalacion simple

------------------------------------------------------------------------

# FASE 0 - Definicion del Producto

## Objetivos funcionales

El sistema debe permitir:

### Biblioteca

-   Multiples bibliotecas
-   Escaneo automatico
-   Watch folders
-   Hash archivos
-   Metadata automatica
-   Agrupacion inteligente
-   Reindexado

### Lectura

Comic:

-   Pagina simple
-   Doble pagina
-   Zoom
-   Scroll continuo
-   Lectura RTL manga
-   Lectura LTR comic
-   Webtoon vertical
-   Marcadores
-   Continuar lectura

EPUB:

-   Fuentes custom
-   Tema oscuro
-   Notas
-   Highlight

PDF:

-   Render progresivo
-   Zoom
-   OCR futuro

### Usuario

-   Multiusuario
-   Restriccion edad
-   Roles
-   Favoritos
-   Want To Read
-   Historial

### Estadisticas

-   Paginas leidas
-   Horas lectura
-   Series terminadas
-   Dias consecutivos

------------------------------------------------------------------------

# FASE 1 - Arquitectura Linux First

## Backend

Go (Golang)

Motivos:

-   Binario unico
-   Excelente Linux
-   Bajo RAM
-   Muy rapido
-   Facil Docker
-   Excelente concurrencia

Framework:

Fiber

Estructura:

backend/

cmd/

internal/

api/

services/

workers/

models/

repository/

middleware/

scanner/

metadata/

reader/

cache/

security/

config/

pkg/

tests/

## Frontend

Vue + Nuxt

Motivos:

-   Ligero
-   Facil mantenimiento
-   Responsive
-   PWA

Estructura:

frontend/

components/

pages/

layouts/

stores/

services/

plugins/

middleware/

assets/

reader/

dashboard/

admin/

## Base datos

PostgreSQL

Motivos:

-   JSON nativo
-   Full Text Search
-   Escalable
-   Excelente Linux

## Cache

Redis

Uso:

-   Cache metadata
-   Sesiones
-   Queue
-   Dashboard cache

## Workers

Go Workers internos

Responsabilidades:

-   Escaneo
-   OCR
-   Metadata
-   Miniaturas
-   Limpieza cache
-   Reindexado

## Docker

Servicios:

frontend

backend

postgres

redis

worker

nginx

Instalacion:

git clone proyecto

docker compose up -d

Objetivo recursos:

2 CPU

2GB RAM

50k comics

------------------------------------------------------------------------

# FASE 2 - Modelo Datos

## Tabla Users

id

username

email

password_hash

avatar

role

created_at

last_login

2fa_enabled

language

theme

## Roles

Admin

Moderador

Adulto

Niño

Invitado

## Library

id

name

type

path

watch_enabled

scan_interval

created_at

## Series

id

title

sort_title

description

cover

year

status

publisher

language

metadata_source

age_rating

## Volume

id

series_id

number

title

## Chapter

id

volume_id

file_path

hash

page_count

size

created_at

## ReadingProgress

user_id

chapter_id

page

percentage

last_read

time_spent

## Metadata

genres

tags

artist

writer

isbn

rating

release_date

## Collections

Favoritos

Pendientes

Leyendo

Completados

Personalizadas

------------------------------------------------------------------------

# FASE 3 - Escaner

Pipeline:

Archivo

↓

SHA256

↓

Analisis extension

↓

Metadata local

↓

Miniatura

↓

OCR opcional

↓

Guardar PostgreSQL

↓

Actualizar cache

Workers:

scanner-worker

metadata-worker

thumbnail-worker

cleanup-worker

ocr-worker

Formatos:

CBZ

CBR

ZIP

RAR

7Z

PDF

EPUB

PNG

WEBP

AVIF

JPG

JPEG

Detecciones:

Series

Volumen

Capitulo

Idioma

Autor

Nombre limpio

------------------------------------------------------------------------

# FASE 4 - Reader Engine

Objetivo:

Crear mejor experiencia posible.

## Comic Reader

Caracteristicas:

-   Zoom inteligente
-   Precarga paginas
-   Render progresivo
-   Cache RAM
-   Cache disco
-   Navegacion teclado
-   Gestos movil
-   Modo pantalla completa

Atajos:

A pagina anterior

D pagina siguiente

F fullscreen

B bookmark

## Webtoon

Scroll vertical

Precarga imagenes

Carga dinamica

## EPUB

Motor EPUB dedicado

Bookmarks

Fuentes custom

Tema oscuro

Tema AMOLED

## PDF

PDF Renderer

Zoom

Busqueda texto

OCR futuro

------------------------------------------------------------------------

# FASE 5 - Metadata

Fuentes futuras:

ComicVine

AniList

MyAnimeList

OpenLibrary

Google Books

Campos:

Titulo

Autor

Genero

Portada

Sinopsis

Editorial

Año

Estado

Cache local:

TTL configurable

Fallback metadata local

------------------------------------------------------------------------

# FASE 6 - Seguridad

JWT

Argon2

HTTPS

Rate Limit

Logs seguridad

2FA

Captcha login futuro

Permisos:

Admin:

todo

Usuario:

lectura

favoritos

Moderador:

biblioteca

Restriccion edad:

13+

16+

18+

Logs:

login

logout

errores

api

acciones admin

------------------------------------------------------------------------

# FASE 7 - Dashboard

Widgets:

Continuar leyendo

Favoritos

Recientes

Pendientes

Mas leidos

Tiempo lectura

Estadisticas

Personalizacion:

drag drop

ocultar widgets

orden custom

filtros inteligentes

------------------------------------------------------------------------

# FASE 8 - API

/api/library

/api/users

/api/search

/api/reader

/api/metadata

/api/stats

/api/auth

Documentacion:

OpenAPI

Swagger

API Keys futuras

Websocket:

progreso lectura

dashboard tiempo real

------------------------------------------------------------------------

# FASE 9 - Funciones Avanzadas

OCR

Offline cache

Smart Filters

Sync dispositivos

Descarga offline

Export progreso

Import progreso

Reading goals

Achievements

Colecciones inteligentes

Sistema recomendaciones futuro

IA metadata futuro

------------------------------------------------------------------------

# FASE 10 - Infraestructura

NGINX

Redis

PostgreSQL

Go Backend

Workers

Docker

Logs:

rotacion

retencion

compresion

Backup:

automatico

manual

incremental

Monitoreo:

CPU

RAM

Usuarios activos

Errores API

Tiempo respuesta

Disco

Cache hit ratio

Integraciones futuras:

Grafana

Prometheus

Zabbix

------------------------------------------------------------------------

# FASE 11 - QA

Unitarias

Integracion

Stress Test

Carga

UX

Seguridad

Objetivo:

10000+

50000+

100000+

archivos

Benchmark:

escaneo

lectura

render

API

------------------------------------------------------------------------

# FASE 12 - Release

Alpha

Beta

RC

Stable

Pipeline:

Git

↓

Build

↓

Tests

↓

Docker

↓

Deploy

CI/CD futuro

Github Actions

------------------------------------------------------------------------

# Roadmap

v0.1

Login

Biblioteca

CBZ

CBR

Docker

Linux

v0.5

PDF

EPUB

Metadata

Dashboard

v1.0

OCR

API

Sync

Estadisticas

v2.0

Apps Android

Apps IOS

Cluster soporte

IA metadata

------------------------------------------------------------------------

# Monorepo Final

lector-comics/

backend/

frontend/

infra/

docker/

docs/

scripts/

config/

tests/

library/

metadata/

cache/

backups/

logs/

------------------------------------------------------------------------

# Objetivo Final

Construir la mejor plataforma self-hosted Linux para lectura digital
inspirada en Kavita pero optimizada para rendimiento, simplicidad y
escalabilidad profesional.
