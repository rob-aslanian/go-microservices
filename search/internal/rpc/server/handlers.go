package serverRPC

import (
	"context"
	"errors"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/searchRPC"
	"gitlab.lan/Rightnao-site/microservices/search/internal/company"
	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
)

// UserSearch ...
func (s Server) UserSearch(ctx context.Context, data *searchRPC.UserSearchRequest) (*searchRPC.Profiles, error) {
	profilesInterface, total, err := s.service.UserSearch(
		ctx,
		userSearchRequestRPCToRequestsUserSearch(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.Profiles{
		Profiles: profilesInterface,
		Total:    total,
	}, nil
}

// CompanySearch ...
func (s Server) CompanySearch(ctx context.Context, data *searchRPC.CompanySearchRequest) (*searchRPC.Companies, error) {
	profilesInterface, total, err := s.service.CompanySearch(
		ctx,
		companySearchRequestRPCToRequestsCompanySearch(data),
	)
	if err != nil {
		return nil, err
	}

	profiles := make([]*companyRPC.Profile, 0)

	if prof, ok := profilesInterface.([]*companyRPC.Profile); ok {
		for _, p := range prof {
			profiles = append(profiles, p)
		}
	} else {
		log.Printf("wrong type assertion. Want []*companyRPC.Profile. Got %T", profilesInterface)
		return nil, errors.New("internal_error")
	}

	return &searchRPC.Companies{
		Results: profiles,
		Total:   total,
	}, nil
}

// JobSearch ...
func (s Server) JobSearch(ctx context.Context, data *searchRPC.JobSearchRequest) (*searchRPC.Jobs, error) {

	jobsInterface, total, err := s.service.JobSearch(
		ctx,
		jobSearchRPCToRequestsJobSearch(data),
	)
	if err != nil {
		return nil, err
	}

	profiles := make([]*searchRPC.JobResult, 0)

	if prof, ok := jobsInterface.([]*searchRPC.JobResult); ok {
		for _, p := range prof {
			profiles = append(profiles, p)
		}
	} else {
		log.Printf("wrong type assertion. Want []*searchRPC.JobResult. Got %T", jobsInterface)
		return nil, errors.New("internal_error")
	}

	return &searchRPC.Jobs{
		JobResults: profiles,
		Total:      total,
	}, nil
}

// CandidateSearch ...
func (s Server) CandidateSearch(ctx context.Context, data *searchRPC.CandidateSearchRequest) (*searchRPC.Candidates, error) {
	profilesInterface, total, err := s.service.CandidateSearch(
		ctx,
		data.GetCompanyID(),
		candidateSearchRequestToRequestsCandidateSearch(data),
	)
	if err != nil {
		return nil, err
	}

	profiles := make([]*searchRPC.CandidateResult, 0)

	if prof, ok := profilesInterface.([]*searchRPC.CandidateResult); ok {
		for _, p := range prof {
			profiles = append(profiles, p)
		}
	} else {
		log.Printf("wrong type assertion. Want []*searchRPC.CandidateResult. Got %T", profilesInterface)
		return nil, errors.New("internal_error")
	}

	return &searchRPC.Candidates{
		CandidateResults: profiles,
		Total:            total,
	}, nil
}

// ServiceSearch ...
func (s Server) ServiceSearch(ctx context.Context, data *searchRPC.ServiceSearchRequest) (*searchRPC.Services, error) {
	ids, total, err := s.service.ServiceSearch(
		ctx,
		serviceSearchRequestToRequestsServiceSearch(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.Services{
		IDs:   ids,
		Total: total,
	}, nil
}

// ServiceRequestSearch ...
func (s Server) ServiceRequestSearch(ctx context.Context, data *searchRPC.ServiceRequestSearchRequest) (*searchRPC.ServiceRequests, error) {
	ids, total, err := s.service.ServiceRequestSearch(
		ctx,
		serviceRequestRPCToServiceRequest(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ServiceRequests{
		IDs:   ids,
		Total: total,
	}, nil
}

// SaveUserSearchFilter ...
func (s Server) SaveUserSearchFilter(ctx context.Context, data *searchRPC.SaveUserSearchFilterRequest) (*searchRPC.ID, error) {
	saveUserSearchFilter, err := s.service.SaveUserSearchFilter(
		ctx,
		saveUserSearchFilterRequestRPCToRequestsSaveUserSearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: saveUserSearchFilter,
	}, nil
}

// SaveJobSearchFilter ...
func (s Server) SaveJobSearchFilter(ctx context.Context, data *searchRPC.SaveJobSearchFilterRequest) (*searchRPC.ID, error) {
	saveJobSearchFilter, err := s.service.SaveJobSearchFilter(
		ctx,
		saveJobSearchFilterRequestRPCToRequestsSaveJobSearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: saveJobSearchFilter,
	}, nil
}

// SaveCompanySearchFilter ...
func (s Server) SaveCompanySearchFilter(ctx context.Context, data *searchRPC.SaveCompanySearchFilterRequest) (*searchRPC.ID, error) {
	saveCompanySearchFilter, err := s.service.SaveCompanySearchFilter(
		ctx,
		SaveCompanySearchFilterRequestRPCToSaveCompanySearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: saveCompanySearchFilter,
	}, nil
}

// SaveServiceSearchFilter ...
func (s Server) SaveServiceSearchFilter(ctx context.Context, data *searchRPC.SaveServiceSearchFilterRequest) (*searchRPC.ID, error) {
	id, err := s.service.SaveServiceSearchFilter(
		ctx,
		SaveServiceSearchFilterRequestRPCToSaveServiceSearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: id,
	}, nil
}

// SaveServiceRequestSearchFilter ...
func (s Server) SaveServiceRequestSearchFilter(ctx context.Context, data *searchRPC.SaveServiceRequestSearchFilterRequest) (*searchRPC.ID, error) {
	id, err := s.service.SaveServiceRequestSearchFilter(
		ctx,
		SaveServiceRequestSearchFilterRequestRPCToSaveServiceRequestSearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: id,
	}, nil
}

// RemoveFilter ...
func (s Server) RemoveFilter(ctx context.Context, data *searchRPC.ID) (*searchRPC.Empty, error) {
	err := s.service.RemoveFilter(
		ctx,
		data.GetID(),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.Empty{}, nil
}

// SaveUserSearchFilterForCompany ...
func (s Server) SaveUserSearchFilterForCompany(ctx context.Context, data *searchRPC.SaveUserSearchFilterRequest) (*searchRPC.ID, error) {
	saveUserSearchFilter, err := s.service.SaveUserSearchFilterForCompany(
		ctx,
		data.GetCompanyID(),
		saveUserSearchFilterRequestRPCToRequestsSaveUserSearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: saveUserSearchFilter,
	}, nil
}

// SaveJobSearchFilterForCompany ...
func (s Server) SaveJobSearchFilterForCompany(ctx context.Context, data *searchRPC.SaveJobSearchFilterRequest) (*searchRPC.ID, error) {
	saveJobSearchFilter, err := s.service.SaveJobSearchFilterForCompany(
		ctx,
		data.GetCompanyID(),
		saveJobSearchFilterRequestRPCToRequestsSaveJobSearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: saveJobSearchFilter,
	}, nil
}

// SaveCompanySearchFilterForCompany ...
func (s Server) SaveCompanySearchFilterForCompany(ctx context.Context, data *searchRPC.SaveCompanySearchFilterRequest) (*searchRPC.ID, error) {
	saveCompanySearchFilter, err := s.service.SaveCompanySearchFilterForCompany(
		ctx,
		data.GetCompanyID(),
		SaveCompanySearchFilterRequestRPCToSaveCompanySearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: saveCompanySearchFilter,
	}, nil
}

// SaveCandidateSearchFilterForCompany ...
func (s Server) SaveCandidateSearchFilterForCompany(ctx context.Context, data *searchRPC.SaveCandidateSearchFilterRequest) (*searchRPC.ID, error) {
	id, err := s.service.SaveCandidateSearchFilterForCompany(
		ctx,
		data.GetCompanyID(),
		SaveCandidateSearchFilterRequestRPCToSaveCandidateSearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.ID{
		ID: id,
	}, nil
}

// GetAllFiltersForCompany ...
func (s Server) GetAllFiltersForCompany(ctx context.Context, data *searchRPC.ID) (*searchRPC.FilterArr, error) {
	filtType, err := s.service.GetAllFiltersForCompany(
		ctx,
		data.GetID(),
	)
	if err != nil {
		return nil, err
	}

	f := make([]*searchRPC.Filter, 0)

	for _, fl := range filtType {

		filt := searchRPC.Filter{}

		switch flt := fl.(type) {
		case requests.UserSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_UserSearch{
				UserSearch: requestsUserSearchToUserSearchRequestRPC(&flt.UserSearch),
			}
			f = append(f, &filt)
		case requests.JobSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_JobSearch{
				JobSearch: requestsJobSearchToJobSearchRequestRPC(&flt.JobSearch),
			}
			f = append(f, &filt)
		case requests.CandidateSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_CandidateSearch{
				CandidateSearch: requestsCandidateSearchToCandidateSearchRequestRPC(&flt.CandidateSearch),
			}
			f = append(f, &filt)
		case requests.CompanySearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_CompanySearch{
				CompanySearch: requestsCompanySearchToCompanySearchRequestRPC(&flt.CompanySearch),
			}
			f = append(f, &filt)
		case requests.ServiceSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_ServiceSearch{
				ServiceSearch: requestsServiceSearchToServiceSearchRequestRPC(&flt.ServiceSearch),
			}
			f = append(f, &filt)
		case requests.ServiceRequestSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_ServiceRequestSearch{
				ServiceRequestSearch: requestsServiceRequestSearchToServiceRequestSearchRequestRPC(&flt.ServiceRequest),
			}
			f = append(f, &filt)
		}
	}

	return &searchRPC.FilterArr{
		Filters: f,
	}, nil
}

// RemoveFilterForCompany ...
func (s Server) RemoveFilterForCompany(ctx context.Context, data *searchRPC.IDs) (*searchRPC.Empty, error) {
	err := s.service.RemoveFilterForCompany(
		ctx,
		data.GetID(),
		data.GetCompanyID(),
	)
	if err != nil {
		return nil, err
	}

	return &searchRPC.Empty{}, nil
}

// GetAllFilters ...
func (s Server) GetAllFilters(ctx context.Context, data *searchRPC.Empty) (*searchRPC.FilterArr, error) {
	filters, err := s.service.GetAllFilters(
		ctx,
	)
	if err != nil {
		return nil, err
	}
	f := make([]*searchRPC.Filter, 0)

	for _, fl := range filters {

		filt := searchRPC.Filter{}

		switch flt := fl.(type) {
		case requests.UserSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_UserSearch{
				UserSearch: requestsUserSearchToUserSearchRequestRPC(&flt.UserSearch),
			}
			f = append(f, &filt)
		case requests.JobSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_JobSearch{
				JobSearch: requestsJobSearchToJobSearchRequestRPC(&flt.JobSearch),
			}
			f = append(f, &filt)
		case requests.CandidateSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_CandidateSearch{
				CandidateSearch: requestsCandidateSearchToCandidateSearchRequestRPC(&flt.CandidateSearch),
			}
			f = append(f, &filt)
		case requests.CompanySearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_CompanySearch{
				CompanySearch: requestsCompanySearchToCompanySearchRequestRPC(&flt.CompanySearch),
			}
			f = append(f, &filt)

		case requests.ServiceSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_ServiceSearch{
				ServiceSearch: requestsServiceSearchToServiceSearchRequestRPC(&flt.ServiceSearch),
			}
			f = append(f, &filt)

		case requests.ServiceRequestSearchFilter:
			filt.ID = flt.GetID()
			filt.Name = flt.Name
			filt.Filters = &searchRPC.Filter_ServiceRequestSearch{
				ServiceRequestSearch: requestsServiceRequestSearchToServiceRequestSearchRequestRPC(&flt.ServiceRequest),
			}
			f = append(f, &filt)

		}
	}

	return &searchRPC.FilterArr{
		Filters: f,
	}, nil
}

// to RPC
// TODO create functions each for jobs, company and candidate to transform request to searchRPC.
// TODo finish getting the filters in switch case (line 199-209) GetFiltersByType()

// requestsUserSearchToUserSearchRequestRPC ...
func requestsUserSearchToUserSearchRequestRPC(data *requests.UserSearch) *searchRPC.UserSearchResult {
	if data == nil {
		return nil
	}

	filter := searchRPC.UserSearchResult{
		First:           data.First,
		ConnectionsOfID: data.ConnectionsOfID,
		CurrentCompany:  data.CurrentCompany,
		Nickname:        data.Nickname,
		Skill:           data.Skill,
		Interest:        data.Interest,
		MyConnections:   data.MyConnections,
		CountryID:       data.CountryID,
		School:          data.School,
		Position:        data.Position,
		Firstname:       data.Firstname,
		Lastname:        data.Lastname,
		After:           data.After,
		Keywords:        data.Keyword,
		Degree:          data.Degree,
		FieldOfStudy:    data.FieldOfStudy,
		IsStudent:       data.IsStudent,
		Industry:        data.Industry,
		IsMale:          data.IsMale,
		City:            cityToCityRPC(data.City),
		CityID:          data.CityID,
		PastCompany:     data.PastCompany,
		IsFemale:        data.IsFemale,
		Language:        data.Language,
		// Birthday:        make([]*searchRPC.Date, 0, len(data.Birthday)),
		MinAge:   data.MinAge,
		MaxAge:   data.MaxAge,
		FullName: data.FullName,
	}

	// for _, date := range data.Birthday {
	// 	filter.Birthday = append(filter.GetBirthday(), requestDataToSearchDateRPC(date))
	// }

	return &filter
}

func cityToCityRPC(data []requests.City) []*searchRPC.City {
	c := make([]*searchRPC.City, 0, len(data))

	for _, city := range data {
		c = append(c, &searchRPC.City{
			City:        city.City,
			Country:     city.Country,
			Subdivision: city.Subdivision,
			ID:          city.ID,
		})
	}

	return c
}

// requestsJobSearchToUserSearchRequestRPC ...
func requestsJobSearchToJobSearchRequestRPC(data *requests.JobSearch) *searchRPC.JobSearchResult {
	if data == nil {
		return nil
	}

	request := searchRPC.JobSearchResult{
		Keywords:           data.Keyword,
		Degree:             data.Degree,
		Subindustry:        data.Subindustry,
		CompanyName:        data.CompanyName,
		WithSalary:         data.WithSalary,
		Currency:           data.Currency,
		MinSalary:          data.MinSalary,
		Skill:              data.Skill,
		IsMinSalaryNull:    data.IsMaxSalaryNull,
		First:              data.First,
		After:              data.After,
		Country:            data.Country,
		City:               cityToCityRPC(data.City),
		CityID:             data.CityID,
		CompanySize:        companySizeToSearchRPCSizeEnum(data.CompanySize),
		Period:             data.Period,
		MaxSalary:          data.MaxSalary,
		WithoutCoverLetter: data.WithoutCoverLetter,
		IsMaxSalaryNull:    data.IsMaxSalaryNull,
		JobType:            data.JobType,
		Language:           data.Language,
		Industry:           data.Industry,
		IsFollowing:        data.IsFollowing,
		DatePosted:         requestsDatePostedEnumToSearchRPCDatePostedEnum(data.DatePosted),
		ExperienceLevel:    requestsExperienceEnumToSearchRPCExperienceEnum(data.ExperienceLevel),
	}

	// for _, exp := range data.ExperienceLevel {
	// 	request.ExperienceLevel = append(request.GetExperienceLevel(), requestsExperienceEnumToSearchRPCExperienceEnum(exp))
	// }

	// for _, date := range data.DatePosted {
	// 	request.DatePosted = append(request.GetDatePosted(), requestDataToSearchDateRPC(date))
	// }

	return &request
}

// requestsCompanySearchToCompanySearchRequestRPC ...
func requestsCompanySearchToCompanySearchRequestRPC(data *requests.CompanySearch) *searchRPC.CompanySearchResult {
	if data == nil {
		return nil
	}

	request := searchRPC.CompanySearchResult{
		After:          data.After,
		City:           cityToCityRPC(data.City),
		CityID:         data.CityID,
		Country:        data.Country,
		First:          data.First,
		FounderIDs:     data.FoundersID,
		FounderNames:   data.FoundersName,
		Industry:       data.Industry,
		IsCompany:      data.IsCompany,
		IsJobOffers:    data.IsJobOffers,
		IsOrganization: data.IsOrganization,
		Keywords:       data.Keyword,
		Name:           data.Name,
		SubIndustry:    data.SubIndustry,
		Type:           companyTypeToCompanyRPCTypeEnum(data.Type),
		Size:           companySizeToCompanyRPCSizeEnum(data.Size),
		BusinessHours:  data.BusinessHours,
		// Rating
	}

	// for _, tt := range data.Type {
	// 	request.Type = append(request.GetType(), accountTypeToCompanyTypeRPC(tt))
	// }

	// for _, size := range data.Size {
	// 	request.Size = append(request.GetSize(), accountSizeToCompanySizeRPC(size))
	// }

	return &request
}

// requestsServiceSearchToServiceSearchRequestRPC ...
func requestsServiceSearchToServiceSearchRequestRPC(data *requests.ServiceSearch) *searchRPC.ServiceSearchRequest {
	if data == nil {
		return nil
	}

	request := searchRPC.ServiceSearchRequest{
		Keywords:     data.Keyword,
		City:         data.CityID,
		Country:      data.CountryID,
		Currency:     data.CurrencyPrice,
		FixedPrice:   data.FixedPrice,
		MaxSalary:    data.MaxPrice,
		MinSalary:    data.MinPrice,
		HourFrom:     data.HourFrom,
		HourTo:       data.HourTo,
		IsAlwaysOpen: data.IsAlwaysOpen,
		Skill:        data.Skills,
		DeliveryTime: deliveryTimeToRPC(data.DeliveryTime),
		LocationType: locationTypeToRPC(data.LocationType),
		Price:        priceTypeToRPC(data.Price),
		ServiceOwner: serviceOwenrToRPC(data.ServiceOwner),
		WeekDays:     weekDaysToRPC(data.WeekDays),
	}

	return &request
}

// requestsServiceRequestSearchToServiceRequestSearchRequestRPC ...
func requestsServiceRequestSearchToServiceRequestSearchRequestRPC(data *requests.ServiceRequest) *searchRPC.ServiceRequestSearchRequest {
	if data == nil {
		return nil
	}

	request := searchRPC.ServiceRequestSearchRequest{
		Keywords:     data.Keyword,
		City:         data.CityID,
		Country:      data.CountryID,
		Currency:     data.CurrencyPrice,
		FixedPrice:   data.FixedPrice,
		MaxSalary:    data.MaxPrice,
		MinSalary:    data.MinPrice,
		Skills:       data.Skills,
		Languages:    data.Languages,
		Tools:        data.Tools,
		DeliveryTime: deliveryTimeToRPC(data.DeliveryTime),
		LocationType: locationTypeToRPC(data.LocationType),
		PriceType:    priceTypeToRPC(data.PriceType),
		ServiceOwner: serviceOwenrToRPC(data.ServiceOwner),
		ProjectType:  projectTypesRPCToRPC(data.ProjectType),
	}

	return &request
}

// requestsCandidateSearchToCandidateSearchRequestRPC ...
func requestsCandidateSearchToCandidateSearchRequestRPC(data *requests.CandidateSearch) *searchRPC.CandidateSearchResult {
	if data == nil {
		return nil
	}

	request := searchRPC.CandidateSearchResult{
		IsWillingToWorkRemotly: data.IsWillingToWorkRemotly,
		First:                  data.First,
		Skill:                  data.Skill,
		School:                 data.School,
		MaxSalary:              data.MaxSalary,
		Industry:               data.Industry,
		IsStudent:              data.IsStudent,
		Country:                data.Country,
		ExperienceLevel:        requestsExperienceEnumToSearchRPCExperienceEnum(data.ExperienceLevel),
		FieldOfStudy:           data.FieldOfStudy,
		IsPossibleToRelocate:   data.IsPossibleToRelocate,
		JobType:                data.JobType,
		Period:                 data.Period,
		MinSalary:              data.MinSalary,
		City:                   cityToCityRPC(data.City),
		CityID:                 data.CityID,
		Keywords:               data.Keyword,
		CurrentCompany:         data.CurrentCompany,
		PastCompany:            data.PastCompany,
		Degree:                 data.Degree,
		IsMaxSalaryNull:        data.IsMaxSalaryNull,
		IsWillingToTravel:      data.IsWillingToTravel,
		IsMinSalaryNull:        data.IsMinSalaryNull,
		After:                  data.After,
		SubIndustry:            data.SubIndustry,
		Language:               data.Language,
		Currency:               data.Currency,
	}

	return &request
}

// // requestDataToSearchDateRPC transforms requests date into RPC date
// func requestDataToSearchDateRPC(data *requests.Date) *searchRPC.Date {
// 	if data == nil {
// 		return nil
// 	}

// 	date := searchRPC.Date{
// 		Day:   data.Day,
// 		Month: data.Month,
// 		Year:  data.Year,
// 	}

// 	return &date
// }

// accountTypeToCompanyTypeRPC transformds company type into RPC type
func accountTypeToCompanyTypeRPC(data company.Type) companyRPC.Type {
	companyType := companyRPC.Type_TYPE_UNKNOWN

	switch data {
	case company.TypePartnership:
		companyType = companyRPC.Type_TYPE_PARTNERSHIP
	case company.TypeSelfEmployed:
		companyType = companyRPC.Type_TYPE_SELF_EMPLOYED
	case company.TypePrivatelyHeld:
		companyType = companyRPC.Type_TYPE_PRIVATELY_HELD
	case company.TypePublicCompany:
		companyType = companyRPC.Type_TYPE_PUBLIC_COMPANY
	case company.TypeGovernmentAgency:
		companyType = companyRPC.Type_TYPE_GOVERNMENT_AGENCY
	case company.TypeSoleProprietorship:
		companyType = companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP
	case company.TypeEducationalInstitution:
		companyType = companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION
	}

	return companyType
}

// accountSizeToCompanySizeRPC transforms company size into RPC size
func accountSizeToCompanySizeRPC(data company.Size) companyRPC.Size {
	size := companyRPC.Size_SIZE_UNKNOWN

	switch data {
	case company.SizeSelfEmployed:
		size = companyRPC.Size_SIZE_SELF_EMPLOYED
	case company.SizeFrom1Till10Employees:
		size = companyRPC.Size_SIZE_1_10_EMPLOYEES
	case company.SizeFrom11Till50Employees:
		size = companyRPC.Size_SIZE_11_50_EMPLOYEES
	case company.SizeFrom51Till200Employees:
		size = companyRPC.Size_SIZE_51_200_EMPLOYEES
	case company.SizeFrom201Till500Employees:
		size = companyRPC.Size_SIZE_201_500_EMPLOYEES
	case company.SizeFrom501Till1000Employees:
		size = companyRPC.Size_SIZE_501_1000_EMPLOYEES
	case company.SizeFrom1001Till5000Employees:
		size = companyRPC.Size_SIZE_1001_5000_EMPLOYEES
	case company.SizeFrom5001Till10000Employees:
		size = companyRPC.Size_SIZE_5001_10000_EMPLOYEES
	case company.SizeFrom10001AndMoreEmployees:
		size = companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES
	}

	return size
}

// from RPC
// userSearchRequestRPCToRequestsUserSearch transforms RPC UerSearch into reuqests UserSearch
func userSearchRequestRPCToRequestsUserSearch(data *searchRPC.UserSearchRequest) *requests.UserSearch {
	if data == nil {
		return nil
	}

	req := requests.UserSearch{
		First:           data.GetFirst(),
		ConnectionsOfID: data.GetConnectionsOfID(),
		CurrentCompany:  data.GetCurrentCompany(),
		Nickname:        data.GetNickname(),
		Skill:           data.GetSkill(),
		Interest:        data.GetInterest(),
		MyConnections:   data.GetMyConnections(),
		CountryID:       data.GetCountryID(),
		School:          data.GetSchool(),
		Position:        data.GetPosition(),
		Firstname:       data.GetFirstname(),
		Lastname:        data.GetLastname(),
		After:           data.GetAfter(),
		Keyword:         data.GetKeywords(),
		Degree:          data.GetDegree(),
		FieldOfStudy:    data.GetFieldOfStudy(),
		IsStudent:       data.GetIsStudent(),
		Industry:        data.GetIndustry(),
		IsMale:          data.GetIsMale(),
		CityID:          data.GetCityID(),
		PastCompany:     data.GetPastCompany(),
		IsFemale:        data.GetIsFemale(),
		Language:        data.GetLanguage(),
		MinAge:          data.GetMinAge(),
		MaxAge:          data.GetMaxAge(),
		// Birthday:        make([]*requests.Date, 0, len(data.GetBirthday())),
		FullName: data.GetFullName(),
	}

	// for _, date := range data.GetBirthday() {
	// 	req.Birthday = append(req.Birthday, searchDateRPCtorequestsDate(date))
	// }

	return &req
}

// saveUserSearchFilterRequestRPCToRequestsSaveUserSearchFilter transforms RPC UerSearch into reuqests UserSearch
func saveUserSearchFilterRequestRPCToRequestsSaveUserSearchFilter(data *searchRPC.SaveUserSearchFilterRequest) *requests.UserSearchFilter {
	if data == nil {
		return nil
	}

	req := requests.UserSearchFilter{
		Name: data.GetName(),
	}
	if data.GetUserSearch() != nil {
		userSearch := userSearchRequestRPCToRequestsUserSearch(data.GetUserSearch())
		req.UserSearch = *userSearch
	}
	return &req
}

// saveJobSearchFilterRequestRPCToRequestsSaveJobSearchFilter ...
func saveJobSearchFilterRequestRPCToRequestsSaveJobSearchFilter(data *searchRPC.SaveJobSearchFilterRequest) *requests.JobSearchFilter {
	if data == nil {
		return nil
	}

	req := requests.JobSearchFilter{
		Name: data.GetName(),
	}
	if data.GetJobSearch() != nil {
		jobSearch := jobSearchRPCToRequestsJobSearch(data.GetJobSearch())
		req.JobSearch = *jobSearch
	}

	return &req
}

// SaveCompanySearchFilterRequestRPCToSaveCompanySearchFilter ...
func SaveCompanySearchFilterRequestRPCToSaveCompanySearchFilter(data *searchRPC.SaveCompanySearchFilterRequest) *requests.CompanySearchFilter {
	if data == nil {
		return nil
	}

	req := requests.CompanySearchFilter{
		Name: data.GetName(),
	}
	if data.GetCompanySearch() != nil {
		companySearch := companySearchRequestRPCToRequestsCompanySearch(data.GetCompanySearch())
		req.CompanySearch = *companySearch
	}

	return &req
}

// SaveServiceSearchFilterRequestRPCToSaveServiceSearchFilter ...
func SaveServiceSearchFilterRequestRPCToSaveServiceSearchFilter(data *searchRPC.SaveServiceSearchFilterRequest) *requests.ServiceSearchFilter {
	if data == nil {
		return nil
	}

	req := requests.ServiceSearchFilter{
		Name: data.GetName(),
	}

	if data.GetCompanyID() != "" {
		req.SetCompanyID(data.GetCompanyID())
	}

	if data.GetServiceSearch() != nil {
		serviceSeach := serviceSearchRequestToRequestsServiceSearch(data.GetServiceSearch())
		req.ServiceSearch = *serviceSeach
	}

	return &req
}

// SaveServiceRequestSearchFilterRequestRPCToSaveServiceRequestSearchFilter ...
func SaveServiceRequestSearchFilterRequestRPCToSaveServiceRequestSearchFilter(data *searchRPC.SaveServiceRequestSearchFilterRequest) *requests.ServiceRequestSearchFilter {
	if data == nil {
		return nil
	}

	req := requests.ServiceRequestSearchFilter{
		Name: data.GetName(),
	}

	if data.GetCompanyID() != "" {
		req.SetCompanyID(data.GetCompanyID())
	}

	if data.GetServiceRequestSearch() != nil {
		serviceSeach := serviceRequestRPCToServiceRequest(data.GetServiceRequestSearch())
		req.ServiceRequest = *serviceSeach
	}

	return &req
}

// SaveCandidateSearchFilterRequestRPCToSaveCandidateSearchFilter ...
func SaveCandidateSearchFilterRequestRPCToSaveCandidateSearchFilter(data *searchRPC.SaveCandidateSearchFilterRequest) *requests.CandidateSearchFilter {
	if data == nil {
		return nil
	}

	req := requests.CandidateSearchFilter{
		Name: data.GetName(),
	}
	if data.GetCandidateSearch() != nil {
		candidateSearch := candidateSearchRequestToRequestsCandidateSearch(data.GetCandidateSearch())
		req.CandidateSearch = *candidateSearch
	}
	return &req
}

// GetFiltersByTypeRPCToRequestsFilterType ...
func GetFiltersByTypeRPCToRequestsFilterType(data searchRPC.FilterTypeRequest) requests.FilterType {
	switch data.Type {
	case searchRPC.FilterTypeRequest_CandidateFilterType:
		return requests.TypeCandidateFilterType
	case searchRPC.FilterTypeRequest_CompanyFilterType:
		return requests.TypeCompanyFilterType
	case searchRPC.FilterTypeRequest_UserFilterType:
		return requests.TypeUserFilterType
	case searchRPC.FilterTypeRequest_JobFilterType:
		return requests.TypeJobFilterType
	}

	return ""
}

// func searchDateRPCtorequestsDate(data *searchRPC.Date) *requests.Date {
// 	if data == nil {
// 		return nil
// 	}

// 	date := requests.Date{
// 		Day:   data.GetDay(),
// 		Month: data.GetMonth(),
// 		Year:  data.GetYear(),
// 	}

// 	return &date
// }

// func filterTypeRPCToRequestsFilterType(data *searchRPC.FilterArr) *[]requests.FilterType {
// 	if data == nil {
// 		return nil
// 	}

// 	arr := make([]requests.FilterType, 0, len(data.GetFilters()))

// 	for _, d := range data.GetFilters() {
// 		arr = append(arr, arr(d))
// 	}

// 	return &arr
// }

func companySearchRequestRPCToRequestsCompanySearch(data *searchRPC.CompanySearchRequest) *requests.CompanySearch {
	if data == nil {
		return nil
	}

	req := requests.CompanySearch{
		After:                  data.GetAfter(),
		CityID:                 data.GetCity(),
		Country:                data.GetCountry(),
		First:                  data.GetFirst(),
		FoundersID:             data.GetFounderIDs(),
		FoundersName:           data.GetFounderNames(),
		Industry:               data.GetIndustry(),
		IsCompany:              data.GetIsCompany(),
		IsJobOffers:            data.GetIsJobOffers(),
		IsOrganization:         data.GetIsOrganization(),
		Keyword:                data.GetKeywords(),
		Name:                   data.GetName(),
		SubIndustry:            data.GetSubIndustry(),
		Type:                   companyRPCToCompanyType(data.GetType()),
		Size:                   companyRPCSizeToCompanySize(data.GetSize()),
		BusinessHours:          data.GetBusinessHours(),
		IsCareerCenterOpenened: data.GetIsCareerCenterOpenened(),
		// Rating
	}

	return &req
}

func companyTypeRPCToAccountType(data companyRPC.Type) company.Type {
	companyType := company.TypeUnknown

	switch data {
	case companyRPC.Type_TYPE_PARTNERSHIP:
		companyType = company.TypePartnership
	case companyRPC.Type_TYPE_SELF_EMPLOYED:
		companyType = company.TypeSelfEmployed
	case companyRPC.Type_TYPE_PRIVATELY_HELD:
		companyType = company.TypePrivatelyHeld
	case companyRPC.Type_TYPE_PUBLIC_COMPANY:
		companyType = company.TypePublicCompany
	case companyRPC.Type_TYPE_GOVERNMENT_AGENCY:
		companyType = company.TypeGovernmentAgency
	case companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP:
		companyType = company.TypeSoleProprietorship
	case companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION:
		companyType = company.TypeEducationalInstitution
	}

	return companyType
}

// func sizeRPCToAccountSize(data companyRPC.Size) company.Size {
// 	size := company.SizeUnknown

// 	switch data {
// 	case companyRPC.Size_SIZE_SELF_EMPLOYED:
// 		size = company.SizeSelfEmployed
// 	case companyRPC.Size_SIZE_1_10_EMPLOYEES:
// 		size = company.SizeFrom1Till10Employees
// 	case companyRPC.Size_SIZE_11_50_EMPLOYEES:
// 		size = company.SizeFrom11Till50Employees
// 	case companyRPC.Size_SIZE_51_200_EMPLOYEES:
// 		size = company.SizeFrom51Till200Employees
// 	case companyRPC.Size_SIZE_201_500_EMPLOYEES:
// 		size = company.SizeFrom201Till500Employees
// 	case companyRPC.Size_SIZE_501_1000_EMPLOYEES:
// 		size = company.SizeFrom501Till1000Employees
// 	case companyRPC.Size_SIZE_1001_5000_EMPLOYEES:
// 		size = company.SizeFrom1001Till5000Employees
// 	case companyRPC.Size_SIZE_5001_10000_EMPLOYEES:
// 		size = company.SizeFrom5001Till10000Employees
// 	case companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES:
// 		size = company.SizeFrom10001AndMoreEmployees
// 	}

// 	return size
// }

func jobSearchRPCToRequestsJobSearch(data *searchRPC.JobSearchRequest) *requests.JobSearch {
	if data == nil {
		return nil
	}

	request := requests.JobSearch{
		Keyword:            data.GetKeywords(),
		Degree:             data.GetDegree(),
		Subindustry:        data.GetSubindustry(),
		CompanyName:        data.GetCompanyName(),
		WithSalary:         data.GetWithSalary(),
		ExperienceLevel:    searchRPCExperienceEnumToRequestsExperienceEnum(data.GetExperienceLevel()),
		Currency:           data.GetCurrency(),
		MinSalary:          data.GetMinSalary(),
		Skill:              data.GetSkill(),
		IsMinSalaryNull:    data.GetIsMinSalaryNull(),
		First:              data.GetFirst(),
		After:              data.GetAfter(),
		Country:            data.GetCountry(),
		CityID:             data.GetCity(),
		CompanySize:        searchRPCSizeToCompanySize(data.GetCompanySize()),
		Period:             data.GetPeriod(),
		MaxSalary:          data.GetMaxSalary(),
		WithoutCoverLetter: data.GetWithoutCoverLetter(),
		IsMaxSalaryNull:    data.GetIsMaxSalaryNull(),
		JobType:            data.GetJobType(),
		Language:           data.GetLanguage(),
		Industry:           data.GetIndustry(),
		IsFollowing:        data.GetIsFollowing(),
		DatePosted:         searchRPCDatePostedEnumTorequestsDatePosted(data.GetDatePosted()),
		CompanyIDs:         data.GetCompanyIDs(),
	}

	// for _, date := range data.GetDatePosted() {
	// 	request.DatePosted = append(request.DatePosted, searchDateRPCtorequestsDate(date))
	// }

	return &request
}

func serviceSearchRequestToRequestsServiceSearch(data *searchRPC.ServiceSearchRequest) *requests.ServiceSearch {
	if data == nil {
		return nil
	}

	return &requests.ServiceSearch{
		First:         data.GetFirst(),
		After:         data.GetAfter(),
		CityID:        data.GetCity(),
		CountryID:     data.GetCountry(),
		CurrencyPrice: data.GetCurrency(),
		Keyword:       data.GetKeywords(),
		DeliveryTime:  deliveryTimeRPCToDeliveryTime(data.GetDeliveryTime()),
		LocationType:  locationTypeRPCToLocationType(data.GetLocationType()),
		Price:         priceTypeRPCToPriceType(data.GetPrice()),
		FixedPrice:    data.GetFixedPrice(),
		MinPrice:      data.GetMinSalary(),
		MaxPrice:      data.GetMaxSalary(),
		Skills:        data.GetSkill(),
		IsAlwaysOpen:  data.GetIsAlwaysOpen(),
		HourFrom:      data.GetHourFrom(),
		HourTo:        data.GetHourTo(),
		WeekDays:      weekDaysRPCToWeekDays(data.GetWeekDays()),
		ServiceOwner:  serviceOwenrRPCToServiceOnwer(data.GetServiceOwner()),
	}
}

func serviceRequestRPCToServiceRequest(data *searchRPC.ServiceRequestSearchRequest) *requests.ServiceRequest {
	if data == nil {
		return nil
	}

	return &requests.ServiceRequest{
		First:         data.GetFirst(),
		After:         data.GetAfter(),
		CityID:        data.GetCity(),
		CountryID:     data.GetCountry(),
		CurrencyPrice: data.GetCurrency(),
		Keyword:       data.GetKeywords(),
		ProjectType:   projectTypesRPCToProjectTypes(data.GetProjectType()),
		DeliveryTime:  deliveryTimeRPCToDeliveryTime(data.GetDeliveryTime()),
		LocationType:  locationTypeRPCToLocationType(data.GetLocationType()),
		PriceType:     priceTypeRPCToPriceType(data.GetPriceType()),
		FixedPrice:    data.GetFixedPrice(),
		MinPrice:      data.GetMinSalary(),
		MaxPrice:      data.GetMaxSalary(),
		Skills:        data.GetSkills(),
		Languages:     data.GetLanguages(),
		Tools:         data.GetTools(),
		ServiceOwner:  serviceOwenrRPCToServiceOnwer(data.GetServiceOwner()),
	}
}

func projectTypesRPCToRPC(data []requests.ProjectType) []servicesRPC.RequestProjectTypeEnum {
	if len(data) <= 0 {
		return nil
	}

	projectTypes := make([]servicesRPC.RequestProjectTypeEnum, 0, len(data))

	for _, pt := range data {
		projectTypes = append(projectTypes, projectTypeToRPC(pt))
	}

	return projectTypes
}

func projectTypeToRPC(data requests.ProjectType) servicesRPC.RequestProjectTypeEnum {

	switch data {
	case requests.ProjectTypeOnGoing:
		return servicesRPC.RequestProjectTypeEnum_On_Going_Project
	case requests.ProjectTypeOneTime:
		return servicesRPC.RequestProjectTypeEnum_One_Time_Project

	}
	return servicesRPC.RequestProjectTypeEnum_Not_Sure
}

func projectTypesRPCToProjectTypes(data []servicesRPC.RequestProjectTypeEnum) []requests.ProjectType {
	if len(data) <= 0 {
		return nil
	}

	projectTypes := make([]requests.ProjectType, 0, len(data))

	for _, pt := range data {
		projectTypes = append(projectTypes, projectTypeRPCToProjectType(pt))
	}

	return projectTypes
}

func projectTypeRPCToProjectType(data servicesRPC.RequestProjectTypeEnum) requests.ProjectType {

	switch data {
	case servicesRPC.RequestProjectTypeEnum_On_Going_Project:
		return requests.ProjectTypeOnGoing
	case servicesRPC.RequestProjectTypeEnum_One_Time_Project:
		return requests.ProjectTypeOneTime

	}
	return requests.ProjectTypeNotSure
}

func serviceOwenrRPCToServiceOnwer(data searchRPC.ServiceOwnerEnum) requests.ServiceOwner {

	switch data {
	case searchRPC.ServiceOwnerEnum_Owner_User:
		return requests.ServiceOwnerUser
	case searchRPC.ServiceOwnerEnum_Owner_Company:
		return requests.ServiceOwnerCompany
	}

	return requests.ServiceOwnerAny
}

func serviceOwenrToRPC(data requests.ServiceOwner) searchRPC.ServiceOwnerEnum {

	switch data {
	case requests.ServiceOwnerUser:
		return searchRPC.ServiceOwnerEnum_Owner_User
	case requests.ServiceOwnerCompany:
		return searchRPC.ServiceOwnerEnum_Owner_Company
	}

	return searchRPC.ServiceOwnerEnum_Any_Owner
}

func locationTypeRPCToLocationType(data servicesRPC.LocationEnum) requests.LocationType {
	switch data {
	case servicesRPC.LocationEnum_On_Site_Work:
		return requests.LocationTypeOnSiteWork
	case servicesRPC.LocationEnum_Remote_only:
		return requests.LocationTypeRemoTeOnly
	}

	return requests.LocationTypeAny
}

func locationTypeToRPC(data requests.LocationType) servicesRPC.LocationEnum {
	switch data {
	case requests.LocationTypeOnSiteWork:
		return servicesRPC.LocationEnum_On_Site_Work
	case requests.LocationTypeRemoTeOnly:
		return servicesRPC.LocationEnum_Remote_only
	}

	return servicesRPC.LocationEnum_Location_Any
}

func weekDaysToRPC(data []requests.WeekDay) []servicesRPC.WeekDays {
	if len(data) <= 0 {
		return nil
	}

	weekdays := make([]servicesRPC.WeekDays, 0, len(data))

	for _, w := range data {
		weekdays = append(weekdays, weekDayToRPC(w))
	}

	return weekdays
}

func weekDayToRPC(data requests.WeekDay) servicesRPC.WeekDays {
	switch data {
	case requests.WeekDayMonday:
		return servicesRPC.WeekDays_MONDAY
	case requests.WeekDayTuesday:
		return servicesRPC.WeekDays_TUESDAY
	case requests.WeekDayWednesday:
		return servicesRPC.WeekDays_WEDNESDAY
	case requests.WeekDayThursday:
		return servicesRPC.WeekDays_THURSDAY
	case requests.WeekDayFriday:
		return servicesRPC.WeekDays_FRIDAY
	case requests.WeekDaySaturday:
		return servicesRPC.WeekDays_SATURDAY

	}
	return servicesRPC.WeekDays_SUNDAY
}

func weekDaysRPCToWeekDays(data []servicesRPC.WeekDays) []requests.WeekDay {
	if len(data) <= 0 {
		return nil
	}

	weekdays := make([]requests.WeekDay, 0, len(data))

	for _, w := range data {
		weekdays = append(weekdays, weekDayRPCToWeekDay(w))
	}

	return weekdays
}

func weekDayRPCToWeekDay(data servicesRPC.WeekDays) requests.WeekDay {
	switch data {
	case servicesRPC.WeekDays_MONDAY:
		return requests.WeekDayMonday
	case servicesRPC.WeekDays_TUESDAY:
		return requests.WeekDayTuesday
	case servicesRPC.WeekDays_WEDNESDAY:
		return requests.WeekDayWednesday
	case servicesRPC.WeekDays_THURSDAY:
		return requests.WeekDayThursday
	case servicesRPC.WeekDays_FRIDAY:
		return requests.WeekDayFriday
	case servicesRPC.WeekDays_SATURDAY:
		return requests.WeekDaySaturday

	}
	return requests.WeekDaySunday
}

func priceTypeRPCToPriceType(data servicesRPC.PriceEnum) requests.Price {

	switch data {
	case servicesRPC.PriceEnum_Price_Fixed:
		return requests.PriceFixed
	case servicesRPC.PriceEnum_Price_Hourly:
		return requests.PriceHourly
	case servicesRPC.PriceEnum_Price_Negotiable:
		return requests.PriceNegotiable
	}

	return requests.PriceAny
}

func priceTypeToRPC(data requests.Price) servicesRPC.PriceEnum {

	switch data {
	case requests.PriceFixed:
		return servicesRPC.PriceEnum_Price_Fixed
	case requests.PriceHourly:
		return servicesRPC.PriceEnum_Price_Hourly
	case requests.PriceNegotiable:
		return servicesRPC.PriceEnum_Price_Negotiable
	}

	return servicesRPC.PriceEnum_Price_Any
}

func deliveryTimeRPCToDeliveryTime(data servicesRPC.DeliveryTimeEnum) requests.DeliveryTime {

	switch data {
	case servicesRPC.DeliveryTimeEnum_Custom:
		return requests.DeliveryCustom
	case servicesRPC.DeliveryTimeEnum_Month_And_More:
		return requests.DeliveryMonthAndMore
	case servicesRPC.DeliveryTimeEnum_Up_To_24_Hours:
		return requests.DeliveryUpTo24Hours
	case servicesRPC.DeliveryTimeEnum_Up_To_3_Days:
		return requests.DeliveryUpTo3Days
	case servicesRPC.DeliveryTimeEnum_Up_To_7_Days:
		return requests.DeliveryUpTo7Days
	case servicesRPC.DeliveryTimeEnum_Weeks_1_2:
		return requests.Delivery12Weeks
	case servicesRPC.DeliveryTimeEnum_Weeks_2_4:
		return requests.Delivery2Weeks
	}

	return requests.DeliveryAny

}

func deliveryTimeToRPC(data requests.DeliveryTime) servicesRPC.DeliveryTimeEnum {

	switch data {
	case requests.DeliveryCustom:
		return servicesRPC.DeliveryTimeEnum_Custom
	case requests.DeliveryMonthAndMore:
		return servicesRPC.DeliveryTimeEnum_Month_And_More
	case requests.DeliveryUpTo24Hours:
		return servicesRPC.DeliveryTimeEnum_Up_To_24_Hours
	case requests.DeliveryUpTo3Days:
		return servicesRPC.DeliveryTimeEnum_Up_To_3_Days
	case requests.DeliveryUpTo7Days:
		return servicesRPC.DeliveryTimeEnum_Up_To_7_Days
	case requests.Delivery12Weeks:
		return servicesRPC.DeliveryTimeEnum_Weeks_1_2
	case requests.Delivery2Weeks:
		return servicesRPC.DeliveryTimeEnum_Weeks_2_4
	}

	return servicesRPC.DeliveryTimeEnum_Delivery_Time_Any

}

func candidateSearchRequestToRequestsCandidateSearch(data *searchRPC.CandidateSearchRequest) *requests.CandidateSearch {
	if data == nil {
		return nil
	}

	request := requests.CandidateSearch{
		IsWillingToWorkRemotly: data.GetIsWillingToWorkRemotly(),
		First:                  data.GetFirst(),
		Skill:                  data.GetSkill(),
		School:                 data.GetSchool(),
		MaxSalary:              data.GetMaxSalary(),
		Industry:               data.GetIndustry(),
		IsStudent:              data.GetIsStudent(),
		Country:                data.GetCountry(),
		ExperienceLevel:        searchRPCExperienceEnumToRequestsExperienceEnum(data.GetExperienceLevel()),
		FieldOfStudy:           data.GetFieldOfStudy(),
		IsPossibleToRelocate:   data.GetIsPossibleToRelocate(),
		JobType:                data.GetJobType(),
		Period:                 data.GetPeriod(),
		MinSalary:              data.GetMinSalary(),
		CityID:                 data.GetCity(),
		Keyword:                data.GetKeywords(),
		CurrentCompany:         data.GetCurrentCompany(),
		PastCompany:            data.GetPastCompany(),
		Degree:                 data.GetDegree(),
		IsMaxSalaryNull:        data.GetIsMaxSalaryNull(),
		IsWillingToTravel:      data.GetIsWillingToTravel(),
		IsMinSalaryNull:        data.GetIsMinSalaryNull(),
		After:                  data.GetAfter(),
		SubIndustry:            data.GetSubIndustry(),
		Language:               data.GetLanguage(),
		Currency:               data.GetCurrency(),
	}

	return &request
}

// func SaveUserSearchFilterRPCToRequestsUserSearch(data *searchRPC.UserSearchRequest) string {
// 	if data == nil {
// 		return nil
// 	}

// 	req := requests.UserSearch{
// 		First: data.GetFirst()
// 	},
// }

func searchRPCExperienceEnumToRequestsExperienceEnum(t searchRPC.ExperienceEnum) requests.ExperienceEnum {
	switch t {
	case searchRPC.ExperienceEnum_WithoutExperience:
		return requests.ExperienceEnumWithoutExperience
	case searchRPC.ExperienceEnum_LessThenOneYear:

		return requests.ExperienceEnumLessThenOneYear
	case searchRPC.ExperienceEnum_OneTwoYears:
		return requests.ExperienceEnumOneTwoYears

		// return requests.ExperienceEnumTwoThreeYears

	case searchRPC.ExperienceEnum_TwoThreeYears:
		return requests.ExperienceEnumTwoThreeYears
	case searchRPC.ExperienceEnum_ThreeFiveYears:
		return requests.ExperienceEnumThreeFiveYears
	case searchRPC.ExperienceEnum_FiveSevenyears:
		return requests.ExperienceEnumFiveSevenYears
	case searchRPC.ExperienceEnum_SevenTenYears:
		return requests.ExperienceEnumSevenTenYears
	case searchRPC.ExperienceEnum_TenYearsAndMore:
		return requests.ExperienceEnumTenYearsAndMore
	}

	return requests.ExperienceEnumExpericenUnknown
}

func requestsExperienceEnumToSearchRPCExperienceEnum(t requests.ExperienceEnum) searchRPC.ExperienceEnum {
	switch t {
	case requests.ExperienceEnumWithoutExperience:

		return searchRPC.ExperienceEnum_WithoutExperience
	case requests.ExperienceEnumLessThenOneYear:
		return searchRPC.ExperienceEnum_LessThenOneYear
	case requests.ExperienceEnumOneTwoYears:
		return searchRPC.ExperienceEnum_OneTwoYears

		// return searchRPC.ExperienceEnum_UnknownExperience
	// case requests.ExperienceEnumLessThenOneYear:
	// 	return searchRPC.ExperienceEnum_LessThenOneYear

	case requests.ExperienceEnumTwoThreeYears:
		return searchRPC.ExperienceEnum_TwoThreeYears
	case requests.ExperienceEnumThreeFiveYears:
		return searchRPC.ExperienceEnum_ThreeFiveYears
	case requests.ExperienceEnumFiveSevenYears:
		return searchRPC.ExperienceEnum_FiveSevenyears
	case requests.ExperienceEnumSevenTenYears:
		return searchRPC.ExperienceEnum_SevenTenYears
	case requests.ExperienceEnumTenYearsAndMore:
		return searchRPC.ExperienceEnum_TenYearsAndMore
	}

	return searchRPC.ExperienceEnum_UnknownExperience
}

func requestsDatePostedEnumToSearchRPCDatePostedEnum(t requests.DatePostedEnum) searchRPC.DatePostedEnum {
	switch t {
	case requests.DateEnumPast24Hours:
		return searchRPC.DatePostedEnum_Past24Hours
	case requests.DateEnumPastWeek:
		return searchRPC.DatePostedEnum_PastWeek
	case requests.DateEnumPastMonth:
		return searchRPC.DatePostedEnum_PastMonth
	}

	return searchRPC.DatePostedEnum_Anytime
}

func searchRPCDatePostedEnumTorequestsDatePosted(t searchRPC.DatePostedEnum) requests.DatePostedEnum {
	switch t {
	case searchRPC.DatePostedEnum_Past24Hours:
		return requests.DateEnumPast24Hours
	case searchRPC.DatePostedEnum_PastWeek:
		return requests.DateEnumPastWeek
	case searchRPC.DatePostedEnum_PastMonth:
		return requests.DateEnumPastMonth
	}

	return requests.DateEnumAnytime
}

func companyTypeToCompanyRPCTypeEnum(t company.Type) companyRPC.Type {
	switch t {
	case company.TypeSelfEmployed:
		return companyRPC.Type_TYPE_SELF_EMPLOYED
	case company.TypeEducationalInstitution:
		return companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION
	case company.TypeGovernmentAgency:
		return companyRPC.Type_TYPE_GOVERNMENT_AGENCY
	case company.TypePartnership:
		return companyRPC.Type_TYPE_PARTNERSHIP
	case company.TypePrivatelyHeld:
		return companyRPC.Type_TYPE_PRIVATELY_HELD
	case company.TypePublicCompany:
		return companyRPC.Type_TYPE_PUBLIC_COMPANY
	case company.TypeSoleProprietorship:
		return companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP
	}

	return companyRPC.Type_TYPE_UNKNOWN
}

func companyRPCToCompanyType(t companyRPC.Type) company.Type {
	switch t {
	case companyRPC.Type_TYPE_SELF_EMPLOYED:
		return company.TypeSelfEmployed
	case companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION:
		return company.TypeEducationalInstitution
	case companyRPC.Type_TYPE_GOVERNMENT_AGENCY:
		return company.TypeGovernmentAgency
	case companyRPC.Type_TYPE_PARTNERSHIP:
		return company.TypePartnership
	case companyRPC.Type_TYPE_PRIVATELY_HELD:
		return company.TypePrivatelyHeld
	case companyRPC.Type_TYPE_PUBLIC_COMPANY:
		return company.TypePublicCompany
	case companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP:
		return company.TypeSoleProprietorship
	}

	return company.TypeUnknown
}

func companySizeToSearchRPCSizeEnum(t company.Size) searchRPC.CompanySizeEnum {
	switch t {
	case company.SizeSelfEmployed:
		return searchRPC.CompanySizeEnum_SIZE_SELF_EMPLOYED
	case company.SizeFrom1Till10Employees:
		return searchRPC.CompanySizeEnum_SIZE_1_10_EMPLOYEES
	case company.SizeFrom11Till50Employees:
		return searchRPC.CompanySizeEnum_SIZE_11_50_EMPLOYEES
	case company.SizeFrom51Till200Employees:
		return searchRPC.CompanySizeEnum_SIZE_51_200_EMPLOYEES
	case company.SizeFrom201Till500Employees:
		return searchRPC.CompanySizeEnum_SIZE_201_500_EMPLOYEES
	case company.SizeFrom501Till1000Employees:
		return searchRPC.CompanySizeEnum_SIZE_501_1000_EMPLOYEES
	case company.SizeFrom10001AndMoreEmployees:
		return searchRPC.CompanySizeEnum_SIZE_10001_PLUS_EMPLOYEES
	case company.SizeFrom5001Till10000Employees:
		return searchRPC.CompanySizeEnum_SIZE_5001_10000_EMPLOYEES
	}

	return searchRPC.CompanySizeEnum_SIZE_UNDEFINED
}

func companySizeToCompanyRPCSizeEnum(t company.Size) companyRPC.Size {
	switch t {
	case company.SizeSelfEmployed:
		return companyRPC.Size_SIZE_SELF_EMPLOYED
	case company.SizeFrom1Till10Employees:
		return companyRPC.Size_SIZE_1_10_EMPLOYEES
	case company.SizeFrom11Till50Employees:
		return companyRPC.Size_SIZE_11_50_EMPLOYEES
	case company.SizeFrom51Till200Employees:
		return companyRPC.Size_SIZE_51_200_EMPLOYEES
	case company.SizeFrom201Till500Employees:
		return companyRPC.Size_SIZE_201_500_EMPLOYEES
	case company.SizeFrom501Till1000Employees:
		return companyRPC.Size_SIZE_501_1000_EMPLOYEES
	case company.SizeFrom10001AndMoreEmployees:
		return companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES
	}

	return companyRPC.Size_SIZE_UNKNOWN
}

func companyRPCSizeToCompanySize(t companyRPC.Size) company.Size {
	switch t {
	case companyRPC.Size_SIZE_SELF_EMPLOYED:
		return company.SizeSelfEmployed
	case companyRPC.Size_SIZE_1_10_EMPLOYEES:
		return company.SizeFrom1Till10Employees
	case companyRPC.Size_SIZE_11_50_EMPLOYEES:
		return company.SizeFrom11Till50Employees
	case companyRPC.Size_SIZE_51_200_EMPLOYEES:
		return company.SizeFrom51Till200Employees
	case companyRPC.Size_SIZE_201_500_EMPLOYEES:
		return company.SizeFrom201Till500Employees
	case companyRPC.Size_SIZE_501_1000_EMPLOYEES:
		return company.SizeFrom5001Till10000Employees
	case companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES:
		return company.SizeFrom10001AndMoreEmployees
	}

	return company.SizeUnknown
}

func searchRPCSizeToCompanySize(t searchRPC.CompanySizeEnum) company.Size {
	switch t {
	case searchRPC.CompanySizeEnum_SIZE_SELF_EMPLOYED:
		return company.SizeSelfEmployed
	case searchRPC.CompanySizeEnum_SIZE_1_10_EMPLOYEES:
		return company.SizeFrom1Till10Employees
	case searchRPC.CompanySizeEnum_SIZE_11_50_EMPLOYEES:
		return company.SizeFrom11Till50Employees
	case searchRPC.CompanySizeEnum_SIZE_51_200_EMPLOYEES:
		return company.SizeFrom51Till200Employees
	case searchRPC.CompanySizeEnum_SIZE_201_500_EMPLOYEES:
		return company.SizeFrom201Till500Employees
	case searchRPC.CompanySizeEnum_SIZE_501_1000_EMPLOYEES:
		return company.SizeFrom5001Till10000Employees
	case searchRPC.CompanySizeEnum_SIZE_10001_PLUS_EMPLOYEES:
		return company.SizeFrom10001AndMoreEmployees
	case searchRPC.CompanySizeEnum_SIZE_5001_10000_EMPLOYEES:
		return company.SizeFrom5001Till10000Employees
	}

	return company.SizeUnknown
}
