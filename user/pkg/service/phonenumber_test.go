package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

func TestPhoneNumberValidator(t *testing.T) {

	tables := []struct {
		Number      string
		CountryCode string
		countryID   string
		Result      error
	}{
		{
			Number:      "598661708",
			CountryCode: "995",
			countryID:   "GE",
			Result:      nil,
		},
		{
			Number:      "598661708",
			CountryCode: "995",
			countryID:   "US",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "5986617",
			CountryCode: "995",
			countryID:   "GE",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "$$##$",
			CountryCode: "995",
			countryID:   "US",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "dfgdfgf",
			CountryCode: "995",
			countryID:   "US",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "598661708",
			CountryCode: "995",
			countryID:   "US",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "2025550184",
			CountryCode: "1",
			countryID:   "US",
			Result:      nil,
		},
		{
			Number:      "20255501",
			CountryCode: "1",
			countryID:   "US",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "9226625196",
			CountryCode: "7",
			countryID:   "RU",
			Result:      nil,
		},
		{
			Number:      "7212438555",
			CountryCode: "7",
			countryID:   "KZ",
			Result:      nil,
		},
		{
			Number:      "9226625196",
			CountryCode: "2",
			countryID:   "RU",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "555013995",
			CountryCode: "994",
			countryID:   "AZ",
			Result:      nil,
		},
		{
			Number:      "555013995",
			CountryCode: "995",
			countryID:   "AZ",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "55501399533",
			CountryCode: "995",
			countryID:   "AZ",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "10531710",
			CountryCode: "374",
			countryID:   "AR",
			Result:      nil,
		},
		{
			Number:      "1053171022",
			CountryCode: "374",
			countryID:   "AR",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "10531710",
			CountryCode: "3724",
			countryID:   "AR",
			Result:      usersErrors.InValidPhone,
		},
		{
			Number:      "0686579014",
			CountryCode: "33",
			countryID:   "FR",
			Result:      nil,
		},
	}

	for _, ta := range tables {
		err := phoneValidator(ta.CountryCode, ta.Number, ta.countryID)
		if err != ta.Result {
			t.Errorf("Error: TestPhoneNumberValidator(%v): got %v, expected: %v", ta.Number, err, nil)
		}
	}
}
