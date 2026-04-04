package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/models"
	"github.com/maynguyen24/sever/pkg/apperr"
	jwtUtil "github.com/maynguyen24/sever/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

const (
	fidoSessionEnrollPrefix = "fido:session:enroll:"
	fidoSessionAuthPrefix   = "fido:session:auth:"
	fidoSessionTTL          = 5 * time.Minute
)

// FIDODeviceRepository defines the DB contract for the FIDO service
type FIDODeviceRepository interface {
	GetByID(ctx context.Context, id, userID int64) (*models.Device, error)
	GetByCredentialID(ctx context.Context, credentialID string) (*models.Device, error)
	UpdateFIDOCredential(ctx context.Context, deviceID int64, credentialID, publicKey, aaguid string, signCount int64) error
	UpdateSignCount(ctx context.Context, deviceID int64, signCount int64) error
}

// FIDOService handles FIDO2 WebAuthn enrollment and step-up authentication
type FIDOService struct {
	deviceRepo FIDODeviceRepository
	redis      *redis.Client
	webauthn   *webauthn.WebAuthn
	cfg        *configs.Config
}

func NewFIDOService(deviceRepo FIDODeviceRepository, redisClient *redis.Client, cfg *configs.Config) (*FIDOService, error) {
	wauthn, err := webauthn.New(&webauthn.Config{
		RPID:          cfg.WebAuthnRPID,
		RPDisplayName: cfg.WebAuthnRPName,
		RPOrigins:     []string{cfg.WebAuthnOrigin},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize webauthn: %w", err)
	}
	return &FIDOService{
		deviceRepo: deviceRepo,
		redis:      redisClient,
		webauthn:   wauthn,
		cfg:        cfg,
	}, nil
}

// --- Enrollment ---

// BeginEnrollment generates a WebAuthn credential creation challenge for a device.
// Returns the PublicKeyCredentialCreationOptions JSON to send to the client.
func (s *FIDOService) BeginEnrollment(ctx context.Context, userID, deviceID int64) (*protocol.CredentialCreation, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device: %w", err)
	}
	if device == nil {
		return nil, fmt.Errorf("%w: Device not found", apperr.ErrNotFound)
	}

	fidoUser := newFIDOUser(userID, device)

	options, session, err := s.webauthn.BeginRegistration(fidoUser)
	if err != nil {
		return nil, fmt.Errorf("failed to begin registration: %w", err)
	}

	if err := s.storeSession(ctx, fidoSessionEnrollPrefix+strconv.FormatInt(deviceID, 10), session); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	return options, nil
}

// FinishEnrollment verifies the attestation response from the client and stores the credential.
func (s *FIDOService) FinishEnrollment(ctx context.Context, userID, deviceID int64, body []byte) (*models.Device, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device: %w", err)
	}
	if device == nil {
		return nil, fmt.Errorf("%w: Device not found", apperr.ErrNotFound)
	}

	session, err := s.loadSession(ctx, fidoSessionEnrollPrefix+strconv.FormatInt(deviceID, 10))
	if err != nil {
		return nil, fmt.Errorf("%w: Challenge expired, please try again", apperr.ErrInvalidInput)
	}
	// Session is consumed — delete it immediately to prevent replay
	s.redis.Del(ctx, fidoSessionEnrollPrefix+strconv.FormatInt(deviceID, 10))

	fidoUser := newFIDOUser(userID, device)

	parsedResponse, err := protocol.ParseCredentialCreationResponseBytes(body)
	if err != nil {
		return nil, fmt.Errorf("%w: Biometric verification failed", apperr.ErrInvalidInput)
	}

	credential, err := s.webauthn.CreateCredential(fidoUser, *session, parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("%w: Biometric verification failed", apperr.ErrInvalidInput)
	}

	credentialID := base64.RawURLEncoding.EncodeToString(credential.ID)
	publicKey := base64.StdEncoding.EncodeToString(credential.PublicKey)
	aaguid := base64.StdEncoding.EncodeToString(credential.Authenticator.AAGUID)

	if err := s.deviceRepo.UpdateFIDOCredential(ctx, deviceID, credentialID, publicKey, aaguid, int64(credential.Authenticator.SignCount)); err != nil {
		return nil, fmt.Errorf("failed to save credential: %w", err)
	}

	// Reload device to return the updated record
	updated, err := s.deviceRepo.GetByID(ctx, deviceID, userID)
	if err != nil || updated == nil {
		return nil, fmt.Errorf("failed to reload device: %w", err)
	}
	return updated, nil
}

