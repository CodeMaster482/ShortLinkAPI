CREATE TABLE IF NOT EXISTS link (
    id            BIGSERIAL,
    original_link TEXT UNIQUE,
    token         TEXT UNIQUE,
    expires_at    TIMESTAMPTZ,
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS token_idx
    ON link (token)