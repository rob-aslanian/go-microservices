package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
)

func TestLanguageValidator(t *testing.T) {

	// t1, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	// if err != nil {
	// 	t.Fatal("wrong date")
	// }
	// t2, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	// descr := "2222"

	tables := []struct {
		Language *profile.KnownLanguage
		Result   error
	}{
		{
			Language: nil,
			Result:   usersErrors.InvalidEnter,
		},
		{
			Language: &profile.KnownLanguage{
				Language: "",
			},
			Result: usersErrors.SpecificRequired,
		},
	}

	for _, ta := range tables {
		err := languageValidator(ta.Language)
		if err != ta.Result {
			t.Errorf("Error: got %v, expected: %v", err, nil)
		}
	}
}
