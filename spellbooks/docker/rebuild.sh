#!/bin/sh
set -e

echo "Rebuilding with no cache..."
docker compose build --no-cache

echo "Starting services..."
docker compose up -d

echo "Done. Running containers:"
docker compose ps
