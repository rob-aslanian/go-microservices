package service

// import (
// 	"errors"
// 	"strconv"
// 	"strings"
// 	"time"
// 	"unicode/utf8"

// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/user_report"

// 	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
// 	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/profile"
// 	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/status"
// 	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"

// 	"github.com/asaskevich/govalidator"
// 	"github.com/ttacon/libphonenumber"
// )
import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/ttacon/libphonenumber"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/profile"
)

const (
	firstNamePattern                  = `^[a-zA-Z\- ]+$`
	namePattern                       = "^[a-zA-Z\\-'\\s]+(?:[a-zA-Z]+)$"
	vatPattern                        = "^[a-zA-Z0-9]+$"
	urlPattern                        = `^[a-zA-Z0-9\.-_]+$`
	alpahebtsNumbersSymbolsPattern    = `^[a-zA-Z0-9\W ]`
	unicodeNamePattern                = `^[\p{L}'-]+`
	numberPattern                     = `^[0-9]+`
	commonPattern                     = `^[a-zA-Z0-9\W ]{0,128}$`
	lettersNumbersSpecialSpacePattern = `^[a-zA-Z0-9\W ]`
	specailPattern                    = "[\\[\\]\\{\\}\\|\\<\\>\\(\\)\\'~!@#$%^&*?.,/\\-+]"
)

func alphabetNumberSymbolValidator(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		if govalidator.StringMatches(str[i], alpahebtsNumbersSymbolsPattern) == false {
			return companyErrors.AlphabetNumberSymbol
		}
	}
	return nil
}

func dashAndSpace(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		if govalidator.StringMatches(str[i], firstNamePattern) == false {
			return companyErrors.DashAndSpace
		}
	}
	return nil
}

func lettersNumbersSpecialSpaceValidator(str ...string) error {
	for _, s := range str {
		strings.TrimSpace(s)
		if govalidator.StringMatches(s, lettersNumbersSpecialSpacePattern) == false {
			return companyErrors.LettersNumbersSpecialSpace
		}
	}
	return nil
}

func fromTwoToSixtyFour(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		if utf8.RuneCountInString(str[i]) > 64 || utf8.RuneCountInString(str[i]) < 2 {
			return companyErrors.FromTwoToSixtyFour
		}
	}

	return nil
}
func eightToSixtyFour(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		if utf8.RuneCountInString(str[i]) > 64 || utf8.RuneCountInString(str[i]) < 8 {
			return companyErrors.EightToSixtyFour
		}
	}

	return nil
}

func commonPatternValidator(str ...string) error {
	for _, s := range str {
		strings.TrimSpace(s)
		if govalidator.StringMatches(s, commonPattern) == false {
			return companyErrors.CommonErrors
		}
	}
	return nil
}

func registerValidator(acc *account.Account) error {

	strings.TrimSpace(acc.Name)
	strings.TrimSpace(acc.URL)
	//checking if name or industry is empty
	err := emptyValidator(acc.Name, acc.URL)
	if err != nil {
		return err
	}

	//checking if name is over 100 characters
	err = length128Validator(acc.Name)
	if err != nil {
		return err
	}

	if len(acc.Emails) > 0 {
		err = emailValidator(acc.Emails[0].Email)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Please Enter Email")
	}
	//checking if it returns a number
	err = isNumber(acc.Industry.Main)
	if err != nil {
		return err
	}
	err = urlValidator(acc.URL)
	if err != nil {
		return err
	}

	//checking if zipcode, street and apartment is empty
	for i := range acc.Addresses {
		strings.TrimSpace(acc.Addresses[i].ZIPCode)
		strings.TrimSpace(acc.Addresses[i].Street)
		strings.TrimSpace(acc.Addresses[i].Apartment)
		err := emptyValidator(acc.Addresses[i].ZIPCode, acc.Addresses[i].Street /*, acc.Addresses[i].Apartment*/)
		if err != nil {
			return err
		}
	}
	//checking if zipcode, street and apartment is over 100 characters
	for i := range acc.Addresses {
		strings.TrimSpace(acc.Addresses[i].ZIPCode)
		strings.TrimSpace(acc.Addresses[i].Street)
		strings.TrimSpace(acc.Addresses[i].Apartment)
		err := length128Validator(acc.Addresses[i].ZIPCode)
		if err != nil {
			return err
		}
		// err = lettersNumbersSpecialSpaceValidator(acc.Addresses[i].Apartment)
		// if err != nil {
		// 	return err
		// }
		err = commonPatternValidator(acc.Addresses[i].Street)
		if err != nil {
			return err
		}
		//checking that zipcode is only numbers
		if govalidator.StringMatches(acc.Addresses[i].ZIPCode, numberPattern) == false {
			return companyErrors.InvalidEnter
		}
	}

	//checking if phone is valid
	for i := range acc.Phones {
		strings.TrimSpace(acc.Phones[i].Number)
		err := phoneValidator(acc.Phones[i].Number, acc.Phones[i].CountryCode.Code, acc.Phones[i].CountryCode.CountryID)
		if err != nil {
			return err
		}
	}

	if acc.VAT != nil {
		if govalidator.StringMatches(*acc.VAT, vatPattern) == false {
			return companyErrors.InvalidVAT
		}
		if utf8.RuneCountInString(*acc.VAT) > 128 {
			return companyErrors.Max128
		}
	}

	return nil
}

