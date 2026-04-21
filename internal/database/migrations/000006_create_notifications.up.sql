-- Migration 006: Create notifications table
-- Supports: multi-source notifications (transaction, fund, auth, system, etc.)

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_source') THEN
        CREATE TYPE notification_source AS ENUM (
            'transaction',
            'fund',
            'auth',
            'system'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_type') THEN
        CREATE TYPE notification_type AS ENUM (
            'info',
            'success',
            'warning',
            'alert'
        );
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS notifications (
    id          BIGINT                  PRIMARY KEY,
    user_id     BIGINT                  NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Source of the notification
    source      notification_source     NOT NULL DEFAULT 'system',
    source_id   BIGINT,                 -- optional FK to the triggering entity (transaction_id, fund_id…)

    -- Content
    type        notification_type       NOT NULL DEFAULT 'info',
    title       VARCHAR(200)            NOT NULL,
    body        TEXT                    NOT NULL,

    -- Flexible extra data per source (e.g. {"amount": 50000, "currency": "VND"})
    metadata    JSONB,

    -- Read state
    is_read     BOOLEAN                 NOT NULL DEFAULT FALSE,
    read_at     TIMESTAMPTZ,

    created_at  TIMESTAMPTZ             NOT NULL DEFAULT NOW()
);

-- Fast lookup of all notifications for a user
CREATE INDEX IF NOT EXISTS idx_notifications_user_id
    ON notifications(user_id, created_at DESC);

-- Fast unread count query: SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND is_read=false
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread
    ON notifications(user_id) WHERE is_read = FALSE;

-- Lookup by source entity (e.g. all notifications related to fund_id=123)
CREATE INDEX IF NOT EXISTS idx_notifications_source
    ON notifications(source, source_id) WHERE source_id IS NOT NULL;
