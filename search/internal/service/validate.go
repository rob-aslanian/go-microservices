package service

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"

	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"

	"gitlab.lan/Rightnao-site/microservices/search/internal/search-errors"
)

const (
	namePattern        = "^[a-zA-Z\\-'\\s]+(?:[a-zA-Z]+)$"
	vatPattern         = "^[a-zA-Z0-9]+$"
	urlPattern         = `^[a-z0-9\.]+$`
	unicodeNamePattern = `^[\p{L}'-]+`
	numberPattern      = `^[0-9]+`
	specailPattern     = "[\\[\\]\\{\\}\\|\\<\\>\\(\\)\\'~!@#$%^&*?.,/\\-+]"
)

// userSearchValidator trims the entered fields and checks for them not to exceed 100 characters
// It also checks that city, country and connections IDs to be numbers which means user must pass correct fields.
// Checks for birthday not to be older then 1960 and newer then current time
func userSearchValidator(data *requests.UserSearch) error {
	err := length100ValidatorForArrays(data.Firstname, data.Language, data.Nickname, data.PastCompany, data.CurrentCompany, data.FieldOfStudy, data.Industry, data.Interest, data.Position, data.Skill, data.Degree, data.School)
	if err != nil {
		return err
	}

	length100ValidatorForArrays(data.Keyword)
	if err != nil {
		return err
	}
	err = isCountry(data.CountryID)
	if err != nil {
		return err
	}

	err = isCity(data.CityID)
	if err != nil {
		return err
	}

	err = isOnlyLetters(data.Firstname, data.Language, data.Lastname)
	if err != nil {
		return err
	}

	err = dateValidator(data)
	if err != nil {
		return err
	}

	return nil
}

// jobSearch trims the entered fields and checks for them not to exceed 100 characters
// It also checks that city, country and connections IDs to be numbers which means user must pass correct fields.
// Checks for birthday not to be older then 1960 and newer then current time
func jobSearchValidator(data *requests.JobSearch) error {

	err := length100ValidatorForArrays(data.Degree, data.Industry, data.JobType, data.Language, data.Skill)
	if err != nil {
		return err
	}
	err = isCountry(data.Country)
	if err != nil {
		return err
	}

	err = isCity(data.CityID)
	if err != nil {
		return err
	}
	if data.MaxSalary > 999999999 || data.MaxSalary != 0 && (data.MinSalary > data.MaxSalary) {
		return searchErrors.InvalidSalary
	}

	// err = dateValidatorForJob(data)
	// if err != nil {
	// 	return err
	// }

	err = length100ValidatorForArrays(data.Keyword)
	if err != nil {
		return err
	}

	return nil
}

// candidateSearchValidator trims the entered fields and checks for them not to exceed 100 characters
// It also checks that city, country and connections IDs to be numbers which means user must pass correct fields.
// Checks for birthday not to be older then 1960 and newer then current time
func candidateSearchValidator(data *requests.CandidateSearch) error {
	err := length100ValidatorForArrays(data.JobType, data.Language, data.PastCompany, data.CurrentCompany, data.CityID, data.Country, data.Degree, data.Industry, data.School, data.Skill, data.SubIndustry)
	if err != nil {
		return err
	}

	if data.MaxSalary > 999999999 || data.MaxSalary != 0 && (data.MinSalary > data.MaxSalary) {
		return searchErrors.InvalidSalary
	}

	err = length100ValidatorForArrays(data.Keyword)
	if err != nil {
		return err
	}

	err = length100Validator(data.Period)
	if err != nil {
		return err
	}
	err = isCountry(data.Country)
	if err != nil {
		return err
	}

	err = isCity(data.CityID)
	if err != nil {
		return err
	}

	return nil
}

