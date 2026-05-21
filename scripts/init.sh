#!/bin/bash

set -e

echo "==================================="
echo "Lector Comics - Setup Script"
echo "==================================="

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

echo ""
echo "[1/6] Checking prerequisites..."

if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "Error: Docker Compose is not installed"
    exit 1
fi

echo "  ✓ Docker found"
echo "  ✓ Docker Compose found"

echo ""
echo "[2/6] Creating necessary directories..."

mkdir -p library metadata cache backups logs

echo "  ✓ Directories created"

echo ""
echo "[3/6] Building Docker images..."

docker compose build

echo "  ✓ Images built"

echo ""
echo "[4/6] Starting services..."

docker compose up -d postgres redis

echo "  Waiting for database to be ready..."
sleep 5

for i in {1..30}; do
    if docker compose exec -T postgres pg_isready -U lector -d lector &> /dev/null; then
        echo "  ✓ Database is ready"
        break
    fi
    echo "  Waiting for database... ($i/30)"
    sleep 2
done

echo ""
echo "[5/6] Starting all services..."

docker compose up -d

echo "  ✓ All services started"

echo ""
echo "[6/6] Checking health..."

sleep 5

if curl -sf http://localhost:3000/health &> /dev/null; then
    echo "  ✓ Backend is healthy"
else
    echo "  ⚠ Backend may still be starting..."
fi

echo ""
echo "==================================="
echo "Setup Complete!"
echo "==================================="
echo ""
echo "Services:"
echo "  - Frontend: http://localhost:8080"
echo "  - API:      http://localhost:3000"
echo "  - API Docs: http://localhost:3000/api/v1"
echo ""
echo "To stop: docker compose down"
echo "To view logs: docker compose logs -f"
echo "==================================="