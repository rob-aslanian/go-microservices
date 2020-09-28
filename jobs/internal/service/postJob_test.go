package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job-errors"
)

func TestPostJobValidator(t *testing.T) {

	tables := []struct {
		Post   *job.Posting
		Result error
	}{
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						"",
						// candidate.JobTypeUnknown,
					},
					JobFunctions: []string{
						"write code",
					},
					Descriptions: []*job.Description{
						{
							Language:    "en-a",
							Description: "222",
						},
					},
				},
			},
			Result: jobsErrors.InvalidJobType,
		},
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						candidate.JobTypeConsultancy,
					},
					JobFunctions: []string{
						"dfdf",
					},
					Descriptions: []*job.Description{
						{
							Language:    "en-a",
							Description: "222",
						},
					},
				},
			},
			Result: nil,
		},
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						candidate.JobTypeConsultancy,
					},
					JobFunctions: []string{
						"dfdfdfdfdfdfdfdfddddddddddddddddddddddddddddddddddddddfdfdfdfdfdfdfdfddddddddddddddddddddddddddddddddddddddfdfdfdfdfdfdfdfddddddddddddddddddddddddddddddddddddddfdfdfdfdfdfdfdfddddddddddddddddddddddddddddddddddddddfdfdfdfdfdfdfdfddddddddddddddddddddddddddddddddddddd",
					},
					Descriptions: []*job.Description{
						{
							Language:    "en-a",
							Description: "222",
						},
					},
				},
			},
			Result: jobsErrors.Max128,
		},
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						candidate.JobTypeConsultancy,
					},
					JobFunctions: []string{"", "", "", ""},
					Descriptions: []*job.Description{
						{
							Language:    "en-a",
							Description: "222",
						},
					},
				},
			},
			Result: jobsErrors.InvalidJobFunctionNumber,
		},
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						candidate.JobTypeConsultancy,
					},
					JobFunctions: []string{},
					Descriptions: []*job.Description{
						{
							Language:    "en-a",
							Description: "222",
						},
					},
				},
			},
			Result: jobsErrors.InvalidJobFunctionNumber,
		},
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						candidate.JobTypeConsultancy,
					},
					JobFunctions: []string{""},
					Descriptions: nil,
				},
			},
			Result: jobsErrors.PointerIsNill,
		},
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						candidate.JobTypeConsultancy,
					},
					JobFunctions: []string{""},
					Descriptions: []*job.Description{
						{
							Language:    "en-a",
							Description: "",
						},
					},
				},
			},
			Result: jobsErrors.SpecificRequired,
		},
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						candidate.JobTypeConsultancy,
					},
					JobFunctions: []string{""},
					Descriptions: []*job.Description{
						{
							Language:    "en-a",
							Description: "ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
						},
					},
				},
			},
			Result: jobsErrors.Max2000,
		},
		{
			Post: &job.Posting{
				JobDetails: job.Details{
					Title:   "dfdfd",
					City:    "2",
					Country: "2",
					EmploymentTypes: []candidate.JobType{
						candidate.JobTypeConsultancy,
					},
					JobFunctions: []string{""},
					Descriptions: []*job.Description{
						{
							Language:    "en-a",
							Description: "good",
						},
					},
					RequiredEducations: []string{
						"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddsssssssss",
					},
				},
			},
			Result: jobsErrors.Max128,
		},
	}

	for _, ta := range tables {
		err := postJobValidator(ta.Post)
		if err != ta.Result {
			t.Errorf("Error: registerValidator(%+v) got %v, expected: %v", ta.Post, err, nil)
		}
	}
}