// companySearchValidator trims the entered fields and checks for them not to exceed 100 characters
// It also checks that city, country and connections IDs to be numbers which means user must pass correct fields.
// Checks for birthday not to be older then 1960 and newer then current time
func companySearchValidator(data *requests.CompanySearch) error {
	err := length100ValidatorForArrays(data.SubIndustry, data.BusinessHours, data.CityID, data.CityID, data.FoundersID, data.Name, data.Rating)
	if err != nil {
		return err
	}
	err = length100ValidatorForArrays(data.Keyword)
	if err != nil {
		return err
	}
	err = isCountry(data.Country)
	if err != nil {
		return err
	}

	err = isCity(data.CityID)
	if err != nil {
		return err
	}

	return nil
}

// saveUserSearchFilterValidator uses the previous userSearchValidator as well as checks for name and UserID not to be mpty
func saveUserSearchFilterValidator(data *requests.UserSearchFilter) error {
	err := emptyValidator(data.Name)
	if err != nil {
		return err
	}
	err = length100Validator(data.Name)
	if err != nil {
		return err
	}

	err = userSearchValidator(&data.UserSearch)
	if err != nil {
		return err
	}

	return nil
}

// saveJobSearchFilterValidator uses the previous jobSearchValidator as well as checks for name and UserID not to be mpty
func saveJobSearchFilterValidator(data *requests.JobSearchFilter) error {
	err := emptyValidator(data.Name)
	if err != nil {
		return err
	}
	err = length100Validator(data.Name)
	if err != nil {
		return err
	}

	err = jobSearchValidator(&data.JobSearch)
	if err != nil {
		return err
	}

	return nil
}

// saveCompanySearchFilterValidator uses the previous companySearchValidator as well as checks for name and UserID not to be mpty
func saveCompanySearchFilterValidator(data *requests.CompanySearchFilter) error {
	err := emptyValidator(data.Name)
	if err != nil {
		return err
	}
	err = length100Validator(data.Name)
	if err != nil {
		return err
	}

	err = companySearchValidator(&data.CompanySearch)
	if err != nil {
		return err
	}

	return nil
}

// saveCandidateSearchFilterValidator uses the previous companySearchValidator as well as checks for name and UserID not to be mpty
func saveCandidateSearchFilterValidator(data *requests.CandidateSearchFilter) error {
	err := emptyValidator(data.Name)
	if err != nil {
		return err
	}
	err = length100Validator(data.Name)
	if err != nil {
		return err
	}

	err = candidateSearchValidator(&data.CandidateSearch)
	if err != nil {
		return err
	}

	return nil
}

// saveUserSearchFilterValidatorForCompany uses the previous userSearchValidator as well as checks for name and UserID not to be mpty
func saveUserSearchFilterValidatorForCompany(data *requests.UserSearchFilter) error {
	err := emptyValidator(data.Name)
	if err != nil {
		return err
	}
	err = length100Validator(data.Name)
	if err != nil {
		return err
	}

	err = userSearchValidator(&data.UserSearch)
	if err != nil {
		return err
	}

	return nil
}

// saveJobSearchFilterValidatorForCompany uses the previous jobSearchValidator as well as checks for name and UserID not to be mpty
func saveJobSearchFilterValidatorForCompany(data *requests.JobSearchFilter) error {
	err := emptyValidator(data.Name)
	if err != nil {
		return err
	}
	err = length100Validator(data.Name)
	if err != nil {
		return err
	}

	err = jobSearchValidator(&data.JobSearch)
	if err != nil {
		return err
	}

	return nil
}

// saveCompanySearchFilterValidatorForCompany uses the previous companySearchValidator as well as checks for name and UserID not to be mpty
func saveCompanySearchFilterValidatorForCompany(data *requests.CompanySearchFilter) error {
	err := emptyValidator(data.Name)
	if err != nil {
		return err
	}
	err = length100Validator(data.Name)
	if err != nil {
		return err
	}

	err = companySearchValidator(&data.CompanySearch)
	if err != nil {
		return err
	}

	return nil
}

