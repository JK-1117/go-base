-- +goose Up

CREATE TABLE login_session (
    session_id VARCHAR(64) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    last_login TIMESTAMP NOT NULL,
    ip_addr INET,
    user_agent TEXT,
    expired_at TIMESTAMP NOT NULL
);
CREATE INDEX index_login_session_session_id ON login_session USING HASH (session_id);
-- +goose Down

DROP TABLE login_session;