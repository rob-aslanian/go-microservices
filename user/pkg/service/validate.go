package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
	userReport "gitlab.lan/Rightnao-site/microservices/user/pkg/user_report"

	"github.com/asaskevich/govalidator"
	"github.com/ttacon/libphonenumber"
	usersErrors "gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

const (
	firstNamePattern                  = `^[a-zA-Z\- ]+$`
	numberPattern                     = `^[0-9]+`
	unicodeNamePattern                = `^[\p{L}'-]+`
	specailPattern                    = "[\\[\\]\\{\\}\\|\\<\\>\\(\\)\\'~!@#$%^&*?.,/\\-+]"
	usernamePattern                   = `^[a-zA-Z0-9_\.]{6,30}$`
	urlPattern                        = `^[a-z0-9\.]+$`
	gradePattern                      = `^[a-zA-z0-9-+%.\- ]+`
	commonPattern                     = `^[a-zA-Z0-9\W ]{0,128}$`
	middlenamePattern                 = `^[a-zA-Z\- ]{0,64}$`
	lettersSymbolsPattern             = `^[a-zA-Z\W ]+$`
	lettersNumbersSpecialPattern      = `^[a-zA-Z0-9\W ]+$`
	lettersNumbersSpecialSpacePattern = `^[a-zA-Z0-9\W ]+$`
)

func fromTwoToSixtyFour(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		if utf8.RuneCountInString(str[i]) > 64 || utf8.RuneCountInString(str[i]) < 2 {
			return usersErrors.FromTwoToSixtyFour
		}
	}

	return nil
}
func eightToSixtyFour(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		if utf8.RuneCountInString(str[i]) > 64 || utf8.RuneCountInString(str[i]) < 8 {
			return usersErrors.EightToSixtyFour
		}
	}

	return nil
}

func fromTwoToHundredTwentyEight(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		if utf8.RuneCountInString(str[i]) > 128 || utf8.RuneCountInString(str[i]) < 2 {
			return usersErrors.FromTwoToHundredTwentyEight
		}
	}

	return nil
}

// check email
func emailValidator(email string) error {
	if govalidator.IsEmail(email) == false {
		return usersErrors.InValidEmail
	}
	if utf8.RuneCountInString(email) > 128 {
		return usersErrors.Max128
	}
	if email == "" {
		return usersErrors.SpecificRequired
	}

	return nil
}

func userNameValidator(str ...string) error {
	for _, t := range str {

		if govalidator.StringMatches(t, usernamePattern) == false {
			return usersErrors.InValidUserName
		}
	}

	return nil
}

func dashAndSpace(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		if govalidator.StringMatches(str[i], firstNamePattern) == false {
			return usersErrors.DashAndSpace
		}
	}
	return nil
}

func lettersSymbols(str ...string) error {
	for _, t := range str {
		strings.TrimSpace(t)
		if govalidator.StringMatches(t, lettersSymbolsPattern) == false {
			return usersErrors.LettersSymbols
		}
	}
	return nil
}

// //middlename not be over 100 characters. must contain alphabets, numbers, dot and underscore
func middlenicknameValidator(middlename *string) error {
	if middlename != nil && *middlename != "" {
		strings.TrimSpace(*middlename)
		if govalidator.StringMatches(*middlename, middlenamePattern) == false {
			return usersErrors.InvalidPatronycName
		}
	}

	return nil
}
func changeNameOnNativeLanguageValidator(name *string, lang *string) error {
	if name == nil || lang == nil {
		return usersErrors.InvalidEnter
	}
	strings.TrimSpace(*name)
	strings.TrimSpace(*lang)
	if *name == "" {
		return usersErrors.SpecificRequired
	}
	if govalidator.StringMatches(*name, unicodeNamePattern) == false {
		return usersErrors.InValidName
	}
	if utf8.RuneCountInString(*name) > 128 {
		return usersErrors.Max128
	}
	if utf8.RuneCountInString(*lang) != 2 {
		return usersErrors.InValidLang
	}

	return nil
}

func birthdayValidator(birthday *time.Time) error {
	ti, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

	if birthday != nil {
		if birthday.Before(ti) || birthday.After(time.Now()) {
			return usersErrors.InValidTime
		}
	}

	return nil
}

