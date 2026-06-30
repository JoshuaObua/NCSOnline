-- Migration: 001_create_users
-- Creates users table with soft-delete support

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id                 TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    email              TEXT NOT NULL UNIQUE,
    password_hash      TEXT NOT NULL,
    first_name         TEXT NOT NULL,
    last_name          TEXT NOT NULL,
    phone              TEXT,
    is_active          BOOLEAN NOT NULL DEFAULT TRUE,
    is_email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at  TIMESTAMPTZ,
    last_login_at      TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