func capitalToLower(str *string) {
	if str != nil {
		*str = strings.ToLower(*str)
	}
}

func addressValidator(address *account.Address) error {
	strings.TrimSpace(address.Apartment)
	strings.TrimSpace(address.ZIPCode)
	strings.TrimSpace(address.Street)
	strings.TrimSpace(address.Location.Country.ID)
	strings.TrimSpace(address.Location.City.Name)

	err := length128Validator(address.Apartment, address.ZIPCode, address.Street)
	if err != nil {
		return err
	}

	err = emptyValidator(address.Apartment, address.ZIPCode, address.Street)
	if err != nil {
		return err
	}

	if govalidator.StringMatches(address.ZIPCode, numberPattern) == false {
		return companyErrors.InvalidEnter
	}

	return nil
}

func aboutUsValidator(aboutUs *profile.AboutUs) error {
	err := length500Validator(aboutUs.Mission)
	if err != nil {
		return err
	}
	err = length2000Validator(aboutUs.Description)
	if err != nil {
		return err
	}

	if aboutUs.FoundationDate.After(time.Now()) {
		return companyErrors.InvalidDate
	}

	if aboutUs.Parking == account.Parking("") || aboutUs.Type == account.Type("") ||
		aboutUs.Size == account.Size("") {
		return companyErrors.InvalidEnter
	}

	return nil
}

func reviewValidator(review *profile.Review) error {
	err := length128Validator(review.Headline)
	if err != nil {
		return err
	}
	err = length400Validator(review.Description)
	if err != nil {
		return err
	}

	return nil
}

func phoneValidator(phone string, countryCode string, countryID string) error {
	err := emptyValidator(phone)
	if err != nil {
		return err
	}

	if govalidator.StringMatches(phone, numberPattern) {

		phone = strings.Replace(phone, " ", "", -1)
		num, err := strconv.Atoi(phone)
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

		strings.TrimSpace(phone)
		//checking for these countries to match with their country phone numbers. Armenia doesn't work with ValidNumberForRegion so it is filtered
		//down with isValidNumber. It gives the right answer.
		if countryID == "GE" || countryID == "RU" || countryID == "US" || countryID == "UK" || countryID == "KZ" || countryID == "AZ" {
			if libphonenumber.IsValidNumberForRegion(&n, countryID) == false {
				return companyErrors.InValidPhone
			}
		} else {
			if libphonenumber.IsValidNumber(&n) == false {
				return companyErrors.InValidPhone
			}
		}
	} else {
		return companyErrors.InValidPhone
	}
	return nil
}

func foundationDateValidator(foundationDate time.Time) error {
	if foundationDate.After(time.Now()) {
		return companyErrors.InvalidDate
	}

	return nil
}

// check email
func emailValidator(email string) error {
	if govalidator.IsEmail(email) == false {
		return companyErrors.InValidEmail
	}
	if utf8.RuneCountInString(email) > 128 {
		return companyErrors.Max128
	}
	if email == "" {
		return companyErrors.InValidEmail
	}

	return nil
}

