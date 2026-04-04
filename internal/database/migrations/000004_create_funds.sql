-- Migration 004: Create funds table
-- Quỹ là "túi tiền ảo" cho phép người dùng phân bổ tiền theo mục tiêu

CREATE TABLE IF NOT EXISTS funds (
    id            BIGINT        PRIMARY KEY,
    user_id       BIGINT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          VARCHAR(100)  NOT NULL,
    description   TEXT,
    target_amount NUMERIC(18,2) NOT NULL DEFAULT 0,  -- 0 = không đặt mục tiêu
    balance       NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency      VARCHAR(10)   NOT NULL DEFAULT 'VND',
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_funds_user_id ON funds(user_id);
