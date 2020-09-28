package service

import (
	"strconv"
	"strings"
	"unicode/utf8"

	jobShared "gitlab.lan/Rightnao-site/microservices/jobs/internal/job-functions"

	"github.com/asaskevich/govalidator"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	jobsErrors "gitlab.lan/Rightnao-site/microservices/jobs/internal/job-errors"
)

const (
	namePattern        = "^[a-zA-Z\\-'\\s]+(?:[a-zA-Z]+)$"
	vatPattern         = "^[a-zA-Z0-9]+$"
	urlPattern         = "^[a-z0-9]+$"
	unicodeNamePattern = `^[\p{L}'-]+`
	languagePattern    = `^[a-z-]{2,5}$`
	numberPattern      = `^[0-9]+`
	specailPattern     = "[\\[\\]\\{\\}\\|\\<\\>\\(\\)\\'~!@#$%^&*?.,/\\-+]"
)

func careerInterestsValidator(data *candidate.CareerInterests) error {
	if data == nil {
		return jobsErrors.PointerIsNill
	}
	for i := range data.Jobs {
		err := length128Validator(data.Jobs[i])
		if err != nil {
			return err
		}
		err = emptyValidator(data.Jobs[i])
		if err != nil {
			return err
		}
	}
	if utf8.RuneCountInString(data.SalaryCurrency) != 3 {
		return jobsErrors.InvalidCurrency
	}

	if data.SalaryMin < 0 || data.SalaryMax > 999999999 || data.SalaryMax != 0 && (data.SalaryMin > data.SalaryMax) {
		return jobsErrors.InvalidSalary
	}

	err := isNumber(data.Industry)
	if err != nil {
		return err
	}

	return nil
}

func postJobValidator(post *job.Posting) error {

	err := emptyValidator(post.JobDetails.Title)
	if err != nil {
		return err
	}

	err = length500Validator(post.JobDetails.Title)
	if err != nil {
		return err
	}

	err = isNumber(post.JobDetails.Country, post.JobDetails.City)
	if err != nil {
		return err
	}

	for _, t := range post.JobDetails.EmploymentTypes {
		if t == candidate.JobType("") || t == candidate.JobTypeUnknown {
			return jobsErrors.InvalidJobType
		}
	}

	for _, t := range post.JobDetails.JobFunctions {
		err := length128Validator(jobShared.String(t))
		if err != nil {
			return err
		}
	}

	jobFunctionsLength := len(post.JobDetails.JobFunctions)
	if jobFunctionsLength > 5 {
		return jobsErrors.InvalidJobFunctionNumber
	}

	if post.JobDetails.Descriptions == nil {
		return jobsErrors.PointerIsNill
	}

	for _, t := range post.JobDetails.Descriptions {
		if govalidator.StringMatches(t.Language, languagePattern) == false {
			return jobsErrors.InvalidLanguage
		}
		err = emptyValidator(t.Description)
		if err != nil {
			return err
		}
		err = length2000Validator(t.Description, t.WhyUs)
		if err != nil {
			return err
		}
	}
	return nil
}

func emptyValidator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if name[i] == "" {
			return jobsErrors.SpecificRequired
		}
	}
	return nil
}

func length128Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 128 {
			return jobsErrors.Max128
		}
	}
	return nil
}

func length500Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 500 {
			return jobsErrors.Max500
		}
	}
	return nil
}

func length2000Validator(name ...string) error {
	for i := range name {
		strings.TrimSpace(name[i])
		if utf8.RuneCountInString(name[i]) > 2000 {
			return jobsErrors.Max2000
		}
	}
	return nil
}

func isNumber(str ...string) error {
	for i := range str {
		_, err := strconv.Atoi(str[i])
		if err != nil {
			return jobsErrors.InvalidEnter
		}
	}
	return nil
}
