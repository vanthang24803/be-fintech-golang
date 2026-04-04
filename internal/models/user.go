package models

import "time"

type User struct {
	ID           int64     `db:"id" json:"id,string"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	GoogleID     *string   `db:"google_id" json:"google_id,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	User *User `json:"user"`
}

type Profile struct {
	ID          int64      `db:"id" json:"id,string"`
	UserID      int64      `db:"user_id" json:"user_id,string"`
	FullName    *string    `db:"full_name" json:"full_name"`
	AvatarURL   *string    `db:"avatar_url" json:"avatar_url"`
	PhoneNumber *string    `db:"phone_number" json:"phone_number"`
	DateOfBirth *time.Time `db:"date_of_birth" json:"date_of_birth"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

type ProfileResponse struct {
	User    *User    `json:"user"`
	Profile *Profile `json:"profile"`
}

// UpdateProfileRequest is the payload for updating a user's profile
type UpdateProfileRequest struct {
	FullName    *string    `json:"full_name"    validate:"omitempty,max=100"`
	AvatarURL   *string    `json:"avatar_url"`
	PhoneNumber *string    `json:"phone_number" validate:"omitempty,e164"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}
