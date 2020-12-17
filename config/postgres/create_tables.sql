
\c hmm
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

CREATE TABLE emails (
    id BIGSERIAL PRIMARY KEY,
    email CITEXT NOT NULL UNIQUE,
    confirmed BOOLEAN DEFAULT FALSE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE phone_numbers (
    id BIGSERIAL PRIMARY KEY,
    "number" CITEXT NOT NULL UNIQUE,
    confirmed BOOLEAN DEFAULT FALSE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    country VARCHAR(20) NOT NULL,
    city VARCHAR(20) NOT NULL,
    state_code VARCHAR(2) NOT NULL,
    street VARCHAR(80) NOT NULL,
    zip_code VARCHAR NOT NULL,
    "type" INT NOT NULL, -- type address and type billing address
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_addresses_update_time BEFORE UPDATE ON addresses FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    email_id BIGINT REFERENCES emails (id) NOT NULL,
    address_id BIGINT REFERENCES addresses (id) NOT NULL,
    phone_number_id BIGINT REFERENCES phone_numbers (id) DEFAULT NULL,
    first_name CITEXT NOT NULL,
    last_name CITEXT NOT NULL,
    dob date NOT NULL,
    gender CITEXT DEFAULT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE CHECK (review_time IS NOT NULL OR NOT active),
    passhash TEXT NOT NULL,
    failed_logins_count INT DEFAULT 0,
    review_time TIMESTAMPTZ DEFAULT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_accounts_update_time BEFORE UPDATE ON accounts FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

CREATE TABLE door_codes (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    door_code TEXT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX door_codes_account_idx ON door_codes (account_id);

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    "name" TEXT UNIQUE,
    permissions JSON DEFAULT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_roles_update_time BEFORE UPDATE ON roles FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

CREATE TABLE account_roles (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    role_id BIGINT REFERENCES roles (id) NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (role_id, account_id)
);
CREATE INDEX account_roles_account_idx ON account_roles (account_id);
CREATE INDEX account_roles_role_idx ON account_roles (role_id);

CREATE TABLE confirmations (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    "type" NUMERIC NOT NULL, -- email, phone number or password reset
    confirmation_target VARCHAR DEFAULT NULL,
    key VARCHAR NOT NULL UNIQUE,
    confirmation_time TIMESTAMPTZ DEFAULT NULL,
    failed_confirmations_count INT DEFAULT 0,
    expire_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + interval '5 hours',
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX confirmations_account_idx ON confirmations (account_id);

CREATE TABLE account_events (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    "type" NUMERIC NOT NULL,
    note VARCHAR,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX account_events_account_idx ON account_events (account_id);
CREATE TRIGGER update_account_events_update_time BEFORE UPDATE ON account_events FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    last_activity_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    token UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    expiration_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 year',
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_sessions_update_time BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

CREATE TABLE equipment (
    id BIGSERIAL PRIMARY KEY,
    "type" INT NOT NULL,
    "description" text,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_equipment_update_time BEFORE UPDATE ON equipment FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

CREATE TABLE authorizations (
    id BIGSERIAL PRIMARY KEY,
    equipment_id BIGINT REFERENCES equipment (id) NOT NULL,
    "type" INT NOT NULL,
    "description" text,
    renewable BOOLEAN NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX authorizations_equipment_idx ON authorizations (equipment_id);
CREATE TRIGGER update_authorizations_update_time BEFORE UPDATE ON authorizations FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

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
CREATE INDEX account_authorizations_account_idx ON account_authorizations (account_id);
CREATE INDEX account_authorizations_authorization_idx ON account_authorizations (authorization_id);
CREATE TRIGGER update_account_authorizations_update_time BEFORE UPDATE ON account_authorizations FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

COMMIT;