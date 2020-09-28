package usersErrors

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

	AllRequired = errors.New("All fields are required")

	InvalidEnter = errors.New("Invalid Enter")

	InValidName = errors.New("Please enter a valid name")

	InvalidUsername = errors.New("Please enter a valid username")

	InvalidAccomplishment = errors.New("Accomplishment is empty")

	InvalidPatronycName = errors.New("Input must be less than 64 characters and can contain _ and .")

	InvalidURL = errors.New("Invalid URL")

	FromTwoToSixtyFour = errors.New("Name must be from 2 to 64 letters")

	FromTwoToHundredTwentyEight = errors.New("Name must be from 2 to 128 letters")

	EightToSixtyFour = errors.New("Name must be from 8 to 64 letters")

	DashAndSpace = errors.New("Name can only contain a-z A-Z, _ and space")

	LettersSymbols = errors.New("Input must only contain letters and special characters")

	InValidUserName = errors.New("Usernamme can be be 6-30 letters, contain numbers, _ and .")

	CommonErrors = errors.New("Input must be less than 128 characters and contain numbers(0-9), special letters and spaces")

	LettersNumbersSpecial = errors.New("Name can contain characters(a-zA-Z), numbers(0-9), and special characters(!@#$%^&*)")

	LettersNumbersSpecialSpace = errors.New("Name can contain characters(a-zA-Z), numbers(0-9), space and special characters(!@#$%^&*)")

	InValidEmail = errors.New("Please enter a valid email")

	InValidLang = errors.New("Please choose valid language")

	InValidPhone = errors.New("Please provide a valid phone number")

	InValidTime = errors.New("Please choose valid time")

	InValidNumber = errors.New("Please choose a valid number")

	TimeRequired = errors.New("Please Enter time")

	SpecificRequired = errors.New("These fileds are required")

	WrongDate = errors.New("Finish date can't be earlier than start date")

	OnlyAlphabets = errors.New("Title of position should only contain alphabets")

	Max32 = errors.New("Text input must be less than 32 characters")

	InvalidGrade = errors.New("Input can only contain +,-,%,.- and space")

	Max64 = errors.New("Text input must be less than 64 characters")

	Max128 = errors.New("Text input must be less than 128 characters")

	Max100 = errors.New("Text input must be less than 100 characters")

	// Max400 = errors.New("Text input must be less than 400 characters")

	Max500 = errors.New("Text input must be less than 500 characters")

	Max1200 = errors.New("Text input must be less than 1200 characters")

	NotLink = errors.New("Link must be type URL")

	MaxArray = errors.New("Too many inputs")

	MaxEachInput = errors.New("Each input must be less than 128 characters")

	EmptyCountry = errors.New("Enter Country")
	EmptyCity    = errors.New("Enter City")
)
