
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


BEGIN;


CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    "name" TEXT UNIQUE,
    permission_bit BIGINT NOT NULL CHECK (permission_bit >= 0),
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account_roles (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT REFERENCES roles (id) NOT NULL,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (role_id, account_id)
);

CREATE TABLE confirmations (
    id BIGSERIAL PRIMARY KEY,
    "type" NUMERIC NOT NULL, -- email, phone number or password reset
    confirmation_target VARCHAR DEFAULT NULL,
    account_id BIGINT REFERENCES accounts (id) NOT NULL,
    key VARCHAR NOT NULL UNIQUE,
    confirm_time TIMESTAMPTZ DEFAULT NULL,
    failed_confirmations_count INT DEFAULT 0,
    expire_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + interval '5 hours',
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