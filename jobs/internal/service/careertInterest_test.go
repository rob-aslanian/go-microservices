package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job-errors"
)

func TestCareerInterestValidator(t *testing.T) {

	tables := []struct {
		Data   *candidate.CareerInterests
		Result error
	}{
		{
			Data: &candidate.CareerInterests{
				Jobs: []string{
					"dfdfdf",
				},
				SalaryCurrency: "saa",
				SalaryMin:      12,
				SalaryMax:      999,
				Industry:       "12",
			},
			Result: nil,
		},
		{
			Data: &candidate.CareerInterests{
				Jobs: []string{
					"dfdfdf",
				},
				SalaryCurrency: "saa2",
				SalaryMin:      12,
				SalaryMax:      999,
				Industry:       "12",
			},
			Result: jobsErrors.InvalidCurrency,
		},
		{
			Data: &candidate.CareerInterests{
				Jobs: []string{
					"dfdfdf",
				},
				SalaryCurrency: "saa",
				SalaryMin:      -5,
				SalaryMax:      999,
				Industry:       "12",
			},
			Result: jobsErrors.InvalidSalary,
		},
		{
			Data: &candidate.CareerInterests{
				Jobs: []string{
					"dfdfdf",
				},
				SalaryCurrency: "saa",
				SalaryMin:      10000,
				SalaryMax:      999,
				Industry:       "12",
			},
			Result: jobsErrors.InvalidSalary,
		},
		{
			Data: &candidate.CareerInterests{
				Jobs: []string{
					"dfdfdf",
				},
				SalaryCurrency: "saa",
				SalaryMin:      12,
				SalaryMax:      999,
				Industry:       "sdsd",
			},
			Result: jobsErrors.InvalidEnter,
		},
		{
			Data: &candidate.CareerInterests{
				Jobs: []string{
					"saasdssssssssssssssssssssssssssssssssssssssssssssssssssssssssssaasdssssssssssssssssssssssssssssssssssssssssssssssssssssssssssaasdssssssssssssssssssssssssssssssssssssssssssssssssssssssssssaasdsssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
				},
				SalaryCurrency: "ss",
				SalaryMin:      12,
				SalaryMax:      999,
				Industry:       "12",
			},
			Result: jobsErrors.Max128,
		},
	}

	for _, ta := range tables {
		err := careerInterestsValidator(ta.Data)
		if err != ta.Result {
			t.Errorf("Error: registerValidator(%+v) got %v, expected: %v", ta.Data, err, nil)
		}
	}
}
