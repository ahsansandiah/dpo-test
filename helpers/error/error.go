package errorHelper

import "errors"

var (
	ErrorDataNotfound = errors.New("data not found")

	// Error customer module
	ErrorFullNameIsRequired    = errors.New("full name is required")
	ErrorAddressIsRequired     = errors.New("address is required")
	ErrorPhoneNumberIsRequired = errors.New("phone number is required")
	ErrorEmailIsRequired       = errors.New("email is required")

	// Error order module
	ErrorCustomerIdRequired   = errors.New("customer is required")
	ErrorOrderDateRequired    = errors.New("order date is required")
	ErrorAmountIsRequired     = errors.New("amount is required")
	ErrorOrderItemsIsRequired = errors.New("order items is required")

	// Error user module
	ErrorUsernameIsRequired        = errors.New("User name is required")
	ErrorPasswordIsRequired        = errors.New("Password is required")
	ErrorPasswordConfirmIsRequired = errors.New("Password confirm is required")
	ErrorPasswordNotMatch          = errors.New("Password not match")
)