// saveCandidateSearchFilterValidatorForCompany uses the previous companySearchValidator as well as checks for name and UserID not to be mpty
func saveCandidateSearchFilterValidatorForCompany(data *requests.CandidateSearchFilter) error {
	err := emptyValidator(data.Name, data.GetCompanyID())
	if err != nil {
		return err
	}
	err = length100Validator(data.Name)
	if err != nil {
		return err
	}

	err = candidateSearchValidator(&data.CandidateSearch)
	if err != nil {
		return err
	}

	return nil
}

// trimArraysOfStrings trims each field inside an array of strings,  I pass an array of arrays, so that it could check all of them at once
// func trimArraysOfStrings(arr ...[]string) {
// 	for _, each := range arr {
// 		for _, item := range each {
// 			strings.TrimSpace(item)
// 		}
// 	}
// }

// trimStrings ...
func trimStrings(str ...string) {
	for _, each := range str {
		strings.TrimSpace(each)
	}
}

// dateFormatter formats Date format into time.Time
func dateFormatter(data *requests.Date) (time.Time, error) {
	if data == nil {
		return time.Time{}, searchErrors.InvalidEnter
	}
	date := time.Date(int(data.Year), time.Month(data.Month), int(data.Day), 0, 0, 0, 0, time.UTC)

	return date, nil

}

// dateValidator checks for date not to be older then 1960 and newer then current time
func dateValidator(data *requests.UserSearch) error {

	if data == nil {
		return searchErrors.InvalidEnter
	}

	if data.MaxAge > 150 || data.MinAge > 150 || (data.MaxAge != 0 && (data.MinAge > data.MaxAge)) {
		return searchErrors.InvalidAge
	}

	return nil
}

// dateValidatorForJob checks for date not to be older then 1960 and newer then current time
// func dateValidatorForJob(data *requests.JobSearch) error {
// 	oldest, _ := time.Parse(time.RFC1123, "Fri, 01 Jan 1960 00:00:00 MST")

// 	if data == nil {
// 		return searchErrors.InvalidEnter
// 	}

// 	for _, each := range data.DatePosted {
// 		date, err := dateFormatter(each)
// 		if err != nil {
// 			return err
// 		}
// 		if date.Before(oldest) || date.After(time.Now()) {
// 			return searchErrors.InvalidDate
// 		}
// 	}

// 	return nil
// }

func isNumber(str ...string) error {
	for i := range str {
		strings.TrimSpace(str[i])
		_, err := strconv.Atoi(str[i])
		if err != nil {
			return searchErrors.InvalidEnter
		}
	}
	return nil
}

func isNumberForArrays(arr ...[]string) error {
	for _, each := range arr {
		for _, item := range each {
			strings.TrimSpace(item)
			_, err := strconv.Atoi(item)
			if err != nil {
				return searchErrors.InvalidEnter
			}
		}
	}
	return nil
}

func emptyValidator(name ...string) error {
	for i := range name {
		if name[i] == "" {
			return searchErrors.SpecificRequired
		}
	}
	return nil
}

func length100ValidatorForArrays(arr ...[]string) error {
	for _, each := range arr {
		for _, item := range each {
			strings.TrimSpace(item)
			if utf8.RuneCountInString(item) > 100 {
				return searchErrors.Max100
			}
		}
	}
	return nil
}

func length100Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 100 {
			return searchErrors.Max100
		}
	}
	return nil
}

func isOnlyLetters(arr ...[]string) error {
	for _, each := range arr {
		for _, item := range each {
			if govalidator.StringMatches(item, namePattern) == false {
				return searchErrors.InvalidName
			}
		}
	}

	return nil
}

func isCountry(counties []string) error {
	for _, c := range counties {
		if utf8.RuneCountInString(c) > 2 || utf8.RuneCountInString(c) == 1 {
			return searchErrors.InvalidCountry
		}
	}

	return nil
}

func isCity(cities []string) error {
	for _, c := range cities {
		_, err := strconv.Atoi(c)
		if err != nil {
			return searchErrors.InvalidCity
		}
	}

	return nil
}
