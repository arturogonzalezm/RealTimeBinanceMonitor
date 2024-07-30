-- database/init.sql

-- Ensure the 'postgres' role exists
DO
$do$
    BEGIN
        IF NOT EXISTS (SELECT
                       FROM pg_catalog.pg_roles
                       WHERE rolname = 'postgres') THEN
            CREATE ROLE postgres;
        END IF;
    END
$do$;

-- Create a new database if it does not exist
DO
$do$
    BEGIN
        IF NOT EXISTS (SELECT
                       FROM pg_database
                       WHERE datname = 'postgres') THEN
            CREATE DATABASE postgres;
        END IF;
    END
$do$;

-- Create a new user with a password if it does not exist
DO
$do$
    BEGIN
        IF NOT EXISTS (SELECT
                       FROM pg_catalog.pg_roles
                       WHERE rolname = 'postgres') THEN
            CREATE USER postgres WITH ENCRYPTED PASSWORD 'postgres';
        END IF;
    END
$do$;

-- Grant all privileges on the new database to the new user
GRANT ALL PRIVILEGES ON DATABASE postgres TO postgres;

-- Connect to the created database to run further SQL commands
\c postgres

-- Create the exchange_info table
CREATE TABLE IF NOT EXISTS exchange_info
(
    symbol       VARCHAR(50),
    status       VARCHAR(50),
    base_asset   VARCHAR(50),
    quote_asset  VARCHAR(50),
    filter_type  VARCHAR(50),
    filter_key   VARCHAR(50),
    filter_value VARCHAR(50),
    CONSTRAINT exchange_info_unique UNIQUE (symbol, filter_type, filter_key)
);

CREATE TABLE IF NOT EXISTS ticker_data
(
    id           SERIAL PRIMARY KEY,
    event_time   BIGINT,
    symbol       TEXT,
    last_price   DOUBLE PRECISION,
    price_change DOUBLE PRECISION,
    high_price   DOUBLE PRECISION,
    low_price    DOUBLE PRECISION,
    volume       DOUBLE PRECISION,
    quote_volume DOUBLE PRECISION,
    open_time    BIGINT,
    close_time   BIGINT,
    trade_count  INT,
    latency      BIGINT
);

