#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER agb WITH ENCRYPTED PASSWORD 'agbpassword';
    CREATE DATABASE agb;
    GRANT ALL PRIVILEGES ON DATABASE agb TO agb;
EOSQL