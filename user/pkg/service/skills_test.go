package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
)

func TestSkillsValidator(t *testing.T) {

	// t1, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	// if err != nil {
	// 	t.Fatal("wrong date")
	// }
	// t2, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	// descr := "2222"

	tables := []struct {
		Skill  []*profile.Skill
		Result error
	}{
		{
			Skill:  []*profile.Skill{},
			Result: usersErrors.SpecificRequired,
		},
		{
			Skill: []*profile.Skill{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
			},
			Result: usersErrors.MaxArray,
		},
		{
			Skill: []*profile.Skill{
				{
					Skill: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				},
			},
			Result: usersErrors.Max64,
		},
	}

	for _, ta := range tables {
		err := skillsValidator(ta.Skill)
		if err != ta.Result {
			t.Errorf("Error: got %v, expected: %v", err, ta.Result)
		}
	}
}
