package service

import (
	"testing"
	"time"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
)

func TestEducationValidator(t *testing.T) {

	t1, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	if err != nil {
		t.Fatal("wrong date")
	}
	t2, err := time.Parse(time.UnixDate, "Mon Jan 2 15:09:05 MST 2006")
	if err != nil {
		t.Fatal("wrong date")
	}

	tables := []struct {
		Edu              *profile.Education
		StartDate        time.Time
		FinishDate       time.Time
		IsCurrentlyStudy bool
		Result           error
	}{
		{
			Edu:    &profile.Education{},
			Result: usersErrors.SpecificRequired,
		},
		{
			Edu: &profile.Education{
				School:           "1212",
				FieldStudy:       "s",
				IsCurrentlyStudy: true,
				StartDate:        t1,
				FinishDate:       t2,
			},
			Result: nil,
		},
		{
			Edu: &profile.Education{
				School:           "21312",
				FieldStudy:       "sdfsdfsfd",
				StartDate:        t1,
				IsCurrentlyStudy: true,
				FinishDate:       t2,
			},
			Result: nil,
		},
	}

	for _, ta := range tables {
		err := educationValidator(ta.Edu)
		if err != ta.Result {
			t.Errorf("Error: educationValidator(%+v) got %v, expecrted: %v", ta.Edu, err, ta.Result)
		}
	}
}
