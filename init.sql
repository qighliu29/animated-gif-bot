CREATE USER agb;
CREATE DATABASE agb;
GRANT ALL PRIVILEGES ON DATABASE agb TO agb;

CREATE TYPE img_format AS ENUM ('jpeg', 'png', 'gif');
CREATE TYPE img_hash AS bytea;
CREATE TABLE gif
    (id img_hash PRIMARY KEY, 
    url varchar[32], 
    characteristics_init integer[], 
    characteristics integer[], 
    match img_hash[], 
    img_size integer, 
    img_type img_format, 
    source varchar[32], 
    create_at timestamp, 
    last_visit timestamp);
CREATE TABLE submit_match
    (id uuid PRIMARY KEY,
    home img_hash,
    away img_hash,
    submitter varchar[32],
    submit_at TIMESTAMP);