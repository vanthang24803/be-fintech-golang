package models

import "time"

// DevicePlatform defines the supported client platforms
type DevicePlatform string

const (
	PlatformIOS     DevicePlatform = "ios"
	PlatformAndroid DevicePlatform = "android"
	PlatformWeb     DevicePlatform = "web"
)

// Device represents a registered client device for a user.
// Used for one-device-one-account enforcement and FIDO2 biometric auth.
type Device struct {
	ID                int64          `db:"id"                 json:"id,string"`
	UserID            int64          `db:"user_id"            json:"user_id,string"`

	// Identity
	DeviceFingerprint string         `db:"device_fingerprint" json:"device_fingerprint"`
	DeviceName        *string        `db:"device_name"        json:"device_name,omitempty"`
	Platform          DevicePlatform `db:"platform"           json:"platform"`

	// Push notification
	PushToken         *string        `db:"push_token"         json:"push_token,omitempty"`

	// FIDO2 / WebAuthn biometric (nil until biometric is enrolled)
	FIDOCredentialID  *string        `db:"fido_credential_id" json:"-"`             // kept server-side only
	FIDOPublicKey     *string        `db:"fido_public_key"    json:"-"`             // kept server-side only
	FIDOSignCount     int64          `db:"fido_sign_count"    json:"-"`
	FIDOAaguid        *string        `db:"fido_aaguid"        json:"fido_aaguid,omitempty"`

	// Status
	IsTrusted         bool           `db:"is_trusted"         json:"is_trusted"`
	IsActive          bool           `db:"is_active"          json:"is_active"`
	LastUsedAt        *time.Time     `db:"last_used_at"       json:"last_used_at,omitempty"`

	CreatedAt         time.Time      `db:"created_at"         json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"         json:"updated_at"`
}

// RegisterDeviceRequest is sent by the client on first launch / login
type RegisterDeviceRequest struct {
	DeviceFingerprint string         `json:"device_fingerprint" validate:"required,min=16,max=512"`
	DeviceName        *string        `json:"device_name"`
	Platform          DevicePlatform `json:"platform"           validate:"required,oneof=ios android web"`
	PushToken         *string        `json:"push_token"`
}

// UpdatePushTokenRequest allows a device to refresh its FCM/APNs token
type UpdatePushTokenRequest struct {
	PushToken string `json:"push_token" validate:"required"`
}

// EnrollBiometricRequest carries FIDO2 registration data from the authenticator
type EnrollBiometricRequest struct {
	CredentialID string `json:"credential_id" validate:"required"`
	PublicKey    string `json:"public_key"    validate:"required"` // COSE-encoded
	Aaguid       string `json:"aaguid"`
	SignCount     int64  `json:"sign_count"`
}

// BiometricAuthRequest carries FIDO2 assertion data for login
type BiometricAuthRequest struct {
	CredentialID string `json:"credential_id" validate:"required"`
	SignCount     int64  `json:"sign_count"    validate:"required"`
	// Signature verification is typically done via a FIDO2 library;
	// raw fields passed through from the WebAuthn assertion response.
	AuthenticatorData string `json:"authenticator_data" validate:"required"`
	ClientDataJSON    string `json:"client_data_json"   validate:"required"`
	Signature         string `json:"signature"          validate:"required"`
}
