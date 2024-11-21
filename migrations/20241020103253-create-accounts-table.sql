-- +migrate Up
CREATE TABLE accounts (
    id VARCHAR(255) PRIMARY KEY,
    cognito_uid VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    image VARCHAR(255),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE IF EXISTS accounts;
