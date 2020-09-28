package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

func TestNativeLanguage(t *testing.T) {

	Invalid := "$$"
	LongName := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	Valid := "ვახხ"

	lang := "en"
	inLang := "ESS"

	tables := []struct {
		Name     *string
		Language *string
		Result   error
	}{
		{
			Name:     &Invalid,
			Language: &lang,
			Result:   usersErrors.InValidName,
		},
		{
			Language: &inLang,
			Name:     &Valid,
			Result:   usersErrors.InValidLang,
		},
		{
			Name:     &Invalid,
			Language: &inLang,
			Result:   usersErrors.InValidName,
		},
		{
			Name:     &LongName,
			Language: &lang,
			Result:   usersErrors.Max128,
		},
	}

	for _, ta := range tables {
		err := changeNameOnNativeLanguageValidator(ta.Name, ta.Language)
		if err != ta.Result {
			t.Errorf("Error:   got %v, expected: %v", err, ta.Result)
		}
	}
}
