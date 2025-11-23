#!/bin/bash

PG_DATABASE_NAME=${PG_DATABASE_NAME:-avito_pr}
PG_USER=${PG_USER:-postgres}
PG_PASSWORD=${PG_PASSWORD:-postgres}
MIGRATION_DIR=${MIGRATION_DIR:-migrations}

export MIGRATION_DSN="host=pg port=5432 dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable"

echo "Waiting for PostgreSQL to be ready..."
sleep 5

echo "Running migrations from ${MIGRATION_DIR}..."
goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v