func phoneValidator(countryCode string, number string, countryID string) error {
	if number == "" {
		return usersErrors.SpecificRequired
	}

	if govalidator.StringMatches(number, numberPattern) {

		number = strings.Replace(number, " ", "", -1)
		num, err := strconv.Atoi(number)
		if err != nil {
			return err
		}
		rCode, err := strconv.Atoi(countryCode)
		if err != nil {
			return err
		}
		regionCode := int32(rCode)

		intNumber := uint64(num)
		n := libphonenumber.PhoneNumber{
			CountryCode:    &regionCode,
			NationalNumber: &intNumber,
		}

		strings.TrimSpace(number)
		//checking for these countries to match with their country phone numbers. Armenia doesn't work with ValidNumberForRegion so it is filtered
		//down with isValidNumber. It gives the right answer.
		if countryID == "GE" || countryID == "RU" || countryID == "US" || countryID == "UK" || countryID == "KZ" || countryID == "AZ" {
			if libphonenumber.IsValidNumberForRegion(&n, countryID) == false {
				return usersErrors.InValidPhone
			}
		} else {
			if libphonenumber.IsValidNumber(&n) == false {
				return usersErrors.InValidPhone
			}
		}
	} else {
		return usersErrors.InValidPhone
	}
	return nil
}