// func emailValidator(email *account.Email) error {
// 	strings.TrimSpace(email.Email)

// 	capitalToLower(&email.Email)

// 	if govalidator.IsEmail(email.Email) == false {
// 		return companyErrors.InvalidEmail
// 	}

// 	err := length128Validator(email.Email)
// 	if err != nil {
// 		return err
// 	}

// 	err = emptyValidator(email.Email)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func awardValidator(award *profile.Award) error {
	time, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

	if award != nil {
		if award.Date.Before(time) {
			return companyErrors.InvalidDate
		}
	}
	strings.TrimSpace(award.Title)
	strings.TrimSpace(award.Issuer)

	err := emptyValidator(award.Title)
	if err != nil {
		return err
	}

	if award.Issuer != "" {
		err = alphabetNumberSymbolValidator(award.Issuer)
		if err != nil {
			return err
		}
	}
	err = alphabetNumberSymbolValidator(award.Title)
	if err != nil {
		return err
	}

	err = length128Validator(award.Title, award.Issuer)
	if err != nil {
		return err
	}

	return nil
}

func milestoneValidator(milestone *profile.Milestone) error {
	time, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

	if milestone != nil {
		if milestone.Date.Before(time) {
			return companyErrors.InValidTime
		}
	}
	strings.TrimSpace(milestone.Title)
	strings.TrimSpace(milestone.Description)

	err := emptyValidator(milestone.Title)
	if err != nil {
		return err
	}
	err = alphabetNumberSymbolValidator(milestone.Title)
	if err != nil {
		return err
	}
	if milestone.Description != "" {
		err = alphabetNumberSymbolValidator(milestone.Description)
		if err != nil {
			return err
		}
	}
	err = length128Validator(milestone.Title)
	if err != nil {
		return err
	}
	err = length600Validator(milestone.Description)
	if err != nil {
		return err
	}

	return nil
}

func productValidator(product *profile.Product) error {
	strings.TrimSpace(product.Name)
	strings.TrimSpace(product.Website)
	err := emptyValidator(product.Name)
	if err != nil {
		return err
	}
	err = length128Validator(product.Name)
	if err != nil {
		return err
	}
	err = alphabetNumberSymbolValidator(product.Name)
	if err != nil {
		return err
	}
	err = length500Validator(product.Website)
	if err != nil {
		return err
	}
	err = urlValidator(product.Website)
	if err != nil {
		return err
	}

	return nil
}

func serviceValidator(service *profile.Service) error {
	strings.TrimSpace(service.Name)
	strings.TrimSpace(service.Website)
	err := emptyValidator(service.Name)
	if err != nil {
		return err
	}
	err = length128Validator(service.Name)
	if err != nil {
		return err
	}
	err = alphabetNumberSymbolValidator(service.Name)
	if err != nil {
		return err
	}
	err = length500Validator(service.Website)
	if err != nil {
		return err
	}
	err = urlValidator(service.Website)
	if err != nil {
		return err
	}

	return nil
}

func urlValidator(name string) error {
	if name != "" {
		if govalidator.StringMatches(name, urlPattern) == false {
			return companyErrors.InvalidURL
		}
	}

	err := length128Validator(name)
	if err != nil {
		return err
	}
	return nil
}

func isNumber(str ...string) error {
	for i := range str {
		_, err := strconv.Atoi(str[i])
		if err != nil {
			return companyErrors.InvalidEnter
		}
	}
	return nil
}

func emptyValidator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if name[i] == "" {
			return companyErrors.SpecificRequired
		}
	}
	return nil
}

// func length128Validator(name ...string) error {
// 	for i := range name {
// 		strings.TrimSpace(name[i])
// 		if utf8.RuneCountInString(name[i]) > 100 {
// 			return companyErrors.Max100
// 		}
// 	}
// 	return nil
// }

func length128Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 128 {
			return companyErrors.Max128
		}
	}
	return nil
}

func length400Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 400 {
			return companyErrors.Max400
		}
	}
	return nil
}

func length500Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 500 {
			return companyErrors.Max500
		}
	}
	return nil
}

func length600Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 600 {
			return companyErrors.Max600
		}
	}
	return nil
}

func length2000Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 2000 {
			return companyErrors.Max2000
		}
	}
	return nil
}
