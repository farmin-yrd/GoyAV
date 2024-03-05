CREATE DATABASE docu_db;

\c docu_db;


CREATE SCHEMA IF NOT EXISTS docu_schema;

CREATE ROLE docu_user WITH LOGIN PASSWORD 'docu_password';

GRANT ALL PRIVILEGES ON SCHEMA docu_schema TO docu_user;

\q
