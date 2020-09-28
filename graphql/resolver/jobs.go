package resolver

import (
	"context"
	"log"
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
)

var jobsEmpty = &jobsRPC.Empty{}

func (_ *Resolver) GetJobProfile(ctx context.Context) (*CandidateProfileResolver, error) {
	res, err := jobs.GetProfile(ctx, jobsEmpty)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	return &CandidateProfileResolver{R: jobs_candidateProfileToGql(ctx, res)}, nil
}

func (_ *Resolver) GetRecommendedJobs(ctx context.Context, input GetRecommendedJobsRequest) ([]JobPostingResolver, error) {
	var first, after int32

	if input.Pagination.After != nil {
		afterInt, _ := strconv.Atoi(*input.Pagination.After)
		after = int32(afterInt)
	}

	first = NullToInt32(input.Pagination.First)

	res, err := jobs.GetRecommendedJobs(ctx, &jobsRPC.Pagination{
		First: first,
		After: after,
	})

	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	companiesIDs := make([]string, 0)

	for _, l := range res.GetList() {
		if l.GetCompanyInfo() != nil {
			companiesIDs = append(companiesIDs, l.GetCompanyInfo().GetCompanyId())
		}
	}

	// companyProfiles, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
	// 	Ids: companiesIDs,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	list := make([]JobPostingResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_jobPostingToResolver(ctx, o)
		// if o.GetCompanyInfo() != nil {
		// 	if cID := o.GetCompanyInfo().GetCompanyId(); cID != "" {
		// 		company := CompanyProfile{}
		//
		// 		for _, j := range companyProfiles.GetProfiles() {
		// 			if j != nil && j.Id == res.List[i].GetId() {
		// 				company = toCompanyProfile(ctx, *j)
		// 				break
		// 			}
		// 		}
		//
		// 		list[i].R.Company = company
		// 	}
		// }
	}

	return list, nil
}

func (_ *Resolver) GetJob(ctx context.Context, input GetJobRequest) (*JobPostingResolver, error) {
	res, err := jobs.GetJob(ctx, &jobsRPC.ID{Id: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	log.Printf("enums: %v\n", res.JobDetails.AdditionalInfo)
	return jobs_jobPostingToResolver(ctx, res), nil
}

func (_ *Resolver) GetJobApplicants(ctx context.Context, input GetJobApplicantsRequest) ([]JobApplicantResolver, error) {
	sort := func(s *string) jobsRPC.GetJobApplicantsRequest_JobApplicantsSort {
		if s != nil {
			switch *s {
			case "first_name":
				return jobsRPC.GetJobApplicantsRequest_Firstname
			case "last_name":
				return jobsRPC.GetJobApplicantsRequest_Lastname
			case "posted_date":
				return jobsRPC.GetJobApplicantsRequest_PostedDate
			case "expeirence_years":
				return jobsRPC.GetJobApplicantsRequest_ExpeirenceYears
			}
		}
		return jobsRPC.GetJobApplicantsRequest_Lastname
	}(input.Sort)

	var first, after int32

	if input.Pagination.After != nil {
		afterInt, _ := strconv.Atoi(*input.Pagination.After)
		after = int32(afterInt)
	}
	first = NullToInt32(input.Pagination.First)

	res, err := jobs.GetJobApplicants(ctx, &jobsRPC.GetJobApplicantsRequest{
		CompanyID: input.CompanyId,
		JobID:     input.JobId,
		Sort:      sort,
		First:     first,
		After:     after,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]JobApplicantResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_jobApplicantToResolver(ctx, o)
	}

	return list, nil
}

func (_ *Resolver) GetPostedJobs(ctx context.Context, input GetPostedJobsRequest) ([]JobPostingResolver, error) {
	res, err := jobs.GetPostedJobs(ctx, &jobsRPC.ID{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]JobPostingResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_jobPostingToResolver(ctx, o)
	}
	return list, nil
}

func (_ *Resolver) GetJobForCompany(ctx context.Context, input GetJobForCompanyRequest) (*JobPostingResolver, error) {
	res, err := jobs.GetJobForCompany(ctx, &jobsRPC.CompanyIdWithJobId{CompanyId: input.CompanyId, JobId: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	return jobs_jobPostingToResolver(ctx, res), nil
}

func (_ *Resolver) GetSavedJobs(ctx context.Context, input GetSavedJobsRequest) ([]JobPostingResolver, error) {
	pagination := jobsRPC.Pagination{}

	if input.Pagination.After != nil {
		after, _ := strconv.Atoi(*input.Pagination.After)
		pagination.After = int32(after)
	}
	pagination.First = NullToInt32(input.Pagination.First)

	res, err := jobs.GetSavedJobs(ctx, &pagination)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	list := make([]JobPostingResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_jobPostingToResolver(ctx, o)
	}
	return list, nil
}

func (_ *Resolver) GetSkippedJobs(ctx context.Context, input GetSkippedJobsRequest) ([]JobPostingResolver, error) {
	pagination := jobsRPC.Pagination{}

	if input.Pagination.After != nil {
		after, _ := strconv.Atoi(*input.Pagination.After)
		pagination.After = int32(after)
	}
	pagination.First = NullToInt32(input.Pagination.First)

	res, err := jobs.GetSkippedJobs(ctx, &pagination)

	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]JobPostingResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_jobPostingToResolver(ctx, o)
	}
	return list, nil
}

func (_ *Resolver) GetAppliedJobs(ctx context.Context, input GetAppliedJobsRequest) ([]JobPostingResolver, error) {
	pagination := jobsRPC.Pagination{}

	if input.Pagination.After != nil {
		after, _ := strconv.Atoi(*input.Pagination.After)
		pagination.After = int32(after)
	}
	pagination.First = NullToInt32(input.Pagination.First)

	res, err := jobs.GetAppliedJobs(ctx, &pagination)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]JobPostingResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_jobPostingToResolver(ctx, o)
	}
	return list, nil
}

func (_ *Resolver) GetInvitedJobs(ctx context.Context, input GetAppliedJobsRequest) ([]JobPostingResolver, error) {
	pagination := jobsRPC.Pagination{}

	if input.Pagination.After != nil {
		after, _ := strconv.Atoi(*input.Pagination.After)
		pagination.After = int32(after)
	}
	pagination.First = NullToInt32(input.Pagination.First)

	res, err := jobs.GetInvitedJobs(ctx, &pagination)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]JobPostingResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_jobPostingToResolver(ctx, o)
	}
	return list, nil
}

