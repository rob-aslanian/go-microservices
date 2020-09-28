package jobsErrors

import "errors"

var (
	PointerIsNill = errors.New("Pointer cannot be nil")

	InvalidCurrency = errors.New("Invalid Currency")

	InvalidSalary = errors.New("Invalid Salary")

	InvalidJobType = errors.New("Invalid Job type")

	InvalidJobFunctionNumber = errors.New("Invalid job function number")

	InvalidLanguage = errors.New("Invalid Language")

	SpecificRequired = errors.New("These fields are required")

	Max128 = errors.New("Text input must be less than 128 characters")

	Max500 = errors.New("Text input must be less than 500 characters")

	Max2000 = errors.New("Text input must be less than 2000 characters")

	InvalidURL = errors.New("Invalid URL")

	InvalidVAT = errors.New("Invalid VAT number")

	InvalidEmail = errors.New("Invalid Email")

	InValidPhone = errors.New("Please provide a valid phone number")

	InvalidEnter = errors.New("Invalid Enter")

	InvalidDate = errors.New("Invalid Date")
)
