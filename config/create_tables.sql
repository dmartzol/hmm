
\c hackerspace
SET client_min_messages TO WARNING;

-- Trigger update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    return NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

BEGIN;

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    "type" NUMERIC NOT NULL,
    "name" TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    dob date NOT NULL,
    gender VARCHAR(1) DEFAULT NULL,
    role_id BIGINT REFERENCES roles (id),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    email CITEXT NOT NULL UNIQUE,
    phone_number VARCHAR(20) UNIQUE DEFAULT NULL,
    passhash TEXT NOT NULL,
    failed_login_attempts INT DEFAULT 0,
    door_code VARCHAR DEFAULT NULL,
    external_payment_customer_id INT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE confirmation_codes (
    id BIGSERIAL PRIMARY KEY,
    "type" NUMERIC NOT NULL, -- email or phone number
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    confirmed_at TIMESTAMP DEFAULT NULL,
    confirmation_code VARCHAR NOT NULL,
    code_expiration_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account_events (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    "type" NUMERIC NOT NULL,
    notes VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    country VARCHAR(20) NOT NULL,
    city VARCHAR(20) NOT NULL,
    state_code VARCHAR(2) NOT NULL,
    street VARCHAR(80) NOT NULL,
    -- zip
    "type" INT NOT NULL, -- type address and type billing address
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    token VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE equipment (
    id BIGSERIAL PRIMARY KEY,
    "type" INT NOT NULL,
    "description" text,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE authorizations (
    id BIGSERIAL PRIMARY KEY,
    equipment_id BIGINT REFERENCES equipment (id) NOT NULL,
    controller_id BIGINT REFERENCES equipment (id),
    "type" INT NOT NULL,
    "description" text,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account_authorizations (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    authorization_id BIGINT REFERENCES authorizations (id) NOT NULL,
    efective TIMESTAMP,
    expires TIMESTAMP,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;