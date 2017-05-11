#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname agb <<-EOSQL
    CREATE TYPE img_format AS ENUM ('jpeg', 'png', 'gif');
    CREATE TABLE gif
        (id bytea PRIMARY KEY, 
        url varchar[32], 
        characteristics_init integer[], 
        characteristics integer[], 
        match bytea[], 
        img_size integer, 
        img_type img_format, 
        source varchar[32], 
        create_at timestamp, 
        last_visit timestamp);
    CREATE TABLE submit_match
        (id uuid PRIMARY KEY,
        home bytea,
        away bytea,
        submitter varchar[32],
        submit_at TIMESTAMP);
EOSQL