package customerDomainEntity

import (
	"time"

	errorHelper "github.com/ahsansandiah/dpo-test/helpers/error"
)

type Customer struct {
	ID          int64     `json:"id"`
	FullName    string    `json:"full_name"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CustomerRequest struct {
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type CustomerFilter struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
	LIMIT       string `json:"limit"`
}

func (r *CustomerRequest) Validate() error {
	if r.FullName == "" {
		return errorHelper.ErrorFullNameIsRequired
	}

	if r.Address == "" {
		return errorHelper.ErrorAddressIsRequired
	}

	if r.PhoneNumber == "" {
		return errorHelper.ErrorPhoneNumberIsRequired
	}

	if r.Email == "" {
		return errorHelper.ErrorEmailIsRequired
	}

	return nil
}