func (_ *Resolver) GetAmountsOfManageCandidates(ctx context.Context, input GetAmountsOfManageCandidatesRequest) (*AmountsOfManageCandidatesResolver, error) {
	res, err := jobs.GetAmountsOfManageCandidates(ctx, &jobsRPC.ID{
		Id: input.Company_id,
	})
	if err != nil {
		return nil, err
	}

	return &AmountsOfManageCandidatesResolver{
		R: &AmountsOfManageCandidates{
			Alerts:  res.GetAlerts(),
			Saved:   res.GetSaved(),
			Skipped: res.GetSkipped(),
		},
	}, nil
}

func (_ *Resolver) GetAmountOfApplicantsPerCategory(ctx context.Context, input GetAmountOfApplicantsPerCategoryRequest) (AmountOfApplicantsPerCategoryResolver, error) {
	res, err := jobs.GetAmountOfApplicantsPerCategory(ctx, &jobsRPC.ID{
		Id: input.Company_id,
	})
	if err != nil {
		return AmountOfApplicantsPerCategoryResolver{}, err
	}

	return AmountOfApplicantsPerCategoryResolver{
		R: &AmountOfApplicantsPerCategory{
			Total:        res.GetTotal(),
			Disqualified: res.GetDisqualified(),
			Favorite:     res.GetFavorite(),
			In_review:    res.GetInReview(),
			Unseen:       res.GetUnseen(),
		},
	}, nil
}

func (_ *Resolver) GetCandidates(ctx context.Context, input GetCandidatesRequest) ([]CandidateProfileResolver, error) {
	pagination := jobsRPC.Pagination{}

	if input.Pagination != nil {
		if input.Pagination.After != nil {
			after, _ := strconv.Atoi(*input.Pagination.After)
			pagination.After = int32(after)
		}
		pagination.First = NullToInt32(input.Pagination.First)
	}

	res, err := jobs.GetCandidates(ctx, &jobsRPC.PaginationWithId{
		Id:         input.CompanyId,
		Pagination: &pagination,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]CandidateProfileResolver, len(res.List))
	for i, o := range res.List {
		list[i] = CandidateProfileResolver{R: jobs_candidateProfileToGql(ctx, o)}
	}
	return list, nil
}

/*func (_ *Resolver) SearchCandidates(ctx context.Context, input SearchCandidatesRequest) ([]CandidateProfileResolver, error) {
	res, err := jobs.SearchCandidates(ctx, &jobsRPC.CompanyIdWithCandidateSearchFilter{
		CompanyId: input.CompanyId,
		Filter:    jobs_candidateSearchFilterInputToRPC(&input.Filter),
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]CandidateProfileResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_candidateProfileToResolver(o)
	}
	return list, nil
}*/

func (_ *Resolver) GetSavedCandidates(ctx context.Context, input GetSavedCandidatesRequest) ([]CandidateProfileResolver, error) {
	res, err := jobs.GetSavedCandidates(ctx, &jobsRPC.PaginationWithId{Id: input.CompanyId, Pagination: &jobsRPC.Pagination{
		First: NullToInt32(input.First),
		After: NullToInt32(input.After),
	}})
	if e, isErr := handleError(err); isErr {
		return []CandidateProfileResolver{}, e
	}
	list := make([]CandidateProfileResolver, len(res.List))
	for i, o := range res.List {
		list[i].R = jobs_candidateProfileToGql(ctx, o)
	}

	return list, nil
}

func (_ *Resolver) GetSkippedCandidates(ctx context.Context, input GetSkippedCandidatesRequest) ([]CandidateProfileResolver, error) {
	res, err := jobs.GetSkippedCandidates(ctx, &jobsRPC.PaginationWithId{Id: input.CompanyId, Pagination: &jobsRPC.Pagination{
		First: NullToInt32(input.First),
		After: NullToInt32(input.After),
	}})
	if e, isErr := handleError(err); isErr {
		return []CandidateProfileResolver{}, e
	}
	list := make([]CandidateProfileResolver, len(res.List))
	for i, o := range res.List {
		list[i].R = jobs_candidateProfileToGql(ctx, o)
	}

	return list, nil
}

