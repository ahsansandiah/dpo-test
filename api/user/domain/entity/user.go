package userDomainEntity

import (
	"time"

	errorHelper "github.com/ahsansandiah/dpo-test/helpers/error"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	PasswordHash    []byte `json:"password_hash"`
	Email           string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiryTime   *time.Time `json:"expiry_time"`
}

func (r *UserRequest) Validate() error {
	if r.Username == "" {
		return errorHelper.ErrorUsernameIsRequired
	}

	if r.Password == "" {
		return errorHelper.ErrorPasswordIsRequired
	}

	if r.PasswordConfirm == "" {
		return errorHelper.ErrorPasswordConfirmIsRequired
	}

	if r.Password != r.PasswordConfirm {
		return errorHelper.ErrorPasswordNotMatch
	}

	if r.Email == "" {
		return errorHelper.ErrorEmailIsRequired
	}

	return nil
}

func (r *LoginRequest) LoginValidate() error {
	if r.Username == "" {
		return errorHelper.ErrorUsernameIsRequired
	}

	if r.Password == "" {
		return errorHelper.ErrorPasswordIsRequired
	}

	return nil
}
