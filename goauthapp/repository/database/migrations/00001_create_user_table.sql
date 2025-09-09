-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at
= NOW();
RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    phone_number VARCHAR(255) UNIQUE NOT NULL,
    created_at   TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_phone_number ON users (phone_number);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON users
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();


-- +migrate Down

DROP TRIGGER IF EXISTS set_timestamp ON users;
DROP INDEX IF EXISTS idx_users_phone_number;
DROP TABLE IF EXISTS users;

-- +migrate StatementBegin
DROP FUNCTION IF EXISTS trigger_set_timestamp();
-- +migrate StatementEnd
