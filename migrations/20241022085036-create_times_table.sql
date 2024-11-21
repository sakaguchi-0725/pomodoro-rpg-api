-- +migrate Up
CREATE TABLE times (
    id VARCHAR(255) PRIMARY KEY,
    focus_time INT NOT NULL,
    execution_date TIMESTAMPTZ NOT NULL,
    account_id VARCHAR(255),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts(id)
);

-- +migrate Down
DROP TABLE IF EXISTS times;