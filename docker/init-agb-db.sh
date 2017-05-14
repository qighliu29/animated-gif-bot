#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname agb <<-EOSQL
    CREATE TYPE img_format AS ENUM ('jpeg', 'png', 'gif');
    CREATE TYPE hot_reason AS ENUM ('query', 'match')
    CREATE TABLE gif
        (id uuid PRIMARY KEY, 
        url varchar[32], 
        characteristics_init nvarchar(32)[], 
        characteristics integer[], 
        match uuid[], 
        img_size integer, 
        img_type img_format,
        img_hash bytea UNIQUE,
        source varchar[32], 
        create_at timestamp, 
        last_update timestamp);
    CREATE TABLE submit_match
        (id uuid PRIMARY KEY,
        home uuid,
        away uuid,
        submitter varchar[32],
        submit_at timestamp);
    CREATE TABLE gif_hot
        (id uuid PRIMARY KEY,
        create_at timestamp,
        reason hot_reason,
        query_hit int,
        match_hit int);
    CREATE TABLE tag
        (id serial PRIMARY KEY,
        content nvarchar(32)
        exp nvarchar(8));
EOSQL