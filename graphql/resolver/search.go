package resolver

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/searchRPC"
)

type Empty struct{}

func (_ *Resolver) SearchUsers(ctx context.Context, input SearchUsersRequest) (SearchUserResultResolver, error) {
	var first uint32
	var after string

	if input.Pagination.First != nil {
		first = uint32(*input.Pagination.First)
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	request := searchRPC.UserSearchRequest{
		First: first,
		After: after,

		Keywords:        NullStringArrayToStringArray(input.Input.Keywords),
		MyConnections:   input.Input.IsMyConnection,
		ConnectionsOfID: NullStringArrayToStringArray(input.Input.ConenctionsOf),
		CountryID:       NullStringArrayToStringArray(input.Input.Country),
		CityID:          NullStringArrayToStringArray(input.Input.City),
		School:          NullStringArrayToStringArray(input.Input.School),
		Degree:          NullStringArrayToStringArray(input.Input.Degree),
		FieldOfStudy:    NullStringArrayToStringArray(input.Input.FiledOfStudy),
		IsStudent:       input.Input.IsStudent,
		CurrentCompany:  NullStringArrayToStringArray(input.Input.CurrentCompany),
		PastCompany:     NullStringArrayToStringArray(input.Input.PastCompany),
		Industry:        NullStringArrayToStringArray(input.Input.Industry),
		Position:        NullStringArrayToStringArray(input.Input.Position),
		Firstname:       NullStringArrayToStringArray(input.Input.Firstname),
		Lastname:        NullStringArrayToStringArray(input.Input.Lastname),
		Nickname:        NullStringArrayToStringArray(input.Input.Nickname),
		IsMale:          input.Input.IsMale,
		IsFemale:        input.Input.IsFemale,
		Skill:           NullStringArrayToStringArray(input.Input.Skill),
		Language:        NullStringArrayToStringArray(input.Input.Language),
		Interest:        NullStringArrayToStringArray(input.Input.Interest),
		FullName:        NullToString(input.Input.Full_name),
		MinAge:          Nullint32ToUint32(input.Input.MinAge),
		MaxAge:          Nullint32ToUint32(input.Input.MaxAge),
	}

	// Changed Birthday to Age. so this is now redandent

	// if input.Input.Birthday != nil {
	// 	bds := make([]*searchRPC.Date, len(*input.Input.Birthday))
	// 	for i := range *input.Input.Birthday {
	// 		bds[i] = &searchRPC.Date{
	// 			Day:   Nullint32ToUint32((*input.Input.Birthday)[i].DayOfBirth),
	// 			Month: Nullint32ToUint32((*input.Input.Birthday)[i].MonthOfBirth),
	// 			Year:  Nullint32ToUint32((*input.Input.Birthday)[i].YearOfBirth),
	// 		}
	// 	}
	// 	request.Birthday = bds
	// }

	answer, _ := search.UserSearch(ctx, &request)

	pr := make([]Profile, len(answer.GetProfiles()))
	for i := range answer.GetProfiles() {
		pr[i] = ToProfile(ctx, answer.GetProfiles()[i])
	}

	return SearchUserResultResolver{
		R: &SearchUserResult{
			Amount_of_results: int32(answer.GetTotal()),
			Profiles:          pr,
		},
	}, nil
}

func (_ *Resolver) SearchJobs(ctx context.Context, input SearchJobsRequest) (SearchJobResultResolver, error) {
	var first uint32
	var after string

	if input.Pagination.First != nil {
		first = uint32(*input.Pagination.First)
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	request := searchRPC.JobSearchRequest{
		First: first,
		After: after,

		Keywords:           NullStringArrayToStringArray(input.Input.Keywords),
		Degree:             NullStringArrayToStringArray(input.Input.Degree),
		Country:            NullStringArrayToStringArray(input.Input.Country),
		City:               NullStringArrayToStringArray(input.Input.City),
		ExperienceLevel:    stringTotExperienceEnumRPC(input.Input.Experience_level),
		DatePosted:         stringToDateEnumRPC(input.Input.Date_posted),
		JobType:            NullStringArrayToStringArray(input.Input.Job_type),
		Language:           NullStringArrayToStringArray(input.Input.Language),
		Industry:           NullStringArrayToStringArray(input.Input.Industry),
		Subindustry:        NullStringArrayToStringArray(input.Input.Subindustry),
		CompanyName:        NullStringArrayToStringArray(input.Input.Company_name),
		CompanySize:        stringToSearchRPC(input.Input.Company_size),
		Currency:           NullToString(input.Input.Currency),
		Period:             input.Input.Period,
		Skill:              NullStringArrayToStringArray(input.Input.Skill),
		IsFollowing:        input.Input.Is_following,
		WithoutCoverLetter: input.Input.Without_cover_letter,
		WithSalary:         input.Input.With_salary,
	}

	// if input.Input.Date_posted != nil {
	// 	bds := make([]*searchRPC.Date, len(*input.Input.Date_posted))
	// 	for i := range *input.Input.Date_posted {
	// 		bds[i] = &searchRPC.Date{
	// 			Day:   Nullint32ToUint32((*input.Input.Date_posted)[i].DayOfBirth),
	// 			Month: Nullint32ToUint32((*input.Input.Date_posted)[i].MonthOfBirth),
	// 			Year:  Nullint32ToUint32((*input.Input.Date_posted)[i].YearOfBirth),
	// 		}
	// 	}
	// 	request.DatePosted = bds
	// }

	// if input.Input.Experience_level == nil {
	// 	request.IsExperienceLevelNull = true
	// }
	// else {
	// 	request.ExperienceLevel = Nullint32ToUint32(input.Input.Experience_level)
	// }

	if input.Input.Min_salary == nil {
		request.IsMinSalaryNull = true
	} else {
		request.MinSalary = uint32(*input.Input.Min_salary)
	}

	if input.Input.Max_salary == nil {
		request.IsMaxSalaryNull = true
	} else {
		request.MaxSalary = uint32(*input.Input.Max_salary)
	}

	if input.Input.Company_ids != nil {
		request.CompanyIDs = *input.Input.Company_ids
	}

	answer, err := search.JobSearch(ctx, &request)
	if err != nil {
		return SearchJobResultResolver{}, err
	}

	jobs := make([]JobSearchResult, 0, len(answer.GetJobResults()))

	for i := range answer.GetJobResults() {
		jobRes := JobSearchResult{
			ID: answer.GetJobResults()[i].GetID(),
			// Job_details: *jobs_jobDetailsToGql(answer.GetJobResults()[i].GetJob()),
			// Company:    toCompanyProfile(ctx, *),
			Is_saved:   answer.GetJobResults()[i].GetIsSaved(),
			Is_applied: answer.GetJobResults()[i].GetIsApplied(),
		}

		if jobs_jobDetailsToGql(answer.GetJobResults()[i].GetJob()) != nil {
			jobRes.Job_details = *jobs_jobDetailsToGql(answer.GetJobResults()[i].GetJob())
		}

		if answer.GetJobResults()[i].GetCompany() != nil {
			jobRes.Company = toCompanyProfile(ctx, *answer.GetJobResults()[i].GetCompany())
		}

		if jobRes.Company.ID == "" {
			log.Println("error: job", jobRes.ID, "has invalid company")
			continue
		}

		jobs = append(jobs, jobRes)
	}

	// for i := range jobs {
	// 	log.Println("company: ", jobs[i].Company.ID jobs[i].Company.Type, jobs[i].Company.Size)
	// }

	return SearchJobResultResolver{
		R: &SearchJobResult{
			Amount_of_results: int32(answer.GetTotal()),
			Job_search_result: jobs,
		}}, nil
}

func (_ *Resolver) SearchCandidate(ctx context.Context, input SearchCandidateRequest) (SearchCandidateResultResolver, error) {
	var first uint32
	var after string

	if input.Pagination.First != nil {
		first = uint32(*input.Pagination.First)
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	request := searchRPC.CandidateSearchRequest{
		First: first,
		After: after,

		Keywords:               NullStringArrayToStringArray(input.Input.Keywords),
		Country:                NullStringArrayToStringArray(input.Input.Country),
		City:                   NullStringArrayToStringArray(input.Input.City),
		CurrentCompany:         NullStringArrayToStringArray(input.Input.Current_company),
		PastCompany:            NullStringArrayToStringArray(input.Input.Past_company),
		Industry:               NullStringArrayToStringArray(input.Input.Industry),
		SubIndustry:            NullStringArrayToStringArray(input.Input.Sub_industry),
		JobType:                NullStringArrayToStringArray(input.Input.Job_type),
		Skill:                  NullStringArrayToStringArray(input.Input.Skill),
		Language:               NullStringArrayToStringArray(input.Input.Language),
		School:                 NullStringArrayToStringArray(input.Input.School),
		Degree:                 NullStringArrayToStringArray(input.Input.Degree),
		FieldOfStudy:           NullStringArrayToStringArray(input.Input.Field_of_study),
		IsStudent:              input.Input.Is_student,
		Currency:               NullToString(input.Input.Currency),
		Period:                 input.Input.Period,
		IsWillingToTravel:      input.Input.Is_willing_to_travel,
		IsWillingToWorkRemotly: input.Input.Is_willing_to_work_remotly,
		IsPossibleToRelocate:   input.Input.Is_possible_to_relocate,
		ExperienceLevel:        stringTotExperienceEnumRPC(input.Input.Experience_level),
	}

	if input.Company_id != nil {
		request.CompanyID = *input.Company_id
	}

	// if input.Input.Experience_level == nil {
	// 	request.IsExperienceLevelNull = true
	// } else {
	// 	request.ExperienceLevel = Nullint32ToUint32(input.Input.Experience_level)
	// }

	if input.Input.Min_salary == nil {
		request.IsMinSalaryNull = true
	} else {
		request.MinSalary = uint32(*input.Input.Min_salary)
	}

	if input.Input.Max_salary == nil {
		request.IsMaxSalaryNull = true
	} else {
		request.MaxSalary = uint32(*input.Input.Max_salary)
	}

	answer, err := search.CandidateSearch(ctx, &request)
	if err != nil {
		return SearchCandidateResultResolver{}, err
	}

	candidates := make([]JobApplicant, 0, len(answer.GetCandidateResults()))

	for i := range answer.GetCandidateResults() {
		j := new(CareerInterests)
		j.Company_size = "size_unknown"
		j.Experience = "experience_unknown"
		j.Salary_interval = "Unknown"

		if answer.GetCandidateResults() != nil && len(answer.GetCandidateResults()) > i {
			if answer.GetCandidateResults()[i].GetCareerInterests() != nil {
				j = jobs_careerInterestsToGql(answer.GetCandidateResults()[i].GetCareerInterests())

			}
		} /*else {
			j = CareerInterests{}
		}*/
		candidates = append(candidates, JobApplicant{
			User:             ToProfile(ctx, answer.GetCandidateResults()[i].GetCandidates()),
			Career_interests: j,
		})
	}

	return SearchCandidateResultResolver{
		R: &SearchCandidateResult{
			Amount_of_results:       int32(answer.GetTotal()),
			Candidate_search_result: candidates,
		},
	}, nil
}

func (_ *Resolver) SearchCompanies(ctx context.Context, input SearchCompaniesRequest) (SearchCompaniesResultResolver, error) {
	var first uint32
	var after string

	if input.Pagination.First != nil {
		first = uint32(*input.Pagination.First)
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	query := searchRPC.CompanySearchRequest{
		After:          after,
		First:          first,
		IsCompany:      input.Input.Search_for_companies,
		IsJobOffers:    input.Input.With_jobs,
		Keywords:       NullStringArrayToStringArray(input.Input.Keywords),
		Size:           stringToCompanySizeRPC(input.Input.Size),
		Type:           stringToCompanyTypeRPC(input.Input.Type),
		IsOrganization: input.Input.Search_for_organizations,
		Name:           NullStringArrayToStringArray(input.Input.Name),
		Industry:       NullStringArrayToStringArray(input.Input.Industry),
		SubIndustry:    NullStringArrayToStringArray(input.Input.Subindustry),
		City:           NullStringArrayToStringArray(input.Input.City),
		Country:        NullStringArrayToStringArray(input.Input.Country),
		FounderIDs:     NullStringArrayToStringArray(input.Input.Founders_id),
		FounderNames:   NullStringArrayToStringArray(input.Input.Founders_name),
		BusinessHours:  NullStringArrayToStringArray(input.Input.Business_hours),
		// Rating:         input.Input.Rating,
	}

	if input.Input.Is_career_center_opened != nil {
		query.IsCareerCenterOpenened = *input.Input.Is_career_center_opened
	}

	// if input.Input.Size != nil {
	// 	query.Size = make([]companyRPC.Size, 0, len(*input.Input.Size))
	// 	for _, s := range *input.Input.Size {
	// 		query.Size = append(query.Size, stringToCompanySizeRPC(s))
	// 	}
	// }

	// if input.Input.Type != nil {
	// 	query.Type = make([]companyRPC.Type, 0, len(*input.Input.Type))
	// 	for _, s := range *input.Input.Type {
	// 		query.Type = append(query.Type, stringToCompanyTypeRPC(s))
	// 	}
	// }

	result, err := search.CompanySearch(ctx, &query)
	if err != nil {
		return SearchCompaniesResultResolver{}, err
	}

	companies := make([]CompanyProfile, 0, len(result.GetResults()))
	for _, res := range result.GetResults() {
		if res != nil {
			companies = append(companies, toCompanyProfile(ctx, *res))
		}
	}

	return SearchCompaniesResultResolver{
		R: &SearchCompaniesResult{
			Amount_of_results: int32(result.GetTotal()),
			Company:           companies,
		},
	}, nil

}

func (r *Resolver) SearchServices(ctx context.Context, input SearchServicesRequest) (SearchServiceResultResolver, error) {
	var first uint32
	var after string

	if input.Pagination.First != nil {
		first = uint32(*input.Pagination.First)
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	query := searchRPC.ServiceSearchRequest{
		After:        after,
		First:        first,
		Keywords:     NullStringArrayToStringArray(input.Input.Keywords),
		City:         NullStringArrayToStringArray(input.Input.City),
		Country:      NullStringArrayToStringArray(input.Input.Country),
		Currency:     NullToString(input.Input.Currency),
		FixedPrice:   NullToInt32(input.Input.Fixed_price_amount),
		MinSalary:    NullToInt32(input.Input.Min_salary),
		MaxSalary:    NullToInt32(input.Input.Max_salary),
		LocationType: locationTypeToServiceRPCLocationTypeEnum(input.Input.Location_type),
		IsAlwaysOpen: input.Input.Is_always_open,
		HourFrom:     NullToString(input.Input.Hour_from),
		HourTo:       NullToString(input.Input.Hour_to),
		DeliveryTime: stringToServiceRPCDeliveryTimeEnum(input.Input.Delivery_time),
		Skill:        NullStringArrayToStringArray(input.Input.Skill),
		Price:        stringToServiceRPCPriceEnum(input.Input.Price),
		ServiceOwner: serviceOwnerRPCToServiceOwner(input.Input.Services_ownwer),
	}

	// Location Type
	if input.Input.Location_type != nil {
		query.LocationType = locationTypeToServiceRPCLocationTypeEnum(input.Input.Location_type)
	}

	// Week days
	if input.Input.Week_days != nil {
		query.WeekDays = weekDaysToRPC(*input.Input.Week_days)
	}

	result, err := search.ServiceSearch(ctx, &query)
	if err != nil {
		return SearchServiceResultResolver{}, err
	}

	services := make([]Service, 0, len(result.GetIDs()))
	for _, res := range result.GetIDs() {
		if res != "" {
			// Get Service by ID
			serv, err := r.GetVOfficeService(ctx, GetVOfficeServiceRequest{
				Service_id: res,
			})

			if err == nil {
				services = append(services, *serv.R)
			}
		}
	}

	return SearchServiceResultResolver{
		R: &SearchServiceResult{
			Amount_of_results:     int32(result.GetTotal()),
			Service_search_result: services,
		},
	}, nil

}

func (r *Resolver) SearchServiceRequests(ctx context.Context, input SearchServiceRequestsRequest) (SearchServiceRequestResultResolver, error) {
	var first uint32
	var after string

	if input.Pagination.First != nil {
		first = uint32(*input.Pagination.First)
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	query := searchRPC.ServiceRequestSearchRequest{
		After:        after,
		First:        first,
		Keywords:     NullStringArrayToStringArray(input.Input.Keywords),
		City:         NullStringArrayToStringArray(input.Input.City),
		Country:      NullStringArrayToStringArray(input.Input.Country),
		Currency:     NullToString(input.Input.Currency),
		FixedPrice:   NullToInt32(input.Input.Fixed_price_amount),
		MinSalary:    NullToInt32(input.Input.Min_salary),
		MaxSalary:    NullToInt32(input.Input.Max_salary),
		DeliveryTime: stringToServiceRPCDeliveryTimeEnum(input.Input.Delivery_time),
		ProjectType:  stringToServiceRPCRequestProjectTypes(input.Input.Project_type),
		Skills:       NullStringArrayToStringArray(input.Input.Skill),
		Languages:    NullStringArrayToStringArray(input.Input.Language),
		LocationType: locationTypeToServiceRPCLocationTypeEnum(input.Input.Location_type),
		Tools:        NullStringArrayToStringArray(input.Input.Tool),
		PriceType:    stringToServiceRPCPriceEnum(input.Input.Price),
		ServiceOwner: serviceOwnerRPCToServiceOwner(input.Input.Services_ownwer),
	}

	// Location Type
	if input.Input.Location_type != nil {
		query.LocationType = locationTypeToServiceRPCLocationTypeEnum(input.Input.Location_type)
	}

	result, err := search.ServiceRequestSearch(ctx, &query)
	if err != nil {
		return SearchServiceRequestResultResolver{}, err
	}

	requests := make([]ServiceRequest, 0, len(result.GetIDs()))
	for _, res := range result.GetIDs() {
		if res != "" {
			// Get Service Request by ID
			req, err := r.GetServiceRequest(ctx, GetServiceRequestRequest{
				Service_id: res,
			})

			if err == nil {
				requests = append(requests, *req.R)
			}
		}
	}

	return SearchServiceRequestResultResolver{
		R: &SearchServiceRequestResult{
			Amount_of_results:             int32(result.GetTotal()),
			Service_request_search_result: requests,
		},
	}, nil

}

func GetFiltersVariable(ctx context.Context, res *searchRPC.FilterArr) (*[]SearchFilterInterfaceResolver, error) {

	filters := make([]SearchFilterInterfaceResolver, 0, len(res.GetFilters()))

	for _, f := range res.GetFilters() {
		var filt SearchFilterInterface

		if f.GetUserSearch() != nil {

			filt = UserSearchFilterFragmentResolver{
				R: &UserSearchFilterFragment{
					ID:             f.GetID(),
					Filter_name:    f.GetName(),
					City:           searchRPCCityArrayToCityArray(f.GetUserSearch().GetCity()),
					City_id:        f.GetUserSearch().GetCityID(),
					ConenctionsOf:  f.GetUserSearch().GetConnectionsOfID(),
					Country:        f.GetUserSearch().GetCountryID(),
					CurrentCompany: f.GetUserSearch().GetCurrentCompany(),
					Degree:         f.GetUserSearch().GetDegree(),
					FiledOfStudy:   f.GetUserSearch().GetFieldOfStudy(),
					Firstname:      f.GetUserSearch().GetFirstname(),
					Industry:       f.GetUserSearch().GetIndustry(),
					Interest:       f.GetUserSearch().GetInterest(),
					IsFemale:       f.GetUserSearch().GetIsFemale(),
					IsMale:         f.GetUserSearch().GetIsMale(),
					IsMyConnection: f.GetUserSearch().GetMyConnections(),
					IsStudent:      f.GetUserSearch().GetIsStudent(),
					Keywords:       f.GetUserSearch().GetKeywords(),
					Language:       f.GetUserSearch().GetLanguage(),
					Nickname:       f.GetUserSearch().GetNickname(),
					PastCompany:    f.GetUserSearch().GetPastCompany(),
					Position:       f.GetUserSearch().GetPosition(),
					School:         f.GetUserSearch().GetSchool(),
					Skill:          f.GetUserSearch().GetSkill(),
					MinAge:         Uint32Toint32(f.GetUserSearch().GetMinAge()),
					MaxAge:         Uint32Toint32(f.GetUserSearch().GetMaxAge()),
					Full_name:      f.GetUserSearch().GetFullName(),
				},
			}
		} else if f.GetJobSearch() != nil {
			filt = SearchJobFilterFragmentResolver{
				R: &SearchJobFilterFragment{
					ID:                   f.GetID(),
					Filter_name:          f.GetName(),
					City:                 searchRPCCityArrayToCityArray(f.GetJobSearch().GetCity()),
					City_id:              f.GetCandidateSearch().GetCityID(),
					Keywords:             f.GetJobSearch().GetKeywords(),
					Degree:               f.GetJobSearch().GetDegree(),
					Country:              f.GetJobSearch().GetCountry(),
					Job_type:             f.GetJobSearch().GetJobType(),
					Language:             f.GetJobSearch().GetLanguage(),
					Industry:             f.GetJobSearch().GetIndustry(),
					Subindustry:          f.GetJobSearch().GetSubindustry(),
					Company_name:         f.GetJobSearch().GetCompanyName(),
					Currency:             f.GetJobSearch().GetCurrency(),
					Period:               f.GetJobSearch().GetPeriod(),
					Skill:                f.GetJobSearch().GetSkill(),
					Without_cover_letter: f.GetJobSearch().GetWithoutCoverLetter(),
					With_salary:          f.GetJobSearch().GetWithSalary(),
					Date_posted:          searchRPCDateEnumToString(f.GetJobSearch().GetDatePosted()),
					Experience_level:     searchExperienceEnumToString(f.GetJobSearch().GetExperienceLevel()),
					Company_size:         searchRPCToString(f.GetJobSearch().GetCompanySize()),
					// Birthday: // TODO:
				},
			}
		} else if f.GetCompanySearch() != nil {
			filt = SearchCompanyFilterFragmentResolver{
				R: &SearchCompanyFilterFragment{
					ID:                       f.GetID(),
					Filter_name:              f.GetName(),
					Search_for_companies:     f.GetCompanySearch().GetIsCompany(),
					Search_for_organizations: f.GetCompanySearch().GetIsOrganization(),
					With_jobs:                f.GetCompanySearch().GetIsJobOffers(),
					Keywords:                 f.GetCompanySearch().GetKeywords(),
					Size:                     companySizeRPCToString(f.GetCompanySearch().GetSize()),
					Type:                     companyTypeRPCToString(f.GetCompanySearch().GetType()),
					Name:                     f.GetCompanySearch().GetName(),
					Industry:                 f.GetCompanySearch().GetIndustry(),
					Subindustry:              f.GetCompanySearch().GetSubIndustry(),
					City:                     searchRPCCityArrayToCityArray(f.GetCompanySearch().GetCity()),
					City_id:                  f.GetCompanySearch().GetCityID(),
					Country:                  f.GetCompanySearch().GetCountry(),
					Founders_id:              f.GetCompanySearch().GetFounderIDs(),
					Founders_name:            f.GetCompanySearch().GetFounderNames(),
					Business_hours:           f.GetCompanySearch().GetBusinessHours(),
				},
			}

		} else if f.GetCandidateSearch() != nil {
			filt = SearchCandidateFilterFragmentResolver{
				R: &SearchCandidateFilterFragment{
					ID:                         f.GetID(),
					Filter_name:                f.GetName(),
					Keywords:                   f.GetCandidateSearch().GetKeywords(),
					Country:                    f.GetCandidateSearch().GetCountry(),
					City:                       searchRPCCityArrayToCityArray(f.GetCandidateSearch().GetCity()),
					City_id:                    f.GetCandidateSearch().GetCityID(),
					Current_company:            f.GetCandidateSearch().GetCurrentCompany(),
					Past_company:               f.GetCandidateSearch().GetPastCompany(),
					Industry:                   f.GetCandidateSearch().GetIndustry(),
					Sub_industry:               f.GetCandidateSearch().GetSubIndustry(),
					Job_type:                   f.GetCandidateSearch().GetJobType(),
					Skill:                      f.GetCandidateSearch().GetSkill(),
					Language:                   f.GetCandidateSearch().GetLanguage(),
					School:                     f.GetCandidateSearch().GetSchool(),
					Degree:                     f.GetCandidateSearch().GetDegree(),
					Field_of_study:             f.GetCandidateSearch().GetFieldOfStudy(),
					Is_student:                 f.GetCandidateSearch().GetIsStudent(),
					Currency:                   f.GetCandidateSearch().GetCurrency(),
					Period:                     f.GetCandidateSearch().GetPeriod(),
					Is_willing_to_travel:       f.GetCandidateSearch().GetIsWillingToTravel(),
					Is_willing_to_work_remotly: f.GetCandidateSearch().GetIsWillingToWorkRemotly(),
					Is_possible_to_relocate:    f.GetCandidateSearch().GetIsPossibleToRelocate(),
					Experience_level:           searchExperienceEnumToString(f.GetCandidateSearch().GetExperienceLevel()),
					// },
				},
			}
		} else if f.GetServiceSearch() != nil {
			service := f.GetServiceSearch()
			filt = SearchServiceFilterFragmentResolver{
				R: &SearchServiceFilterFragment{
					ID:                 f.GetID(),
					Filter_name:        f.GetName(),
					Keywords:           service.GetKeywords(),
					City:               service.GetCity(),
					Country:            service.GetCountry(),
					Is_always_open:     service.GetIsAlwaysOpen(),
					Currency:           service.GetCurrency(),
					Fixed_price_amount: service.GetFixedPrice(),
					Min_price_amount:   service.GetMinSalary(),
					Max_price_amount:   service.GetMaxSalary(),
					Hour_to:            service.GetHourTo(),
					Hour_from:          service.GetHourFrom(),
					Skill:              service.GetSkill(),
					Delivery_time:      servicesRPCDeliveryTimeToDeliveryTime(service.GetDeliveryTime()),
					Location_type:      servicesRPCLocationTypeToLocationType(service.GetLocationType()),
					Price:              servicesRPCPriceToPrice(service.GetPrice()),
					Services_ownwer:    serviceOwnerRPCToString(service.GetServiceOwner()),
					Week_days:          serviceWeekDaysRPCToWeekDays(service.GetWeekDays()),
				},
			}
		} else if f.GetServiceRequestSearch() != nil {
			service := f.GetServiceRequestSearch()
			filt = SearchServiceRequestFilterFragmentResolver{
				R: &SearchServiceRequestFilterFragment{
					ID:                 f.GetID(),
					Filter_name:        f.GetName(),
					Keywords:           service.GetKeywords(),
					City:               service.GetCity(),
					Country:            service.GetCountry(),
					Currency:           service.GetCurrency(),
					Fixed_price_amount: service.GetFixedPrice(),
					Min_price_amount:   service.GetMinSalary(),
					Max_price_amount:   service.GetMaxSalary(),
					Skill:              service.GetSkills(),
					Language:           service.GetLanguages(),
					Tool:               service.GetTools(),
					Project_type:       projectTypesToRPC(service.GetProjectType()),
					Delivery_time:      servicesRPCDeliveryTimeToDeliveryTime(service.GetDeliveryTime()),
					Location_type:      servicesRPCLocationTypeToLocationType(service.GetLocationType()),
					Price:              servicesRPCPriceToPrice(service.GetPriceType()),
					Services_ownwer:    serviceOwnerRPCToString(service.GetServiceOwner()),
				},
			}
		}

		filters = append(filters, SearchFilterInterfaceResolver{r: &filt})
	}

	return &filters, nil
}

// GetAllFilters ...
func (_ *Resolver) GetAllFilters(ctx context.Context) (*[]SearchFilterInterfaceResolver, error) {
	res, err := search.GetAllFilters(ctx, &searchRPC.Empty{})
	if err != nil {
		return nil, err
	}

	filters, err := GetFiltersVariable(ctx, res)
	if err != nil {
		return nil, err
	}

	return filters, nil
}

// GetAllFilters ...
func (_ *Resolver) GetAllFiltersForCompany(ctx context.Context, input GetAllLabelsForCompanyRequest) (*[]SearchFilterInterfaceResolver, error) {
	res, err := search.GetAllFiltersForCompany(ctx, &searchRPC.ID{
		ID: input.CompanyId,
	})
	if err != nil {
		return nil, err
	}

	filters, err := GetFiltersVariable(ctx, res)
	if err != nil {
		return nil, err
	}

	return filters, nil
}

// GetlistOfFiltersByTypeForCompany ...
// func (_ *Resolver) GetListOfFiltersByTypeForCompany(ctx context.Context, input GetListOfFiltersByTypeForCompanyRequest) (*[]SearchFilterInterfaceResolver, error) {
// 	res, err := search.GetFiltersByTypeForCompany(ctx, &searchRPC.FilterTypeRequest{
// 		Type:      stringToFilterToTypeRequestRPC(input.Filter_type),
// 		CompanyID: input.Company_id,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	filters, err := GetFiltersVariable(ctx, res)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return filters, nil
// }

// // RemoveFilter ...
func (_ *Resolver) RemoveFilter(ctx context.Context, input RemoveFilterRequest) (*SuccessResolver, error) {
	_, err := search.RemoveFilter(ctx, &searchRPC.ID{
		ID: input.Filter_id,
	},
	)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
			R: &Success{
				Success: true,
			},
		},
		nil
}

// // RemoveFilterForCompany ...
func (_ *Resolver) RemoveFilterForCompany(ctx context.Context, input RemoveFilterForCompanyRequest) (*SuccessResolver, error) {
	_, err := search.RemoveFilterForCompany(ctx, &searchRPC.IDs{
		ID:        input.Filter_id,
		CompanyID: input.CompanyID,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
			R: &Success{
				Success: true,
			},
		},
		nil
}

func stringToFilterToTypeRequestRPC(s string) searchRPC.FilterTypeRequest_FilterType {
	switch s {
	case "SearchCandidateFilterType":
		return searchRPC.FilterTypeRequest_CandidateFilterType
	case "SearchJobFilterType":
		return searchRPC.FilterTypeRequest_JobFilterType
	case "SearchCompanyFilterType":
		return searchRPC.FilterTypeRequest_CompanyFilterType
	}

	return searchRPC.FilterTypeRequest_UserFilterType
}

func filterTypeRequestRPCToString(s searchRPC.FilterTypeRequest_FilterType) string {
	switch s {
	case searchRPC.FilterTypeRequest_CandidateFilterType:
		return "SearchCandidateFilterType"
	case searchRPC.FilterTypeRequest_JobFilterType:
		return "SearchJobFilterType"
	case searchRPC.FilterTypeRequest_CompanyFilterType:
		return "SearchCompanyFilterType"
	}

	return "SearchUserFilterType"
}

// SaveUserSearchFilter ...
func (_ *Resolver) SaveUserSearchFilter(ctx context.Context, input SaveUserSearchFilterRequest) (SuccessResolver, error) {
	res, err := search.SaveUserSearchFilter(ctx, &searchRPC.SaveUserSearchFilterRequest{
		Name: input.User_filter.Name,
		UserSearch: &searchRPC.UserSearchRequest{
			Keywords:        NullStringArrayToStringArray(input.User_filter.Filter.Keywords),
			MyConnections:   input.User_filter.Filter.IsMyConnection,
			ConnectionsOfID: NullStringArrayToStringArray(input.User_filter.Filter.ConenctionsOf),
			CountryID:       NullStringArrayToStringArray(input.User_filter.Filter.Country),
			CityID:          NullStringArrayToStringArray(input.User_filter.Filter.City),
			School:          NullStringArrayToStringArray(input.User_filter.Filter.School),
			Degree:          NullStringArrayToStringArray(input.User_filter.Filter.Degree),
			FieldOfStudy:    NullStringArrayToStringArray(input.User_filter.Filter.FiledOfStudy),
			IsStudent:       input.User_filter.Filter.IsStudent,
			CurrentCompany:  NullStringArrayToStringArray(input.User_filter.Filter.CurrentCompany),
			PastCompany:     NullStringArrayToStringArray(input.User_filter.Filter.PastCompany),
			Industry:        NullStringArrayToStringArray(input.User_filter.Filter.Industry),
			Position:        NullStringArrayToStringArray(input.User_filter.Filter.Position),
			Firstname:       NullStringArrayToStringArray(input.User_filter.Filter.Firstname),
			Lastname:        NullStringArrayToStringArray(input.User_filter.Filter.Lastname),
			Nickname:        NullStringArrayToStringArray(input.User_filter.Filter.Nickname),
			IsMale:          input.User_filter.Filter.IsMale,
			IsFemale:        input.User_filter.Filter.IsFemale,
			Skill:           NullStringArrayToStringArray(input.User_filter.Filter.Skill),
			Language:        NullStringArrayToStringArray(input.User_filter.Filter.Language),
			Interest:        NullStringArrayToStringArray(input.User_filter.Filter.Interest),
			MaxAge:          Nullint32ToUint32(input.User_filter.Filter.MaxAge),
			MinAge:          Nullint32ToUint32(input.User_filter.Filter.MinAge),
			FullName:        NullToString(input.User_filter.Filter.Full_name),
		},
	})
	if err != nil {
		return SuccessResolver{}, err
	}

	return SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// SaveJobSearchFilter ...
func (_ *Resolver) SaveJobSearchFilter(ctx context.Context, input SaveJobSearchFilterRequest) (SuccessResolver, error) {
	res, err := search.SaveJobSearchFilter(ctx, &searchRPC.SaveJobSearchFilterRequest{
		Name: input.Job_filter.Name,
		JobSearch: &searchRPC.JobSearchRequest{
			Keywords:           NullStringArrayToStringArray(input.Job_filter.Filter.Keywords),
			Degree:             NullStringArrayToStringArray(input.Job_filter.Filter.Degree),
			Country:            NullStringArrayToStringArray(input.Job_filter.Filter.Country),
			City:               NullStringArrayToStringArray(input.Job_filter.Filter.City),
			JobType:            NullStringArrayToStringArray(input.Job_filter.Filter.Job_type),
			Language:           NullStringArrayToStringArray(input.Job_filter.Filter.Language),
			Industry:           NullStringArrayToStringArray(input.Job_filter.Filter.Industry),
			Subindustry:        NullStringArrayToStringArray(input.Job_filter.Filter.Subindustry),
			CompanyName:        NullStringArrayToStringArray(input.Job_filter.Filter.Company_name),
			CompanySize:        stringToSearchRPC(input.Job_filter.Filter.Company_size),
			DatePosted:         stringToDateEnumRPC(input.Job_filter.Filter.Date_posted),
			Currency:           NullToString(input.Job_filter.Filter.Currency),
			Period:             NullToString(input.Job_filter.Filter.Period),
			Skill:              NullStringArrayToStringArray(input.Job_filter.Filter.Skill),
			IsFollowing:        input.Job_filter.Filter.Is_following,
			WithoutCoverLetter: input.Job_filter.Filter.Without_cover_letter,
			WithSalary:         input.Job_filter.Filter.With_salary,
			ExperienceLevel:    stringTotExperienceEnumRPC(input.Job_filter.Filter.Experience_level),
		},
	},
	)

	if err != nil {
		return SuccessResolver{}, err
	}

	return SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// SaveServiceSearchFilter ...
func (_ *Resolver) SaveServiceSearchFilter(ctx context.Context, input SaveServiceSearchFilterRequest) (SuccessResolver, error) {
	res, err := search.SaveServiceSearchFilter(ctx, &searchRPC.SaveServiceSearchFilterRequest{
		Name:      NullToString(input.Service_filter.Name),
		CompanyID: NullToString(input.Company_id),
		ServiceSearch: &searchRPC.ServiceSearchRequest{
			Keywords:     NullStringArrayToStringArray(input.Service_filter.Keywords),
			City:         NullStringArrayToStringArray(input.Service_filter.City),
			Country:      NullStringArrayToStringArray(input.Service_filter.Country),
			Currency:     NullToString(input.Service_filter.Currency),
			FixedPrice:   NullToInt32(input.Service_filter.Fixed_price_amount),
			MinSalary:    NullToInt32(input.Service_filter.Min_salary),
			MaxSalary:    NullToInt32(input.Service_filter.Max_salary),
			LocationType: locationTypeToServiceRPCLocationTypeEnum(input.Service_filter.Location_type),
			IsAlwaysOpen: input.Service_filter.Is_always_open,
			HourFrom:     NullToString(input.Service_filter.Hour_from),
			HourTo:       NullToString(input.Service_filter.Hour_to),
			DeliveryTime: stringToServiceRPCDeliveryTimeEnum(input.Service_filter.Delivery_time),
			Skill:        NullStringArrayToStringArray(input.Service_filter.Skill),
			Price:        stringToServiceRPCPriceEnum(input.Service_filter.Price),
			ServiceOwner: serviceOwnerRPCToServiceOwner(input.Service_filter.Services_ownwer),
		},
	},
	)

	if err != nil {
		return SuccessResolver{}, err
	}

	return SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// SaveServiceRequestSearchFilter ...
func (_ *Resolver) SaveServiceRequestSearchFilter(ctx context.Context, input SaveServiceRequestSearchFilterRequest) (SuccessResolver, error) {
	res, err := search.SaveServiceRequestSearchFilter(ctx, &searchRPC.SaveServiceRequestSearchFilterRequest{
		Name:      NullToString(input.Service_request_filter.Name),
		CompanyID: NullToString(input.Company_id),
		ServiceRequestSearch: &searchRPC.ServiceRequestSearchRequest{
			Keywords:     NullStringArrayToStringArray(input.Service_request_filter.Keywords),
			City:         NullStringArrayToStringArray(input.Service_request_filter.City),
			Country:      NullStringArrayToStringArray(input.Service_request_filter.Country),
			Currency:     NullToString(input.Service_request_filter.Currency),
			FixedPrice:   NullToInt32(input.Service_request_filter.Fixed_price_amount),
			MinSalary:    NullToInt32(input.Service_request_filter.Min_salary),
			MaxSalary:    NullToInt32(input.Service_request_filter.Max_salary),
			LocationType: locationTypeToServiceRPCLocationTypeEnum(input.Service_request_filter.Location_type),
			DeliveryTime: stringToServiceRPCDeliveryTimeEnum(input.Service_request_filter.Delivery_time),
			Skills:       NullStringArrayToStringArray(input.Service_request_filter.Skill),
			Languages:    NullStringArrayToStringArray(input.Service_request_filter.Language),
			Tools:        NullStringArrayToStringArray(input.Service_request_filter.Tool),
			PriceType:    stringToServiceRPCPriceEnum(input.Service_request_filter.Price),
			ProjectType:  stringToServiceRPCRequestProjectTypes(input.Service_request_filter.Project_type),
			ServiceOwner: serviceOwnerRPCToServiceOwner(input.Service_request_filter.Services_ownwer),
		},
	},
	)

	if err != nil {
		return SuccessResolver{}, err
	}

	return SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// SaveCompanySearchFilter ...
func (_ *Resolver) SaveCompanySearchFilter(ctx context.Context, input SaveCompanySearchFilterRequest) (SuccessResolver, error) {
	res, err := search.SaveCompanySearchFilter(ctx, &searchRPC.SaveCompanySearchFilterRequest{
		Name: input.Company_filter.Name,
		CompanySearch: &searchRPC.CompanySearchRequest{
			IsCompany:      input.Company_filter.Filter.Search_for_companies,
			IsJobOffers:    input.Company_filter.Filter.With_jobs,
			Keywords:       NullStringArrayToStringArray(input.Company_filter.Filter.Keywords),
			Size:           stringToCompanySizeRPC(input.Company_filter.Filter.Size),
			Type:           stringToCompanyTypeRPC(input.Company_filter.Filter.Type),
			IsOrganization: input.Company_filter.Filter.Search_for_organizations,
			Name:           NullStringArrayToStringArray(input.Company_filter.Filter.Name),
			Industry:       NullStringArrayToStringArray(input.Company_filter.Filter.Industry),
			SubIndustry:    NullStringArrayToStringArray(input.Company_filter.Filter.Subindustry),
			City:           NullStringArrayToStringArray(input.Company_filter.Filter.City),
			Country:        NullStringArrayToStringArray(input.Company_filter.Filter.Country),
			FounderIDs:     NullStringArrayToStringArray(input.Company_filter.Filter.Founders_id),
			FounderNames:   NullStringArrayToStringArray(input.Company_filter.Filter.Founders_name),
			BusinessHours:  NullStringArrayToStringArray(input.Company_filter.Filter.Business_hours),
			// Rating:         input.Input.Rating,
		}})
	if err != nil {
		return SuccessResolver{}, err
	}

	return SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// // SaveUserSearchFilterForCompany ...
func (_ *Resolver) SaveUserSearchFilterForCompany(ctx context.Context, input SaveUserSearchFilterForCompanyRequest) (*SuccessResolver, error) {
	res, err := search.SaveUserSearchFilterForCompany(ctx, &searchRPC.SaveUserSearchFilterRequest{
		Name:      input.User_filter.Name,
		CompanyID: input.User_filter.CompanyID,
		UserSearch: &searchRPC.UserSearchRequest{
			Keywords:        NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Keywords),
			MyConnections:   input.User_filter.UserSearchFilter.IsMyConnection,
			ConnectionsOfID: NullStringArrayToStringArray(input.User_filter.UserSearchFilter.ConenctionsOf),
			CountryID:       NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Country),
			CityID:          NullStringArrayToStringArray(input.User_filter.UserSearchFilter.City),
			School:          NullStringArrayToStringArray(input.User_filter.UserSearchFilter.School),
			Degree:          NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Degree),
			FieldOfStudy:    NullStringArrayToStringArray(input.User_filter.UserSearchFilter.FiledOfStudy),
			IsStudent:       input.User_filter.UserSearchFilter.IsStudent,
			CurrentCompany:  NullStringArrayToStringArray(input.User_filter.UserSearchFilter.CurrentCompany),
			PastCompany:     NullStringArrayToStringArray(input.User_filter.UserSearchFilter.PastCompany),
			Industry:        NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Industry),
			Position:        NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Position),
			Firstname:       NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Firstname),
			Lastname:        NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Lastname),
			Nickname:        NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Nickname),
			IsMale:          input.User_filter.UserSearchFilter.IsMale,
			IsFemale:        input.User_filter.UserSearchFilter.IsFemale,
			Skill:           NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Skill),
			Language:        NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Language),
			Interest:        NullStringArrayToStringArray(input.User_filter.UserSearchFilter.Interest),
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// // SaveCompanySearchFilterForCompany ...
func (_ *Resolver) SaveCompanySearchFilterForCompany(ctx context.Context, input SaveCompanySearchFilterForCompanyRequest) (*SuccessResolver, error) {
	res, err := search.SaveCompanySearchFilterForCompany(ctx, &searchRPC.SaveCompanySearchFilterRequest{
		CompanyID: input.Company_filter.CompanyID,
		Name:      input.Company_filter.Name,
		CompanySearch: &searchRPC.CompanySearchRequest{
			IsCompany:      input.Company_filter.CompanySearchFilter.Search_for_companies,
			IsJobOffers:    input.Company_filter.CompanySearchFilter.With_jobs,
			Keywords:       NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.Keywords),
			Size:           stringToCompanySizeRPC(input.Company_filter.CompanySearchFilter.Size),
			Type:           stringToCompanyTypeRPC(input.Company_filter.CompanySearchFilter.Type),
			IsOrganization: input.Company_filter.CompanySearchFilter.Search_for_organizations,
			Name:           NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.Name),
			Industry:       NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.Industry),
			SubIndustry:    NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.Subindustry),
			City:           NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.City),
			Country:        NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.Country),
			FounderIDs:     NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.Founders_id),
			FounderNames:   NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.Founders_name),
			BusinessHours:  NullStringArrayToStringArray(input.Company_filter.CompanySearchFilter.Business_hours),
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// SaveJobSearchFilterForCompany ...
func (_ *Resolver) SaveJobSearchFilterForCompany(ctx context.Context, input SaveJobSearchFilterForCompanyRequest) (*SuccessResolver, error) {
	res, err := search.SaveJobSearchFilterForCompany(ctx, &searchRPC.SaveJobSearchFilterRequest{
		CompanyID: input.Job_filter.CompanyID,
		Name:      input.Job_filter.Name,
		JobSearch: &searchRPC.JobSearchRequest{
			Keywords:           NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Keywords),
			Degree:             NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Degree),
			Country:            NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Country),
			City:               NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.City),
			JobType:            NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Job_type),
			Language:           NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Language),
			Industry:           NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Industry),
			Subindustry:        NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Subindustry),
			CompanyName:        NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Company_name),
			CompanySize:        stringToSearchRPC(input.Job_filter.JobSearchFilter.Company_size),
			Currency:           NullToString(input.Job_filter.JobSearchFilter.Currency),
			Period:             NullToString(input.Job_filter.JobSearchFilter.Period),
			Skill:              NullStringArrayToStringArray(input.Job_filter.JobSearchFilter.Skill),
			IsFollowing:        input.Job_filter.JobSearchFilter.Is_following,
			WithoutCoverLetter: input.Job_filter.JobSearchFilter.Without_cover_letter,
			WithSalary:         input.Job_filter.JobSearchFilter.With_salary,
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// // SaveCandidateSearchFilterForCompany ...
func (_ *Resolver) SaveCandidateSearchFilterForCompany(ctx context.Context, input SaveCandidateSearchFilterForCompanyRequest) (*SuccessResolver, error) {
	res, err := search.SaveCandidateSearchFilterForCompany(ctx, &searchRPC.SaveCandidateSearchFilterRequest{
		CompanyID: input.Candidate_filter.CompanyID,
		Name:      input.Candidate_filter.Name,
		CandidateSearch: &searchRPC.CandidateSearchRequest{
			Keywords:               NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Keywords),
			Country:                NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Country),
			City:                   NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.City),
			CurrentCompany:         NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Current_company),
			PastCompany:            NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Past_company),
			Industry:               NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Industry),
			SubIndustry:            NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Sub_industry),
			JobType:                NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Job_type),
			Skill:                  NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Skill),
			Language:               NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Language),
			School:                 NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.School),
			Degree:                 NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Degree),
			FieldOfStudy:           NullStringArrayToStringArray(input.Candidate_filter.CandidateSearchFilter.Field_of_study),
			IsStudent:              input.Candidate_filter.CandidateSearchFilter.Is_student,
			Currency:               NullToString(input.Candidate_filter.CandidateSearchFilter.Currency),
			Period:                 input.Candidate_filter.CandidateSearchFilter.Period,
			IsWillingToTravel:      input.Candidate_filter.CandidateSearchFilter.Is_willing_to_travel,
			IsWillingToWorkRemotly: input.Candidate_filter.CandidateSearchFilter.Is_willing_to_work_remotly,
			IsPossibleToRelocate:   input.Candidate_filter.CandidateSearchFilter.Is_possible_to_relocate,
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// wan't generated automatically

