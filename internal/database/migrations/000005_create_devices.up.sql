-- Migration 005: Create devices table
-- Supports: one-device-one-account enforcement and FIDO2/WebAuthn biometric auth
--
-- Design decisions:
--   • device_fingerprint is GLOBALLY UNIQUE → prevents one device from being
--     linked to multiple accounts (one device = one account rule).
--   • FIDO2 fields are nullable → only populated when user enrolls biometrics.
--   • fido_sign_count acts as a monotonic counter to prevent FIDO replay attacks.
--   • push_token is stored here so deposit/transaction notifications can be
--     delivered to the correct device-specific FCM / APNs endpoint.

CREATE TABLE IF NOT EXISTS devices (
    id                  BIGINT          PRIMARY KEY,
    user_id             BIGINT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Device identity
    device_fingerprint  VARCHAR(512)    NOT NULL UNIQUE, -- hardware/browser fingerprint (globally unique)
    device_name         VARCHAR(150),                    -- e.g. "iPhone 15 Pro", "Chrome on macOS"
    platform            VARCHAR(20)     NOT NULL         -- 'ios' | 'android' | 'web'
                            CHECK (platform IN ('ios', 'android', 'web')),

    -- Push notification token (FCM / APNs)
    push_token          TEXT,

    -- FIDO2 / WebAuthn biometric fields (nullable until biometric is enrolled)
    fido_credential_id  TEXT            UNIQUE,          -- credential ID returned by authenticator
    fido_public_key     TEXT,                            -- COSE-encoded public key
    fido_sign_count     BIGINT          NOT NULL DEFAULT 0, -- replay-attack monotonic counter
    fido_aaguid         VARCHAR(64),                     -- authenticator AAGUID (device type hint)

    -- Status
    is_trusted          BOOLEAN         NOT NULL DEFAULT FALSE, -- user explicitly trusted this device
    is_active           BOOLEAN         NOT NULL DEFAULT TRUE,  -- soft-disable without deleting
    last_used_at        TIMESTAMPTZ,

    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_devices_user_id        ON devices(user_id);
CREATE INDEX IF NOT EXISTS idx_devices_fingerprint    ON devices(device_fingerprint);
CREATE INDEX IF NOT EXISTS idx_devices_fido_cred      ON devices(fido_credential_id) WHERE fido_credential_id IS NOT NULL;
