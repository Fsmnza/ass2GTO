CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    fname         VARCHAR(255),
    sname         VARCHAR(255),
    email         citext UNIQUE               NOT NULL,
    role          VARCHAR(50),
    password_hash BYTEA                       NOT NULL,
    activated     BOOLEAN                     NOT NULL,
    version       INTEGER                     NOT NULL DEFAULT 1
);