func (r SearchFilterInterfaceResolver) ID() graphql.ID {
	return (*r.r).ID()
}

func (r SearchFilterInterfaceResolver) Filter_name() string {
	return (*r.r).Filter_name()
}

func (r SearchFilterInterfaceResolver) ToUserSearchFilterFragment() (*UserSearchFilterFragmentResolver, bool) {
	res, ok := (*r.r).(UserSearchFilterFragmentResolver)
	return &res, ok
}

func (r SearchFilterInterfaceResolver) ToSearchJobFilterFragment() (*SearchJobFilterFragmentResolver, bool) {
	res, ok := (*r.r).(SearchJobFilterFragmentResolver)
	return &res, ok
}

func (r SearchFilterInterfaceResolver) ToSearchCompanyFilterFragment() (*SearchCompanyFilterFragmentResolver, bool) {
	res, ok := (*r.r).(SearchCompanyFilterFragmentResolver)
	return &res, ok
}

func (r SearchFilterInterfaceResolver) ToSearchCandidateFilterFragment() (*SearchCandidateFilterFragmentResolver, bool) {
	res, ok := (*r.r).(SearchCandidateFilterFragmentResolver)
	return &res, ok
}

func (r SearchFilterInterfaceResolver) ToSearchServiceFilterFragment() (*SearchServiceFilterFragmentResolver, bool) {
	res, ok := (*r.r).(SearchServiceFilterFragmentResolver)
	return &res, ok
}

func (r SearchFilterInterfaceResolver) ToSearchServiceRequestFilterFragment() (*SearchServiceRequestFilterFragmentResolver, bool) {
	res, ok := (*r.r).(SearchServiceRequestFilterFragmentResolver)
	return &res, ok
}
