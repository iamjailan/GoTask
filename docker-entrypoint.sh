#!/bin/sh
set -eu

: "${DATABASE_URL:?DATABASE_URL must be set}"

echo "Applying database migrations..."
./migrate

echo "Starting GoTask API..."
exec ./gotask
