package searchErrors

import "errors"

var (
	// ErrNotAuthenticated ...
	ErrNotAuthenticated = errors.New("not_authenticated")

	// ErrAlreadyExists ...
	ErrAlreadyExists = errors.New("already_exists")

	// ErrNotFound ...
	ErrNotFound = errors.New("not_found")

	// ErrWrongArgument ...
	ErrWrongArgument = errors.New("wrong_argument")

	// ErrInternalError ...
	ErrInternalError = errors.New("internal_error")

	// ErrWrongID  ...
	ErrWrongID = errors.New("wrong_id")

	SpecificRequired = errors.New("These fields are required")

	Max100 = errors.New("Text input must be less than 100 characters")

	Max400 = errors.New("Text input must be less than 400 characters")

	Max500 = errors.New("Text input must be less than 500 characters")

	Max2000 = errors.New("Text input must be less than 2000 characters")

	InvalidURL = errors.New("Invalid URL")

	InvalidVAT = errors.New("Invalid VAT number")

	InvalidEmail = errors.New("Invalid Email")

	InvalidSalary = errors.New("Invalid Salary")

	InValidPhone = errors.New("Please provide a valid phone number")

	InvalidName = errors.New("This field should only contain letters")

	InvalidEnter = errors.New("Invalid Enter")

	InvalidDate = errors.New("Invalid Date")

	InvalidAge = errors.New("Please enter correct age")

	InvalidCountry = errors.New("Invalid Country Enter")

	InvalidCity = errors.New("Invalid City Enter")
)
