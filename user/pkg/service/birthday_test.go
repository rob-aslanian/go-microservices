package service

import (
	"testing"
	"time"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

func TestBirthdayValiator(t *testing.T) {

	t1, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	if err != nil {
		t.Fatal("wrong date")
	}

	t2, err := time.Parse(time.UnixDate, "Mon Jul 2 15:04:05 MST 1945")
	if err != nil {
		t.Fatal("wrong date")
	}
	if err != nil {
		t.Fatal("wrong date")
	}

	tables := []struct {
		Birthday *time.Time
		Result   error
	}{
		{
			Birthday: nil,
			Result:   nil,
		},
		{
			Birthday: &t1,
			Result:   nil,
		},
		{
			Birthday: &t2,
			Result:   usersErrors.InValidTime,
		},
	}

	for _, ta := range tables {
		err := birthdayValidator(ta.Birthday)
		if err != ta.Result {
			t.Errorf("Error: TestBirthdayValiator(%v): got %v, expected: %v", ta.Birthday, err, nil)
		}
	}
}