func (_ *Resolver) SetOpenFlag(ctx context.Context, input SetOpenFlagRequest) (*bool, error) {
	_, err := jobs.SetOpenFlag(ctx, &jobsRPC.BoolValue{Value: input.Open})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SetCareerInterests(ctx context.Context, input SetCareerInterestsRequest) (*bool, error) {
	_, err := jobs.SetCareerInterests(ctx, jobs_careerInterestsInputToRPC(&input.Interests))
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetPlanPrices(ctx context.Context, input GetPlanPricesRequest) (*[]PlanPriceResolver, error) {
	result, err := jobs.GetPlanPrices(ctx, &jobsRPC.GetPlanPricesRequest{
		CompanyID: input.Company_id,
		Countries: input.Countries,
		Currency:  input.Currency,
	})
	if err != nil {
		return nil, err
	}

	planPrices := make([]PlanPriceResolver, 0, len(result.GetPrices()))

	for _, r := range result.GetPrices() {
		planPrices = append(planPrices, PlanPriceResolver{
			R: &PlanPrice{
				Country:  r.GetCountry(),
				Currency: r.GetCurrency(),
				Features: Features{
					Anonymously: float64(r.GetFeatures().GetAnonymously()),
					Language:    float64(r.GetFeatures().GetLanguage()),
					Renewal:     float32ArrayToFloat64Array(r.GetFeatures().GetRenewal()),
				},
				Price_per_plan: PricePerPlan{
					Basic:            float64(r.GetPricesPerPlan().GetBasic()),
					Exclusive:        float64(r.GetPricesPerPlan().GetExclusive()),
					Premium:          float64(r.GetPricesPerPlan().GetPremium()),
					Professional:     float64(r.GetPricesPerPlan().GetProfessional()),
					ProfessionalPlus: float64(r.GetPricesPerPlan().GetProfessionalPlus()),
					Standard:         float64(r.GetPricesPerPlan().GetStandard()),
					Start:            float64(r.GetPricesPerPlan().GetStart()),
				},
			},
		})
	}

	return &planPrices, nil
}

func (_ *Resolver) GetPricingFor(ctx context.Context, input GetPricingForRequest) (TotalPricingResultResolver, error) {
	res, err := jobs.GetPricingFor(ctx, &jobsRPC.GetPricingRequest{
		CompanyId: input.CompanyId,
		Meta: &jobsRPC.JobMeta{
			// JobPlan:                jobsRPC.JobPlan(jobsRPC.JobPlan_value[input.Meta.Job_plan]),
			AmountOfDays:           input.Meta.Amount_of_days,
			Anonymous:              input.Meta.Anonymous,
			Renewal:                input.Meta.Renewal,
			AdvertisementCountries: input.Meta.Advertisement_countries,
			NumOfLanguages:         input.Meta.Num_of_languages,
			Currency:               input.Meta.Currency,
		},
	})
	if e, isErr := handleError(err); isErr {
		return TotalPricingResultResolver{}, e
	}
	result := TotalPricingResultResolver{
		R: &TotalPricingResult{
			Total:     float64(res.Total),
			Currency:  res.Currency,
			Countries: make([]*PricingResultByCountry, len(res.Countries)),
		},
	}
	for i, pr := range res.Countries {
		result.R.Countries[i] = &PricingResultByCountry{
			Country:                   pr.Country,
			Language_price:            float64(pr.LanguagePrice),
			Plan_price:                float64(pr.PlanPrice),
			Publish_anonymously_price: float64(pr.PublishAnonymouslyPrice),
			Renewal_price:             float64(pr.RenewalPrice),
			Total_price:               float64(pr.TotalPrice),
		}
	}

	return result, nil
}

func (_ *Resolver) PostJob(ctx context.Context, input PostJobRequest) (SuccessResolver, error) {
	details := input.Details
	jobTypes := make([]jobsRPC.JobType, len(details.Employment_types))
	for i, t := range details.Employment_types {
		jobTypes[i] = jobsRPC.JobType(jobsRPC.JobType_value[t])
	}
	descriptions := make([]*jobsRPC.JobDescription, len(details.Descriptions))
	for i, d := range details.Descriptions {
		descriptions[i] = &jobsRPC.JobDescription{
			Language:    d.Language,
			Description: d.Description,
			WhyUs:       d.Why_us,
		}
	}

	benefits := []jobsRPC.JobDetails_JobBenefit{}

	if details.Benefits != nil {
		benefits = make([]jobsRPC.JobDetails_JobBenefit, len(*details.Benefits))
		for i, t := range *details.Benefits {
			benefits[i] = jobsRPC.JobDetails_JobBenefit(jobsRPC.JobDetails_JobBenefit_value[t])
		}
	}

	id, err := jobs.PostJob(ctx, &jobsRPC.PostJobRequest{
		CompanyId: input.CompanyId,
		Metadata: &jobsRPC.JobMeta{
			// JobPlan:                jobsRPC.JobPlan(jobsRPC.JobPlan_value[input.Meta.Job_plan]),
			AmountOfDays:           input.Meta.Amount_of_days,
			Anonymous:              input.Meta.Anonymous,
			Highlight:              stringToHighlightEnum(input.Meta.Highlight),
			Renewal:                input.Meta.Renewal,
			AdvertisementCountries: input.Meta.Advertisement_countries,
			NumOfLanguages:         input.Meta.Num_of_languages,
			Currency:               input.Meta.Currency,
		},
		Details: &jobsRPC.JobDetails{
			Title:                  details.Title,
			Country:                NullToString(details.Country),
			Region:                 NullToString(details.Region),
			City:                   NullToString(details.City),
			LocationType:           stringLocationTypeToEnum(details.Location_type),
			AdditionalInfo:         additionalInfoToRPC(details.Additional_info),
			JobFunctions:           jobFunctionArrayToRPC(details.Job_functions),
			EmploymentTypes:        jobTypes,
			Descriptions:           descriptions,
			Required:               qualificationInputToRPC(details.Required),
			Preterred:              qualificationInputToRPC(details.Preterred),
			SalaryCurrency:         NullToString(details.Salary_currency),
			SalaryMin:              NullToInt32(details.Salary_min),
			SalaryMax:              NullToInt32(details.Salary_max),
			AdditionalCompensation: AddtionCompensationToRPC(details.Additional_compensation),
			SalaryInterval:         jobsRPC.SalaryInterval(jobsRPC.SalaryInterval_value[NullToString(details.Salary_interval)]),
			Benefits:               benefits, // TODO:
			NumberOfPositions:      details.Number_of_positions,
			CoverLetter:            details.Cover_letter,
			IsWillingToWorkRemotly: details.Is_willing_to_work_remotly,
		},
	})
	if e, isErr := handleError(err); isErr {
		return SuccessResolver{
			R: &Success{
				Success: false,
			},
		}, e
	}
	return SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetId(),
		},
	}, nil
}

