#!/bin/bash

# This script applies database migrations
# It assumes you have psql installed and the DATABASE_URL environment variable set

set -e

if [ -z "$DATABASE_URL" ]; then
    echo "DATABASE_URL environment variable is not set"
    exit 1
fi

echo "Applying database migrations..."

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_DIR="$SCRIPT_DIR/../internal/infra/postgres/migrations"

# Apply migrations in order
for migration in "$MIGRATIONS_DIR"/*.sql; do
    if [ -f "$migration" ]; then
        echo "Applying migration: $(basename "$migration")"
        psql "$DATABASE_URL" -f "$migration"
    fi
done

echo "All migrations applied successfully!"