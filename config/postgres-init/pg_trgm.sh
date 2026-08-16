#!/bin/bash
set -e

# Ensure pg_trgm is created in the configured database during init
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
CREATE EXTENSION IF NOT EXISTS pg_trgm;
EOSQL