func (_ *Resolver) ChangePost(ctx context.Context, input ChangePostRequest) (SuccessResolver, error) {
	details := input.Details
	jobTypes := make([]jobsRPC.JobType, len(details.Employment_types))
	for i, t := range details.Employment_types {
		jobTypes[i] = jobsRPC.JobType(jobsRPC.JobType_value[t])
	}
	descriptions := make([]*jobsRPC.JobDescription, len(details.Descriptions))
	for i, d := range details.Descriptions {
		descriptions[i] = &jobsRPC.JobDescription{
			Language:    d.Language,
			Description: d.Description,
			WhyUs:       d.Why_us,
		}
	}
	// var benefits []jobsRPC.JobDetails_JobBenefit

	// if input.Details.Benefits != nil {
	// 	benefits := make([]jobsRPC.JobDetails_JobBenefit, 0, len(*input.Details.Benefits))
	// 	for _, b := range *(input.Details.Benefits) {
	// 		benefits = append(benefits, jobBenefitToString(b))
	// 	}
	// }

	benefits := make([]jobsRPC.JobDetails_JobBenefit, len(*details.Benefits))
	for i, t := range *details.Benefits {
		benefits[i] = jobsRPC.JobDetails_JobBenefit(jobsRPC.JobDetails_JobBenefit_value[t])
	}
	id, err := jobs.ChangePost(ctx, &jobsRPC.PostJobRequest{
		CompanyId: input.CompanyId,
		DraftId:   input.DraftId,
		Details: &jobsRPC.JobDetails{
			Title:                  details.Title,
			Country:                NullToString(details.Country),
			Region:                 NullToString(details.Region),
			City:                   NullToString(details.City),
			JobFunctions:           jobFunctionArrayToRPC(details.Job_functions),
			EmploymentTypes:        jobTypes,
			Descriptions:           descriptions,
			Required:               qualificationInputToRPC(details.Required),
			Preterred:              qualificationInputToRPC(details.Preterred),
			AdditionalCompensation: AddtionCompensationToRPC(details.Additional_compensation),
			SalaryCurrency:         NullToString(details.Salary_currency),
			SalaryMin:              NullToInt32(details.Salary_min),
			SalaryMax:              NullToInt32(details.Salary_max),
			SalaryInterval:         jobsRPC.SalaryInterval(jobsRPC.SalaryInterval_value[NullToString(details.Salary_interval)]),
			Benefits:               benefits, // TODO:
			NumberOfPositions:      details.Number_of_positions,
			CoverLetter:            details.Cover_letter,
			IsWillingToWorkRemotly: details.Is_willing_to_work_remotly,
		},
	})
	if e, isErr := handleError(err); isErr {
		return SuccessResolver{
			R: &Success{
				Success: false,
			},
		}, e
	}
	return SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetId(),
		},
	}, nil
}

