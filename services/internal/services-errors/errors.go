package serviceErrors

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

	InValidEmail = errors.New("Please enter a valid email")

	Max128 = errors.New("Text input must be less than 128 characters")

	Max400 = errors.New("Text input must be less than 400 characters")

	Max500 = errors.New("Text input must be less than 500 characters")

	Max600 = errors.New("Text input must be less than 600 characters")

	Max2000 = errors.New("Text input must be less than 2000 characters")

	InvalidURL = errors.New("Invalid URL")

	InvalidVAT = errors.New("Invalid VAT number")

	InvalidEmail = errors.New("Invalid Email")

	InValidPhone = errors.New("Please provide a valid phone number")

	InvalidEnter = errors.New("Invalid Enter")

	InvalidDate = errors.New("Invalid Date")

	InValidTime = errors.New("Please choose valid time")

	FromTwoToSixtyFour = errors.New("Name must be from 2 to 64 letters")

	FromTwoToHundredTwentyEight = errors.New("Name must be from 2 to 128 letters")

	EightToSixtyFour = errors.New("Name must be from 8 to 64 letters")

	DashAndSpace = errors.New("Name can only contain a-z A-Z, - and space")

	AlphabetNumberSymbol = errors.New("Input can only contain Alphabets(a-zA-Z), numbers(0-9) and symbols")

	InValidUserName = errors.New("Usernamme can be be 6-30 letters, contain numbers, _ and .")

	CommonErrors = errors.New("Input must be less than 128 characters and contain numbers(0-9), special letters and spaces")

	LettersNumbersSpecial = errors.New("Name can contain characters(a-zA-Z), numbers(0-9), and special characters(!@#$%^&*)")

	LettersNumbersSpecialSpace = errors.New("Name can contain characters(a-zA-Z), numbers(0-9), space and special characters(!@#$%^&*)")
)