func addressValidator(address *account.MyAddress) error {
	if address != nil {
		strings.TrimSpace(address.Apartment)
		strings.TrimSpace(address.Street)
		if address.Location.Country == nil {
			return usersErrors.EmptyCountry
		}
		if address.Location.City == nil {
			return usersErrors.EmptyCity
		}
		if address.Location.Country != nil {
			strings.TrimSpace(address.Location.Country.ID)
		}
		if address.Location.City != nil {
			strings.TrimSpace(address.Location.City.Name)
		}
		err := length128Validator(address.Apartment, address.Street, address.ZIP)
		if err != nil {
			return err
		}
		if address.Apartment == "" || address.Street == "" || address.ZIP == "" {
			return usersErrors.SpecificRequired
		}
		err = emptyValidator(address.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// //skills array not be over 100 in length or empty. each element not to be over 100 characters
func skillsValidator(skills []*profile.Skill) error {
	if skills != nil {
		if len(skills) > 50 {
			return usersErrors.MaxArray
		}

		if len(skills) == 0 {
			return usersErrors.SpecificRequired
		}

		for i := range skills {
			if skills[i] != nil {
				strings.TrimSpace(skills[i].Skill)
				if utf8.RuneCountInString(skills[i].Skill) > 64 {
					return usersErrors.Max64
				}
			}
		}
	} else {
		return usersErrors.SpecificRequired
	}

	return nil
}

func commonPatternValidator(str ...string) error {
	for _, s := range str {
		strings.TrimSpace(s)
		if govalidator.StringMatches(s, commonPattern) == false {
			return usersErrors.CommonErrors
		}
	}
	return nil
}

func validPassword(str string) error {
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"numeric":    {unicode.Number, unicode.Digit},
		// "special":    {unicode.Symbol, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range str {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("Password must have at least one %s character", name)
	}

	return nil
}

// //control language input not to be empty as it is mendotary
func languageValidator(prof *profile.KnownLanguage) error {
	if prof == nil {
		return usersErrors.InvalidEnter
	}
	strings.TrimSpace(prof.Language)
	if prof.Language == "" {
		return usersErrors.SpecificRequired
	}
	if prof.Rank > 5 || prof.Rank < 0 {
		return usersErrors.InvalidEnter
	}

	return nil
}

func headlineValidator(headline string) error {
	strings.TrimSpace(headline)
	if utf8.RuneCountInString(headline) > 128 {
		return usersErrors.Max128
	}

	return nil
}

func storyValidator(story string) error {
	strings.TrimSpace(story)
	if utf8.RuneCountInString(story) > 1200 {
		return usersErrors.Max1200
	}

	return nil
}

// //report description not to be over 500 characters
func reportValidator(report *userReport.Report) error {
	if report != nil {
		strings.TrimSpace(report.Description)
		if utf8.RuneCountInString(report.Description) > 500 {
			return usersErrors.Max500
		}
	}

	return nil
}

func accomplishmentValidator(accomplishment *profile.Accomplishment) error {
	if accomplishment == nil {
		return usersErrors.InvalidAccomplishment
	}
	if accomplishment.Name != "" {
		switch accomplishment.Type {
		case profile.AccomplishmentTypeCertificate:

			err := length128Validator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = emptyValidator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = pointerEmptyValidator(accomplishment.Issuer)
			if err != nil {
				return err
			}
			err = lettersNumbersSpecialValidator(accomplishment.Name)
			if err != nil {
				return err
			}

			if accomplishment.LicenseNumber != nil && govalidator.StringMatches(*accomplishment.LicenseNumber, commonPattern) == false {
				return usersErrors.Max128
			}

			if accomplishment.Issuer != nil && utf8.RuneCountInString(*accomplishment.Issuer) > 128 {
				return usersErrors.Max128
			}

			// time, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

			// if accomplishment.FinishDate != nil {
			// 	if accomplishment.FinishDate.Before(time) {
			// 		return usersErrors.InValidTime
			// 	}
			// }

			if accomplishment.StartDate == nil {
				return errors.New("Please Enter Start Date")
			}

			//if user checked not expires checkbox, then we don't need expire date

			if accomplishment.IsExpire != nil && *accomplishment.IsExpire {
				accomplishment.FinishDate = nil
			}

			if accomplishment.FinishDate != nil && !accomplishment.FinishDate.IsZero() && accomplishment.StartDate.After(*accomplishment.FinishDate) {
				return usersErrors.WrongDate
			}

			if accomplishment.Description != nil {
				strings.TrimSpace(*accomplishment.Description)
				if *accomplishment.Description != "" {
					err = length500Validator(*accomplishment.Description)
					if err != nil {
						return err
					}
					if govalidator.StringMatches(*accomplishment.Description, lettersNumbersSpecialSpacePattern) == false {
						return usersErrors.LettersNumbersSpecialSpace
					}
				}
			}

			if accomplishment.URL != nil {
				strings.TrimSpace(*accomplishment.URL)
				err = pointerLenght128Validator(accomplishment.URL)
				if err != nil {
					return err
				}
			}

		case profile.AccomplishmentTypeAward:

			err := length128Validator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = emptyValidator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = pointerEmptyValidator(accomplishment.Issuer)
			if err != nil {
				return err
			}
			err = lettersNumbersSpecialValidator(accomplishment.Name)
			if err != nil {
				return err
			}

			if accomplishment.FinishDate == nil {
				return errors.New("Please Enter Date")
			}

			if accomplishment.LicenseNumber != nil && govalidator.StringMatches(*accomplishment.LicenseNumber, commonPattern) == false {
				return usersErrors.Max128
			}

			if accomplishment.Issuer != nil && utf8.RuneCountInString(*accomplishment.Issuer) > 128 {
				return usersErrors.Max128
			}

			// time, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

			// if accomplishment.FinishDate != nil {
			// 	if accomplishment.FinishDate.Before(time) {
			// 		return usersErrors.InValidTime
			// 	}
			// }

			if accomplishment.Description != nil {
				strings.TrimSpace(*accomplishment.Description)
				if *accomplishment.Description != "" {
					err = length500Validator(*accomplishment.Description)
					if err != nil {
						return err
					}
					if govalidator.StringMatches(*accomplishment.Description, lettersNumbersSpecialSpacePattern) == false {
						return usersErrors.LettersNumbersSpecialSpace
					}
				}
			}

		case profile.AccomplishmentTypeProject:

			err := length128Validator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = emptyValidator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = lettersNumbersSpecialValidator(accomplishment.Name)
			if err != nil {
				return err
			}

			if accomplishment.LicenseNumber != nil && govalidator.StringMatches(*accomplishment.LicenseNumber, commonPattern) == false {
				return usersErrors.Max128
			}

			// time, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

			// if accomplishment.FinishDate != nil {
			// 	if accomplishment.FinishDate.Before(time) {
			// 		return usersErrors.InValidTime
			// 	}
			// }

			if accomplishment.StartDate == nil {
				return errors.New("Please Enter Start Date")
			}

			if accomplishment.FinishDate != nil && !accomplishment.FinishDate.IsZero() && accomplishment.StartDate.After(*accomplishment.FinishDate) {
				return usersErrors.WrongDate
			}

			if accomplishment.URL != nil {
				strings.TrimSpace(*accomplishment.URL)
				err = pointerLenght128Validator(accomplishment.URL)
				if err != nil {
					return err
				}
			}

			if accomplishment.Description != nil {
				strings.TrimSpace(*accomplishment.Description)
				if *accomplishment.Description != "" {
					err = length500Validator(*accomplishment.Description)
					if err != nil {
						return err
					}
					if govalidator.StringMatches(*accomplishment.Description, lettersNumbersSpecialSpacePattern) == false {
						return usersErrors.LettersNumbersSpecialSpace
					}
				}
			}
		case profile.AccomplishmentTypeLicense:

			strings.TrimSpace(accomplishment.Name)
			strings.TrimSpace(*accomplishment.LicenseNumber)
			strings.TrimSpace(*accomplishment.Issuer)

			err := length128Validator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = emptyValidator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = pointerEmptyValidator(accomplishment.Issuer)
			if err != nil {
				return err
			}
			err = lettersNumbersSpecialValidator(accomplishment.Name)
			if err != nil {
				return err
			}

			if accomplishment.LicenseNumber != nil && govalidator.StringMatches(*accomplishment.LicenseNumber, commonPattern) == false {
				return usersErrors.Max128
			}

			if accomplishment.Issuer != nil && utf8.RuneCountInString(*accomplishment.Issuer) > 128 {
				return usersErrors.Max128
			}

			// time, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

			// if accomplishment.FinishDate != nil {
			// 	if accomplishment.FinishDate.Before(time) {
			// 		return usersErrors.InValidTime
			// 	}
			// }

			if accomplishment.StartDate == nil {
				return errors.New("Please Enter Start Date")
			}

			if accomplishment.FinishDate != nil && !accomplishment.FinishDate.IsZero() && accomplishment.StartDate.After(*accomplishment.FinishDate) {
				return usersErrors.WrongDate
			}

			//if user checked not expires checkbox, then we don't need expire date

			if accomplishment.IsExpire != nil && *accomplishment.IsExpire {
				accomplishment.FinishDate = nil
			}

		case profile.AccomplishmentTypePublication:

			err := length128Validator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = emptyValidator(accomplishment.Name)
			if err != nil {
				return err
			}
			err = lettersNumbersSpecialValidator(accomplishment.Name)
			if err != nil {
				return err
			}

			if accomplishment.LicenseNumber != nil && govalidator.StringMatches(*accomplishment.LicenseNumber, commonPattern) == false {
				return usersErrors.Max128
			}

			if accomplishment.Issuer != nil && utf8.RuneCountInString(*accomplishment.Issuer) > 128 {
				return usersErrors.Max128
			}

			// time, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

			// if accomplishment.FinishDate != nil {
			// 	if accomplishment.FinishDate.Before(time) {
			// 		return usersErrors.InValidTime
			// 	}
			// }

			if accomplishment.FinishDate == nil {
				return errors.New("Please Enter Date")
			}

			if accomplishment.FinishDate == nil {
				return usersErrors.TimeRequired
			}

			if accomplishment.Description != nil {
				strings.TrimSpace(*accomplishment.Description)
				if *accomplishment.Description != "" {
					err = length500Validator(*accomplishment.Description)
					if err != nil {
						return err
					}
					if govalidator.StringMatches(*accomplishment.Description, lettersNumbersSpecialSpacePattern) == false {
						return usersErrors.LettersNumbersSpecialSpace
					}
				}
			}

			if accomplishment.URL != nil {
				strings.TrimSpace(*accomplishment.URL)
				err = pointerLenght128Validator(accomplishment.URL)
				if err != nil {
					return err
				}
			}

		case profile.AccomplishmentTypeTest:
			if accomplishment != nil && accomplishment.Score != nil {

				score := fmt.Sprintf("%f", *accomplishment.Score)

				if len(score) > 16 {
					return errors.New("Score cannot be over 16 digits")
				}

				// time, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

				// if accomplishment.FinishDate != nil {
				// 	if accomplishment.FinishDate.Before(time) {
				// 		return usersErrors.InValidTime
				// 	}
				// }

				if accomplishment.Description != nil {
					strings.TrimSpace(*accomplishment.Description)
					if *accomplishment.Description != "" {
						err := length500Validator(*accomplishment.Description)
						if err != nil {
							return err
						}
						if govalidator.StringMatches(*accomplishment.Description, lettersNumbersSpecialSpacePattern) == false {
							return usersErrors.LettersNumbersSpecialSpace
						}
					}
				}
			}
			if accomplishment.Score == nil {
				return errors.New("Please fill in the score")
			}
			if accomplishment.Name == "" {
				return errors.New("Please enter tittle")
			}
			if accomplishment.FinishDate == nil {
				return errors.New("Please Enter Date")
			}
		}

	} else {
		return usersErrors.SpecificRequired
	}

	return nil
}

// //check education inputs to not be empty, startDate to not be after finishDate.
// //control input text length
func educationValidator(prof *profile.Education) error {

	strings.TrimSpace(prof.School)
	strings.TrimSpace(prof.FieldStudy)
	if prof == nil {
		return usersErrors.InvalidEnter
	}

	if prof.Description != nil && prof.Grade != nil {
		strings.TrimSpace(*prof.Description)
		strings.TrimSpace(*prof.Grade)
	}

	if prof.School == "" || prof.FieldStudy == "" || prof.StartDate.IsZero() {
		return usersErrors.SpecificRequired
	}
	if utf8.RuneCountInString(prof.FieldStudy) < 128 {
		err := lettersSymbols(prof.FieldStudy)
		if err != nil {
			return err
		}
	} else {
		return usersErrors.Max128
	}
	if !prof.FinishDate.IsZero() && prof.StartDate.After(prof.FinishDate) {
		return usersErrors.WrongDate
	}

	if prof.Degree != nil && *prof.Degree != "" {
		if utf8.RuneCountInString(*prof.Degree) < 64 {
			err := dashAndSpace(*prof.Degree)
			if err != nil {
				return err
			}
		} else {
			return usersErrors.Max64
		}
	}

	if prof.Grade != nil && *prof.Grade != "" {
		if utf8.RuneCountInString(*prof.Grade) < 32 {
			err := gradeValidator(prof.Grade)
			if err != nil {
				return err
			}
		} else {
			return usersErrors.Max32
		}
	}

	err := commonPatternValidator(prof.School)
	if err != nil {
		return err
	}

	if prof.Description != nil && utf8.RuneCountInString(*prof.Description) > 500 {
		return usersErrors.Max500
	}

	for _, url := range prof.Links {
		if url != nil && govalidator.IsURL((*url).URL) == false {
			return usersErrors.NotLink
		}
	}

	return nil
}

// //check experience inputs to not be empty, startDate to not be after finishDate.
// //control input text length
func experienceValidator(prof *profile.Experience) error {

	strings.TrimSpace(prof.Position)
	strings.TrimSpace(prof.Company)

	if prof == nil {
		return usersErrors.InvalidEnter
	}

	if prof.Position == "" || prof.Company == "" || prof.StartDate.IsZero() {
		return usersErrors.SpecificRequired
	}
	if prof.FinishDate != nil && !prof.FinishDate.IsZero() && prof.StartDate.After(*prof.FinishDate) {
		return usersErrors.WrongDate
	}

	if govalidator.StringMatches(prof.Position, commonPattern) == false {
		return usersErrors.CommonErrors
	}

	if govalidator.StringMatches(prof.Company, commonPattern) == false {
		return usersErrors.CommonErrors
	}

	if prof.Description != nil && *prof.Description != "" && utf8.RuneCountInString(*prof.Description) > 500 {

		strings.TrimSpace(*prof.Description)
		return usersErrors.Max500
	}

	for _, url := range prof.Links {
		if url != nil && govalidator.IsURL((*url).URL) == false {
			return usersErrors.NotLink
		}
	}

	return nil
}

func gradeValidator(str *string) error {
	if str != nil {
		if govalidator.StringMatches(*str, gradePattern) == false {
			return usersErrors.InvalidGrade
		}
	}
	return nil
}

func lettersNumbersSpecialValidator(str ...string) error {
	for _, s := range str {
		if s != "" {
			strings.TrimSpace(s)
			if govalidator.StringMatches(s, lettersNumbersSpecialPattern) == false {
				return usersErrors.LettersNumbersSpecial
			}
		}
	}
	return nil
}

func pointerLettersNumbersSpecialValidator(str ...*string) error {
	for _, s := range str {
		if s != nil {
			if *s != "" {
				strings.TrimSpace(*s)
				if govalidator.StringMatches(*s, lettersNumbersSpecialPattern) == false {
					return usersErrors.LettersNumbersSpecial
				}
			}
		}
	}
	return nil
}

func pointerLettersNumbersSpecialSpaceValidator(str ...*string) error {
	for _, s := range str {
		if s != nil {
			if *s != "" {
				strings.TrimSpace(*s)
				if govalidator.StringMatches(*s, lettersNumbersSpecialSpacePattern) == false {
					return usersErrors.LettersNumbersSpecialSpace
				}
			}
		}
	}
	return nil
}

func lettersNumbersSpecialSpaceValidator(str ...string) error {
	for _, s := range str {
		if s != "" {
			strings.TrimSpace(s)
			if govalidator.StringMatches(s, lettersNumbersSpecialSpacePattern) == false {
				return usersErrors.LettersNumbersSpecialSpace
			}
		}
	}
	return nil
}

func capitalToLower(str *string) {
	if str != nil {
		*str = strings.ToLower(*str)
	}
}

func urlValidator(name string) error {
	capitalToLower(&name)

	if name != "" && utf8.RuneCountInString(name) > 500 {
		return usersErrors.Max500
	}
	if govalidator.StringMatches(name, urlPattern) == false {
		return usersErrors.InvalidURL
	}
	return nil
}

func isNumber(str ...string) error {
	for i := range str {
		_, err := strconv.Atoi(str[i])
		if err != nil {
			return usersErrors.InvalidEnter
		}
	}
	return nil
}

func emptyValidator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if name[i] == "" {
			return usersErrors.SpecificRequired
		}
	}
	return nil
}