func (_ *Resolver) DeleteExpiredPost(ctx context.Context, in DeleteExpiredPostRequest) (*SuccessResolver, error) {
	_, err := jobs.DeleteExpiredPost(ctx, &jobsRPC.DeleteExpiredPostRequest{
		Id:     in.CompanyId,
		PostID: in.PostId,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) GetPost(ctx context.Context, input GetPostRequest) (JobDetailsResolver, error) {
	res, err := jobs.GetDraft(ctx, &jobsRPC.PostIDs{
		PostID:    input.DraftId,
		CompanyID: input.CompanyId,
	})
	if e, isErr := handleError(err); isErr {
		return JobDetailsResolver{}, e
	}

	details := jobs_jobDetailsToGql(res.GetDetails())

	// details := JobDetails{
	// 	Benefits:             make([]string, 0, len(res.GetDetails().GetBenefits())), // res.GetDetails().GetBenefits(),
	// 	City:                 res.GetDetails().GetCity(),
	// 	Country:              res.GetDetails().GetCountry(),
	// 	Cover_letter:         res.GetDetails().GetCoverLetter(),
	// 	Deadline_day:         res.GetDetails().GetDeadlineDay(),
	// 	Deadline_month:       res.GetDetails().GetDeadlineMonth(),
	// 	Deadline_year:        res.GetDetails().GetDeadlineYear(),
	// 	Descriptions:         make([]JobDescription, 0, len(res.GetDetails().GetDescriptions())), // res.GetDetails().GetBenefits(),
	// 	Employment_types:     make([]string, 0, len(res.GetDetails().GetEmploymentTypes())),
	// 	Header_url:           res.GetDetails().GetHeaderUrl(),
	// 	Hiring_day:           res.GetDetails().GetHiringDay(),
	// 	Hiring_month:         res.GetDetails().GetHiringMonth(),
	// 	Hiring_year:          res.GetDetails().GetHiringYear(),
	// 	Job_functions:        res.GetDetails().GetJobFunctions(),
	// 	Number_of_positions:  res.GetDetails().GetNumberOfPositions(),
	// 	Publish_day:          res.GetDetails().GetPublishDay(),
	// 	Publish_month:        res.GetDetails().GetPublishMonth(),
	// 	Publish_year:         res.GetDetails().GetPublishYear(),
	// 	Region:               res.GetDetails().GetRegion(),
	// 	Required_educations:  res.GetDetails().GetRequiredEducations(),
	// 	Required_eligibility: res.GetDetails().GetRequiredEligibility(),
	// 	Required_experience:  CareerInterestExperienceEnumRPC(res.GetDetails().GetRequiredExperience()),
	// 	Required_languages:   res.GetDetails().GetRequiredLanguages(),
	// 	Required_licenses:    res.GetDetails().GetRequiredLicenses(),
	// 	Required_skills:      res.GetDetails().GetRequiredSkills(),
	// 	Salary_currency:      res.GetDetails().GetSalaryCurrency(),
	// 	Salary_interval:      jobsRPCSalaryIntervalToString(res.GetDetails().GetSalaryInterval()),
	// 	Salary_max:           res.GetDetails().GetSalaryMax(),
	// 	Salary_min:           res.GetDetails().GetSalaryMin(),
	// 	Title:                res.GetDetails().GetTitle(),
	// }

	for _, b := range res.GetDetails().GetBenefits() {
		details.Benefits = append(details.Benefits, jobBenefitRPCToJobBenefit(b))
	}

	for _, d := range res.GetDetails().GetDescriptions() {
		details.Descriptions = append(details.Descriptions, jobsRPCDescriptionToJBDSCR(d))
	}

	for _, e := range res.GetDetails().GetEmploymentTypes() {
		details.Employment_types = append(details.Employment_types, jobsRPCJobTypeToString(e))
	}

	return JobDetailsResolver{
		R: details,
	}, nil
}

func (_ *Resolver) GetDraft(ctx context.Context, input GetDraftRequest) (JobInfoResolver, error) {
	res, err := jobs.GetDraft(ctx, &jobsRPC.PostIDs{
		PostID:    input.DraftId,
		CompanyID: input.CompanyId,
	})
	if e, isErr := handleError(err); isErr {
		return JobInfoResolver{}, e
	}

	details := jobs_jobDetailsToGql(res.GetDetails())

	// details := JobDetails{
	// 	Benefits:             make([]string, 0, len(res.GetDetails().GetBenefits())), // res.GetDetails().GetBenefits(),
	// 	City:                 res.GetDetails().GetCity(),
	// 	Country:              res.GetDetails().GetCountry(),
	// 	Cover_letter:         res.GetDetails().GetCoverLetter(),
	// 	Deadline_day:         res.GetDetails().GetDeadlineDay(),
	// 	Deadline_month:       res.GetDetails().GetDeadlineMonth(),
	// 	Deadline_year:        res.GetDetails().GetDeadlineYear(),
	// 	Descriptions:         make([]JobDescription, 0, len(res.GetDetails().GetDescriptions())), // res.GetDetails().GetBenefits(),
	// 	Employment_types:     make([]string, 0, len(res.GetDetails().GetEmploymentTypes())),
	// 	Header_url:           res.GetDetails().GetHeaderUrl(),
	// 	Hiring_day:           res.GetDetails().GetHiringDay(),
	// 	Hiring_month:         res.GetDetails().GetHiringMonth(),
	// 	Hiring_year:          res.GetDetails().GetHiringYear(),
	// 	Job_functions:        res.GetDetails().GetJobFunctions(),
	// 	Number_of_positions:  res.GetDetails().GetNumberOfPositions(),
	// 	Publish_day:          res.GetDetails().GetPublishDay(),
	// 	Publish_month:        res.GetDetails().GetPublishMonth(),
	// 	Publish_year:         res.GetDetails().GetPublishYear(),
	// 	Region:               res.GetDetails().GetRegion(),
	// 	Required_educations:  res.GetDetails().GetRequiredEducations(),
	// 	Required_eligibility: res.GetDetails().GetRequiredEligibility(),
	// 	Required_experience:  CareerInterestExperienceEnumRPC(res.GetDetails().GetRequiredExperience()),
	// 	Required_languages:   res.GetDetails().GetRequiredLanguages(),
	// 	Required_licenses:    res.GetDetails().GetRequiredLicenses(),
	// 	Required_skills:      res.GetDetails().GetRequiredSkills(),
	// 	Salary_currency:      res.GetDetails().GetSalaryCurrency(),
	// 	Salary_interval:      jobsRPCSalaryIntervalToString(res.GetDetails().GetSalaryInterval()),
	// 	Salary_max:           res.GetDetails().GetSalaryMax(),
	// 	Salary_min:           res.GetDetails().GetSalaryMin(),
	// 	Title:                res.GetDetails().GetTitle(),
	// }
	meta := JobMeta{
		Advertisement_countries: res.GetMetadata().GetAdvertisementCountries(),
		Renewal:                 res.GetMetadata().GetRenewal(),
		// Job_plan:                jobPlanRPCToString(res.GetMetadata().GetJobPlan()),
		Amount_of_days:   res.GetMetadata().GetAmountOfDays(),
		Anonymous:        res.GetMetadata().GetAnonymous(),
		Num_of_languages: res.GetMetadata().GetNumOfLanguages(),
		Currency:         res.GetMetadata().GetCurrency(),
	}

	for _, b := range res.GetDetails().GetBenefits() {
		details.Benefits = append(details.Benefits, jobBenefitRPCToJobBenefit(b))
	}

	for _, d := range res.GetDetails().GetDescriptions() {
		details.Descriptions = append(details.Descriptions, jobsRPCDescriptionToJBDSCR(d))
	}

	for _, e := range res.GetDetails().GetEmploymentTypes() {
		details.Employment_types = append(details.Employment_types, jobsRPCJobTypeToString(e))
	}

	ji := JobInfo{
		// Details: details,
		Meta: meta,
	}

	if details != nil {
		ji.Details = *details
	}

	return JobInfoResolver{
		R: &ji,
	}, nil
}

func (_ *Resolver) SaveDraft(ctx context.Context, input SaveDraftRequest) (SuccessResolver, error) {
	details := input.Details
	jobTypes := make([]jobsRPC.JobType, len(details.Employment_types))
	for i, t := range details.Employment_types {
		jobTypes[i] = jobsRPC.JobType(jobsRPC.JobType_value[t])
	}
	descriptions := make([]*jobsRPC.JobDescription, len(details.Descriptions))
	for i, d := range details.Descriptions {
		descriptions[i] = &jobsRPC.JobDescription{
			Language:    d.Language,
			Description: d.Description,
			WhyUs:       d.Why_us,
		}
	}
	// var benefits []jobsRPC.JobDetails_JobBenefit

	// if input.Details.Benefits != nil {
	// 	benefits := make([]jobsRPC.JobDetails_JobBenefit, 0, len(*input.Details.Benefits))
	// 	for _, b := range *(input.Details.Benefits) {
	// 		benefits = append(benefits, jobBenefitToString(b))
	// 	}
	// }
	benefits := make([]jobsRPC.JobDetails_JobBenefit, len(*details.Benefits))
	for i, t := range *details.Benefits {
		benefits[i] = jobsRPC.JobDetails_JobBenefit(jobsRPC.JobDetails_JobBenefit_value[t])
	}

	id, err := jobs.SaveDraft(ctx, &jobsRPC.PostJobRequest{
		CompanyId: input.CompanyId,
		Metadata: &jobsRPC.JobMeta{
			// JobPlan:                jobsRPC.JobPlan(jobsRPC.JobPlan_value[input.Meta.Job_plan]),
			AmountOfDays:           input.Meta.Amount_of_days,
			Anonymous:              input.Meta.Anonymous,
			Renewal:                input.Meta.Renewal,
			AdvertisementCountries: input.Meta.Advertisement_countries,
			NumOfLanguages:         input.Meta.Num_of_languages,
			Currency:               input.Meta.Currency,
		},
		Details: &jobsRPC.JobDetails{
			Title:                  details.Title,
			Country:                NullToString(details.Country),
			Region:                 NullToString(details.Region),
			City:                   NullToString(details.City),
			JobFunctions:           jobFunctionArrayToRPC(details.Job_functions),
			AdditionalInfo:         additionalInfoToRPC(details.Additional_info),
			AdditionalCompensation: AddtionCompensationToRPC(details.Additional_compensation),
			EmploymentTypes:        jobTypes,
			Descriptions:           descriptions,
			Required:               qualificationInputToRPC(details.Required),
			Preterred:              qualificationInputToRPC(details.Preterred),
			SalaryCurrency:         NullToString(details.Salary_currency),
			SalaryMin:              NullToInt32(details.Salary_min),
			SalaryMax:              NullToInt32(details.Salary_max),
			SalaryInterval:         jobsRPC.SalaryInterval(jobsRPC.SalaryInterval_value[NullToString(details.Salary_interval)]),
			Benefits:               benefits, // TODO:
			NumberOfPositions:      details.Number_of_positions,
			CoverLetter:            details.Cover_letter,
			IsWillingToWorkRemotly: details.Is_willing_to_work_remotly,
		},
	})
	if e, isErr := handleError(err); isErr {
		return SuccessResolver{
			R: &Success{
				Success: false,
			},
		}, e
	}
	return SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetId(),
		},
	}, nil
}

