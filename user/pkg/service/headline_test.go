package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

func TestHeadlineValidator(t *testing.T) {

	tables := []struct {
		Headline string
		Result   error
	}{
		{
			Headline: "",
			Result:   nil,
		},
		{
			Headline: "99559861709",
			Result:   nil,
		},
		{
			Headline: "99559861709dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd99559861709ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
			Result:   usersErrors.Max128,
		},
	}

	for _, ta := range tables {
		err := headlineValidator(ta.Headline)
		if err != ta.Result {
			t.Errorf("Error: TestHeadlineValidator(%v): got %v, expected: %v", ta.Headline, err, nil)
		}
	}
}