func pointerEmptyValidator(name ...*string) error {
	for i := range name {
		if name[i] != nil {
			strings.TrimSpace(*name[i])
			if *name[i] == "" {
				return usersErrors.SpecificRequired
			}
		}
	}
	return nil
}

func length64Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 64 {
			return usersErrors.Max64
		}
	}
	return nil
}

func length100Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 100 {
			return usersErrors.Max100
		}
	}
	return nil
}

func length128Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 128 {
			return usersErrors.Max128
		}
	}
	return nil
}

func pointerLenght128Validator(name ...*string) error {
	for i := range name {
		if name[i] != nil {
			strings.TrimSpace(*name[i])
			if utf8.RuneCountInString(*name[i]) > 128 {
				return usersErrors.Max128
			}
		}
	}
	return nil
}

func length500Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 500 {
			return usersErrors.Max500
		}
	}
	return nil
}

func length1200Validator(str ...string) error {
	for _, s := range str {
		strings.TrimSpace(s)
		if utf8.RuneCountInString(s) > 1200 {
			return usersErrors.Max1200
		}
	}

	return nil
}

// //control recommendation request text not to exceed 500 characters or be empty
// func recommendationRequestValidator(prof *profile.RecommendationRequest) error {
// 	if prof == nil {
// 		return usersErrors.InvalidEnter
// 	}
// 	if utf8.RuneCountInString(prof.Text) > 500 {
// 		return usersErrors.Max500
// 	}
// 	if prof.Text == "" {
// 		return usersErrors.SpecificRequired
// 	}

// 	return nil
// }

// // //control recommendation text not to exceed 500 characters or be empty
// func recommendationValidator(prof *profile.Recommendation) error {
// 	if prof == nil {
// 		return usersErrors.InvalidEnter
// 	}
// 	if utf8.RuneCountInString(prof.Text) > 500 {
// 		return usersErrors.Max500
// 	}
// 	if prof.Text == "" {
// 		return usersErrors.SpecificRequired
// 	}

// 	return nil
// }