func (_ *Resolver) ChangeDraft(ctx context.Context, input ChangeDraftRequest) (SuccessResolver, error) {
	details := input.Details
	jobTypes := make([]jobsRPC.JobType, len(details.Employment_types))
	for i, t := range details.Employment_types {
		jobTypes[i] = jobsRPC.JobType(jobsRPC.JobType_value[t])
	}
	descriptions := make([]*jobsRPC.JobDescription, len(details.Descriptions))
	for i, d := range details.Descriptions {
		descriptions[i] = &jobsRPC.JobDescription{
			Language:    d.Language,
			Description: d.Description,
			WhyUs:       d.Why_us,
		}
	}
	// var benefits []jobsRPC.JobDetails_JobBenefit

	// if input.Details.Benefits != nil {
	// 	benefits := make([]jobsRPC.JobDetails_JobBenefit, 0, len(*input.Details.Benefits))
	// 	for _, b := range *(input.Details.Benefits) {
	// 		benefits = append(benefits, jobBenefitToString(b))
	// 	}
	// }
	benefits := make([]jobsRPC.JobDetails_JobBenefit, len(*details.Benefits))
	for i, t := range *details.Benefits {
		benefits[i] = jobsRPC.JobDetails_JobBenefit(jobsRPC.JobDetails_JobBenefit_value[t])
	}
	id, err := jobs.ChangeDraft(ctx, &jobsRPC.PostJobRequest{
		CompanyId: input.CompanyId,
		DraftId:   input.DraftId,
		Metadata: &jobsRPC.JobMeta{
			// JobPlan:                jobsRPC.JobPlan(jobsRPC.JobPlan_value[input.Meta.Job_plan]),
			AmountOfDays:           input.Meta.Amount_of_days,
			Anonymous:              input.Meta.Anonymous,
			Renewal:                input.Meta.Renewal,
			AdvertisementCountries: input.Meta.Advertisement_countries,
			NumOfLanguages:         input.Meta.Num_of_languages,
			Currency:               input.Meta.Currency,
		},
		Details: &jobsRPC.JobDetails{
			Title:                  details.Title,
			Country:                NullToString(details.Country),
			Region:                 NullToString(details.Region),
			City:                   NullToString(details.City),
			JobFunctions:           jobFunctionArrayToRPC(details.Job_functions),
			AdditionalInfo:         additionalInfoToRPC(details.Additional_info),
			AdditionalCompensation: AddtionCompensationToRPC(details.Additional_compensation),
			EmploymentTypes:        jobTypes,
			Descriptions:           descriptions,
			SalaryCurrency:         NullToString(details.Salary_currency),
			SalaryMin:              NullToInt32(details.Salary_min),
			SalaryMax:              NullToInt32(details.Salary_max),
			SalaryInterval:         jobsRPC.SalaryInterval(jobsRPC.SalaryInterval_value[NullToString(details.Salary_interval)]),
			Benefits:               benefits, // TODO:
			NumberOfPositions:      details.Number_of_positions,
			CoverLetter:            details.Cover_letter,
			IsWillingToWorkRemotly: details.Is_willing_to_work_remotly,
		},
	})
	if e, isErr := handleError(err); isErr {
		return SuccessResolver{
			R: &Success{
				Success: false,
			},
		}, e
	}
	return SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetId(),
		},
	}, nil
}

