
\c hmmm
SET client_min_messages TO WARNING;

-- Trigger update update_time column
CREATE OR REPLACE FUNCTION update_update_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = now();
    return NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

BEGIN;

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    "type" NUMERIC NOT NULL,
    "name" TEXT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT REFERENCES roles (id) NOT NULL,
    permission BIGINT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    dob date NOT NULL,
    gender VARCHAR DEFAULT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    reviewed BOOLEAN NOT NULL DEFAULT FALSE,
    email CITEXT NOT NULL UNIQUE,
    confirmed_email BOOLEAN DEFAULT FALSE CHECK  (confirmed_email OR NOT active),
    phone_number VARCHAR UNIQUE DEFAULT NULL,
    confirmed_phone BOOLEAN DEFAULT FALSE,
    passhash TEXT NOT NULL,
    failed_logins_count INT DEFAULT 0,
    door_code VARCHAR DEFAULT NULL,
    external_payment_customer_id INT DEFAULT NULL,
    review_time TIMESTAMPTZ DEFAULT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE accounts ADD CONSTRAINT external_payment_restrictions CHECK
(
    (confirmed_email AND reviewed)
    OR
    (NOT confirmed_email AND external_payment_customer_id IS NULL)
    OR
    (NOT reviewed AND external_payment_customer_id IS NULL)
);
ALTER TABLE accounts ADD CONSTRAINT active_restrictions CHECK
(
    (confirmed_email AND reviewed)
    OR
    (NOT confirmed_email AND NOT active)
    OR
    (NOT reviewed AND NOT active)
);

CREATE TABLE account_roles (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT REFERENCES roles (id) NOT NULL,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE confirmations (
    id BIGSERIAL PRIMARY KEY,
    "type" NUMERIC NOT NULL, -- email, phone number or password reset
    confirmation_target VARCHAR DEFAULT NULL,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    key VARCHAR NOT NULL UNIQUE,
    confirm_time TIMESTAMPTZ DEFAULT NULL,
    failed_confirmations_count INT DEFAULT 0,
    expire_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + interval '5 minutes',
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account_events (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    "type" NUMERIC NOT NULL,
    note VARCHAR,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) UNIQUE NOT NULL,
    country VARCHAR(20) NOT NULL,
    city VARCHAR(20) NOT NULL,
    state_code VARCHAR(2) NOT NULL,
    street VARCHAR(80) NOT NULL,
    zip_code VARCHAR NOT NULL,
    "type" INT NOT NULL, -- type address and type billing address
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    last_activity_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    session_id UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE equipment (
    id BIGSERIAL PRIMARY KEY,
    "type" INT NOT NULL,
    "description" text,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE authorizations (
    id BIGSERIAL PRIMARY KEY,
    equipment_id BIGINT REFERENCES equipment (id) NOT NULL,
    controller_id BIGINT REFERENCES equipment (id),
    "type" INT NOT NULL,
    "description" text,
    renewable BOOLEAN NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account_authorizations (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    authorization_id BIGINT REFERENCES authorizations (id) NOT NULL,
    active BOOLEAN NOT NULL,
    efective_time TIMESTAMPTZ,
    renewal_time TIMESTAMPTZ,
    expire_time TIMESTAMPTZ,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;