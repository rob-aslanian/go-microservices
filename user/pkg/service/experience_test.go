package service

import (
	"testing"
	"time"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
)

func TestExperienceValidator(t *testing.T) {

	t1, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	if err != nil {
		t.Fatal("wrong date")
	}
	t2, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2007")
	if err != nil {
		t.Fatal("wrong date")
	}
	str := "fdfdf"

	tables := []struct {
		Exp        *profile.Experience
		StartDate  time.Time
		FinishDate time.Time
		Result     error
	}{
		{
			Exp:    &profile.Experience{},
			Result: usersErrors.SpecificRequired,
		},
		{
			Exp: &profile.Experience{
				Position:    "1212",
				Company:     "%@@!",
				StartDate:   t1,
				FinishDate:  &t2,
				Description: &str,
			},
			Result: nil,
		},
		{
			Exp: &profile.Experience{
				Position:    "21312",
				Company:     "1232131",
				StartDate:   t1,
				Description: &str,
			},
			Result: nil,
		},
	}

	for _, ta := range tables {
		err := experienceValidator(ta.Exp)
		if err != ta.Result {
			t.Errorf("Error: got %v, expecrted: %v", err, nil)
		}
	}
}