// --- Step-up Authentication ---

// BeginAuthentication generates a WebAuthn assertion challenge for a known credential.
// Returns the PublicKeyCredentialRequestOptions JSON to send to the client.
func (s *FIDOService) BeginAuthentication(ctx context.Context, credentialID string) (*protocol.CredentialAssertion, error) {
	device, err := s.deviceRepo.GetByCredentialID(ctx, credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device: %w", err)
	}
	if device == nil || device.FIDOCredentialID == nil {
		return nil, fmt.Errorf("%w: Credential not found", apperr.ErrNotFound)
	}

	fidoUser := newFIDOUser(device.UserID, device)

	options, session, err := s.webauthn.BeginLogin(fidoUser)
	if err != nil {
		return nil, fmt.Errorf("failed to begin login: %w", err)
	}

	if err := s.storeSession(ctx, fidoSessionAuthPrefix+credentialID, session); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	return options, nil
}

// FinishAuthentication verifies the assertion and returns a FIDO-verified access token.
func (s *FIDOService) FinishAuthentication(ctx context.Context, body []byte) (string, error) {
	// Extract credential ID from the raw body to look up the device + session
	var raw struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &raw); err != nil || raw.ID == "" {
		return "", fmt.Errorf("%w: Invalid request body", apperr.ErrInvalidInput)
	}

	device, err := s.deviceRepo.GetByCredentialID(ctx, raw.ID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch device: %w", err)
	}
	if device == nil || device.FIDOCredentialID == nil {
		return "", fmt.Errorf("%w: Credential not found", apperr.ErrNotFound)
	}

	session, err := s.loadSession(ctx, fidoSessionAuthPrefix+raw.ID)
	if err != nil {
		return "", fmt.Errorf("%w: Challenge expired, please try again", apperr.ErrInvalidInput)
	}
	// Session is consumed — delete it immediately to prevent replay
	s.redis.Del(ctx, fidoSessionAuthPrefix+raw.ID)

	fidoUser := newFIDOUser(device.UserID, device)

	parsedResponse, err := protocol.ParseCredentialRequestResponseBytes(body)
	if err != nil {
		return "", fmt.Errorf("%w: Biometric verification failed", apperr.ErrInvalidInput)
	}

	credential, err := s.webauthn.ValidateLogin(fidoUser, *session, parsedResponse)
	if err != nil {
		return "", fmt.Errorf("%w: Biometric verification failed", apperr.ErrUnauthorized)
	}

	// Replay protection: sign count must advance
	if credential.Authenticator.CloneWarning {
		return "", fmt.Errorf("%w: Authenticator replay detected", apperr.ErrUnauthorized)
	}

	if err := s.deviceRepo.UpdateSignCount(ctx, device.ID, int64(credential.Authenticator.SignCount)); err != nil {
		return "", fmt.Errorf("failed to update sign count: %w", err)
	}

	accessToken, err := jwtUtil.GenerateAccessTokenFIDO(device.UserID, s.cfg)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return accessToken, nil
}

// --- Helpers ---

func (s *FIDOService) storeSession(ctx context.Context, key string, session *webauthn.SessionData) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, data, fidoSessionTTL).Err()
}

func (s *FIDOService) loadSession(ctx context.Context, key string) (*webauthn.SessionData, error) {
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var session webauthn.SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// fidoUser adapts a device record to the webauthn.User interface
type fidoUser struct {
	id          int64
	credentials []webauthn.Credential
}

func newFIDOUser(userID int64, device *models.Device) *fidoUser {
	u := &fidoUser{id: userID}
	if device.FIDOCredentialID != nil && device.FIDOPublicKey != nil {
		credID, err1 := base64.RawURLEncoding.DecodeString(*device.FIDOCredentialID)
		pubKey, err2 := base64.StdEncoding.DecodeString(*device.FIDOPublicKey)
		if err1 == nil && err2 == nil {
			u.credentials = []webauthn.Credential{
				{
					ID:        credID,
					PublicKey: pubKey,
					Authenticator: webauthn.Authenticator{
						SignCount: uint32(device.FIDOSignCount),
					},
				},
			}
		}
	}
	return u
}

func (u *fidoUser) WebAuthnID() []byte {
	return []byte(strconv.FormatInt(u.id, 10))
}

func (u *fidoUser) WebAuthnName() string {
	return strconv.FormatInt(u.id, 10)
}

func (u *fidoUser) WebAuthnDisplayName() string {
	return strconv.FormatInt(u.id, 10)
}

func (u *fidoUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}
