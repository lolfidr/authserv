DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS apps CASCADE;

CREATE TABLE users (
    id        BIGSERIAL PRIMARY KEY,  
    email     TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE apps (
    id        BIGSERIAL PRIMARY KEY,
    name      TEXT NOT NULL UNIQUE,
    secret    TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);