func (_ *Resolver) ActivateJob(ctx context.Context, input ActivateJobRequest) (*bool, error) {
	_, err := jobs.ActivateJob(ctx, &jobsRPC.CompanyIdWithJobId{CompanyId: input.CompanyId, JobId: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) PauseJob(ctx context.Context, input PauseJobRequest) (*bool, error) {
	_, err := jobs.PauseJob(ctx, &jobsRPC.CompanyIdWithJobId{CompanyId: input.CompanyId, JobId: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) ApplyJob(ctx context.Context, input ApplyJobRequest) (*bool, error) {

	req := jobsRPC.ApplyJobRequest{
		JobId:       input.Application.Job_id,
		Email:       input.Application.Email,
		Phone:       input.Application.Phone,
		CoverLetter: input.Application.Cover_letter,
		// Documents:   input.Application.Document_id,
	}

	if input.Application.Document_id != nil {
		req.Documents = *input.Application.Document_id
	}

	_, err := jobs.ApplyJob(ctx, &req)

	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) IgnoreInvitation(ctx context.Context, input IgnoreInvitationRequest) (*bool, error) {
	_, err := jobs.IgnoreInvitation(ctx, &jobsRPC.ID{
		Id: input.Job_id,
	})

	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) AddJobView(ctx context.Context, input AddJobViewRequest) (*bool, error) {
	_, err := jobs.AddJobView(ctx, &jobsRPC.ID{Id: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SaveJob(ctx context.Context, input SaveJobRequest) (*bool, error) {
	_, err := jobs.SaveJob(ctx, &jobsRPC.ID{Id: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnsaveJob(ctx context.Context, input UnsaveJobRequest) (*bool, error) {
	_, err := jobs.UnsaveJob(ctx, &jobsRPC.ID{Id: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SkipJob(ctx context.Context, input SkipJobRequest) (*bool, error) {
	_, err := jobs.SkipJob(ctx, &jobsRPC.ID{Id: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnskipJob(ctx context.Context, input UnskipJobRequest) (*bool, error) {
	_, err := jobs.UnskipJob(ctx, &jobsRPC.ID{Id: input.JobId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SetJobApplicationSeen(ctx context.Context, input SetJobApplicationSeenRequest) (*bool, error) {
	_, err := jobs.SetJobApplicationSeen(ctx, &jobsRPC.SetJobApplicationSeenRequest{
		CompanyId:   input.CompanyId,
		JobId:       input.JobId,
		ApplicantId: input.ApplicationId,
		Seen:        input.Seen,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SetJobApplicationCategory(ctx context.Context, input SetJobApplicationCategoryRequest) (*bool, error) {
	_, err := jobs.SetJobApplicationCategory(ctx, &jobsRPC.SetJobApplicationCategoryRequest{
		CompanyId:   input.CompanyId,
		JobId:       input.JobId,
		ApplicantId: input.ApplicationId,
		Category:    stringToJobsRPCApplicantCategoryEnum(input.Category),
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SaveCandidate(ctx context.Context, input SaveCandidateRequest) (*bool, error) {
	_, err := jobs.SaveCandidate(ctx, &jobsRPC.CompanyIdWithCandidateId{
		CompanyId:   input.CompanyId,
		CandidateId: input.CandidateId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnsaveCandidate(ctx context.Context, input UnsaveCandidateRequest) (*bool, error) {
	_, err := jobs.UnsaveCandidate(ctx, &jobsRPC.CompanyIdWithCandidateId{
		CompanyId:   input.CompanyId,
		CandidateId: input.CandidateId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SkipCandidate(ctx context.Context, input SkipCandidateRequest) (*bool, error) {
	_, err := jobs.SkipCandidate(ctx, &jobsRPC.CompanyIdWithCandidateId{
		CompanyId:   input.CompanyId,
		CandidateId: input.CandidateId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnskipCandidate(ctx context.Context, input UnskipCandidateRequest) (*bool, error) {
	_, err := jobs.UnskipCandidate(ctx, &jobsRPC.CompanyIdWithCandidateId{
		CompanyId:   input.CompanyId,
		CandidateId: input.CandidateId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) InviteUserToApply(ctx context.Context, input InviteUserToApplyRequest) (*bool, error) {
	_, err := jobs.InviteUserToApply(ctx, &jobsRPC.InviteUserToApplyRequest{
		CompanyId:      input.CompanyId,
		UserId:         input.UserId,
		JobId:          input.JobId,
		InvitationText: input.Text,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) ReportJob(ctx context.Context, input ReportJobRequest) (*bool, error) {
	_, err := jobs.ReportJob(ctx, &jobsRPC.ReportJobRequest{
		JobId: input.JobId,
		Type:  stringToJobsRPCReportJobTypeEnum(input.Type),
		Text:  NullToString(input.Text),
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) ReportCandidate(ctx context.Context, input ReportCandidateRequest) (*bool, error) {
	_, err := jobs.ReportCandidate(ctx, &jobsRPC.ReportCandidateRequest{
		CompanyId:   input.CompanyId,
		CandidateId: input.CandidateId,
		Text:        input.Text,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

// func (_ *Resolver) SaveJobSearchFilter(ctx context.Context, input SaveJobSearchFilterRequest) (*bool, error) {
// 	_, err := jobs.SaveJobSearchFilter(ctx, &jobsRPC.NamedJobSearchFilter{
// 		// Name:   input.Name,
// 		// Filter: jobs_jobSearchFilterInputToRPC(&input.Filter),
// 	})
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
// 	return nil, nil
// }

// func (_ *Resolver) SaveJobAlert(ctx context.Context, input SaveJobAlertRequest) (*bool, error) {
// 	_, err := jobs.SaveJobAlert(ctx, &jobsRPC.JobAlert{
// 		Name:               input.Name,
// 		Interval:           input.Interval,
// 		NotifyEmail:        input.Notify_email,
// 		NotifyNotification: input.Notify_notification,
// 		Filter:             jobs_jobSearchFilterInputToRPC(&input.Filter),
// 	})
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
// 	return nil, nil
// }

func (_ *Resolver) SaveCandidateSearchFilter(ctx context.Context, input SaveCandidateSearchFilterRequest) (*bool, error) {
	_, err := jobs.SaveCandidateSearchFilter(ctx, &jobsRPC.SaveCandidateSearchFilterRequest{
		CompanyId: input.CompanyId,
		Filter: &jobsRPC.NamedCandidateSearchFilter{
			Name:   input.Name,
			Filter: jobs_candidateSearchFilterInputToRPC(&input.Filter),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

// func (_ *Resolver) SaveCandidateAlert(ctx context.Context, input SaveCandidateAlertRequest) (graphql.ID, error) {
// 	id, err := jobs.SaveCandidateAlert(ctx, &jobsRPC.SaveCandidateAlertRequest{
// 		CompanyId: input.CompanyId,
// 		Alert: &jobsRPC.CandidateAlert{
// 			Name:               input.Name,
// 			Interval:           input.Interval,
// 			NotifyEmail:        input.Notify_email,
// 			NotifyNotification: input.Notify_notification,
// 			Filter:             jobs_candidateSearchFilterInputToRPC(&input.Filter),
// 		},
// 	})
// 	if e, isErr := handleError(err); isErr {
// 		return "", e
// 	}
// 	return graphql.ID(id.Id), nil
// }
//
// func (_ *Resolver) UpdateCandidateAlert(ctx context.Context, input UpdateCandidateAlertRequest) (*bool, error) {
// 	_, err := jobs.SaveCandidateAlert(ctx, &jobsRPC.SaveCandidateAlertRequest{
// 		CompanyId: input.CompanyId,
// 		Alert: &jobsRPC.CandidateAlert{
// 			Id:                 input.AlertId,
// 			Name:               input.Name,
// 			Interval:           input.Interval,
// 			NotifyEmail:        input.Notify_email,
// 			NotifyNotification: input.Notify_notification,
// 			Filter:             jobs_candidateSearchFilterInputToRPC(&input.Filter),
// 		},
// 	})
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
// 	return nil, nil
// }
//
// func (_ *Resolver) DeleteCandidateAlert(ctx context.Context, input DeleteCandidateAlertRequest) (*bool, error) {
// 	_, err := jobs.DeleteCandidateAlert(ctx, &jobsRPC.CompanyIdWithId{
// 		CompanyId: input.CompanyId,
// 		Id:        input.AlertId,
// 	})
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
// 	return nil, nil
// }

/*func (_ *Resolver) SearchJob(ctx context.Context, input SearchJobRequest) ([]JobPostingResolver, error) {
	res, err := jobs.SearchJob(ctx, jobs_jobSearchFilterInputToRPC(&input.Filter))
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]JobPostingResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobs_jobPostingToResolver(o)
	}
	return list, nil
}*/

// func (_ *Resolver) GetSavedJobSearchFilters(ctx context.Context) ([]NamedJobSearchFilterResolver, error) {
// 	res, err := jobs.GetSavedJobSearchFilters(ctx, jobsEmpty)
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
// 	list := make([]NamedJobSearchFilterResolver, len(res.List))
// 	for i, o := range res.List {
// 		list[i] = *jobs_NamedJobSearchFilterToResolver(o)
// 	}
// 	return list, nil
// }

// func (_ *Resolver) GetJobAlerts(ctx context.Context) ([]JobAlertResolver, error) {
// 	res, err := jobs.GetJobAlerts(ctx, jobsEmpty)
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
// 	list := make([]JobAlertResolver, len(res.List))
// 	for i, o := range res.List {
// 		list[i] = *jobs_jobAlertToResolver(o)
// 	}
// 	return list, nil
// }

// func (_ *Resolver) GetSavedCandidateSearchFilters(ctx context.Context, input GetSavedCandidateSearchFiltersRequest) ([]NamedCandidateSearchFilterResolver, error) {
// 	res, err := jobs.GetSavedCandidateSearchFilters(ctx, &jobsRPC.ID{Id: input.CompanyId})
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
// 	list := make([]NamedCandidateSearchFilterResolver, len(res.List))
// 	for i, o := range res.List {
// 		list[i] = *jobs_NamedCandidateSearchFilterToResolver(o)
// 	}
// 	return list, nil
// }

// func (_ *Resolver) GetCandidateAlerts(ctx context.Context, input GetCandidateAlertsRequest) ([]CandidateAlertResolver, error) {
// 	res, err := jobs.GetCandidateAlerts(ctx, &jobsRPC.ID{Id: input.CompanyId})
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
// 	list := make([]CandidateAlertResolver, len(res.List))
// 	for i, o := range res.List {
// 		list[i] = *jobs_candidateAlertToResolver(o)
// 	}
// 	return list, nil
// }

func (_ *Resolver) GetListOfJobsWithSeenStat(ctx context.Context, input GetListOfJobsWithSeenStatRequest) (*[]JobWithSeenStatResolver, error) {
	var first, after int32
	if input.Pagination.After != nil {
		afterInt, _ := strconv.Atoi(*input.Pagination.After)
		after = int32(afterInt)
	}

	res, err := jobs.GetListOfJobsWithSeenStat(ctx, &jobsRPC.PaginationWithId{
		Id: input.Company_id,
		Pagination: &jobsRPC.Pagination{
			First: first,
			After: after,
		},
	})
	if err != nil {
		return nil, err
	}

	list := make([]JobWithSeenStatResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *jobsJobWithSeenStatToResolver(o)
	}
	return &list, nil
}

// AddCVInCareerCenter ...
func (_ *Resolver) AddCVInCareerCenter(ctx context.Context, input AddCVInCareerCenterRequest) (*bool, error) {
	_, err := jobs.AddCVInCareerCenter(ctx, &jobsRPC.AddCVInCareerCenterRequest{
		CompanyID:                input.Company_id,
		ExpierencedProfessionals: input.Options.ExpierencedProfessionals,
		NewJobSeekers:            input.Options.NewJobSeekers,
		YoungProfessionals:       input.Options.YoungProfessionals,
	})
	if err != nil {
		return nil, err
	}

	t := true
	return &t, nil
}

// GetSavedCVs ...
func (_ *Resolver) GetSavedCVs(ctx context.Context, input GetSavedCVsRequest) ([]CandidateProfileResolver, error) {
	pagination := jobsRPC.Pagination{}

	if input.Pagination != nil {
		if input.Pagination.After != nil {
			after, _ := strconv.Atoi(*input.Pagination.After)
			pagination.After = int32(after)
		}
		pagination.First = NullToInt32(input.Pagination.First)
	}

	res, err := jobs.GetSavedCVs(ctx, &jobsRPC.GetSavedCVsRequest{
		Id:         input.CompanyId,
		Pagination: &pagination,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]CandidateProfileResolver, len(res.List))
	for i, o := range res.List {
		list[i] = CandidateProfileResolver{R: jobs_candidateProfileToGql(ctx, o)}
	}
	return list, nil
}

// RemoveCVs ...
func (_ *Resolver) RemoveCVs(ctx context.Context, input RemoveCVsRequest) (*bool, error) {
	_, err := jobs.RemoveCVs(ctx, &jobsRPC.IDs{
		ID:  input.CompanyId,
		Ids: input.Ids,
	})
	if err != nil {
		return nil, err
	}

	t := true
	return &t, nil
}

// MakeFavoriteCVs ...
func (_ *Resolver) MakeFavoriteCVs(ctx context.Context, input MakeFavoriteCVsRequest) (*bool, error) {
	_, err := jobs.MakeFavoriteCVs(ctx, &jobsRPC.MakeFavoriteCVsRequest{
		ID:          input.CompanyId,
		IDs:         input.Ids,
		IsFavourite: input.Is_favourite,
	})
	if err != nil {
		return nil, err
	}

	t := true
	return &t, nil
}

func stringToJobsRPCReportJobTypeEnum(s string) jobsRPC.ReportJobRequest_ReportJobTypeEnum {
	switch s {
	case "scam":
		return jobsRPC.ReportJobRequest_scam
	case "offensive":
		return jobsRPC.ReportJobRequest_offensive
	case "incorrect":
		return jobsRPC.ReportJobRequest_incorrect
	case "expired":
		return jobsRPC.ReportJobRequest_expired
	}

	return jobsRPC.ReportJobRequest_other
}

func stringToJobExperienceEnum(s string) jobsRPC.ExperienceEnum {

	switch s {
	case "without_experience":
		return jobsRPC.ExperienceEnum_WithoutExperience
	case "less_then_one_year":
		return jobsRPC.ExperienceEnum_LessThenOneYear
	case "one_two_years":
		return jobsRPC.ExperienceEnum_OneTwoYears
	case "two_three_years":
		return jobsRPC.ExperienceEnum_TwoThreeYears
	case "three_five_years":
		return jobsRPC.ExperienceEnum_ThreeFiveYears
	case "five_seven_years":
		return jobsRPC.ExperienceEnum_FiveSevenyears
	case "seven_ten_years":
		return jobsRPC.ExperienceEnum_SevenTenYears
	case "ten_years_and_more":
		return jobsRPC.ExperienceEnum_TenYearsAndMore
	}
	return jobsRPC.ExperienceEnum_UnknownExperience
}
