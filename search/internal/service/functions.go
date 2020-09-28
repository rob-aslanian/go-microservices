package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/searchRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
)

// UserSearch searches for the user by the input fields
func (s Service) UserSearch(ctx context.Context, data *requests.UserSearch) ([]*userRPC.Profile, int64, error) {
	span := s.tracer.MakeSpan(ctx, "UserSearch")
	defer span.Finish()

	connectionIDs := make([]string, 0)

	err := userSearchValidator(data)
	if err != nil {
		return nil, 0, err // TODO: check
	}

	// get connections ID
	if data.MyConnections || len(data.ConnectionsOfID) > 0 {
		//  get user id
		userID, err := s.authRPC.GetUserID(ctx)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, 0, err
		}

		// get my friends
		if data.MyConnections {
			myFrineds, err := s.networkRPC.GetIDsOfFriends(ctx, userID)
			if err != nil {
				s.tracer.LogError(span, err)
				return nil, 0, err
			}
			connectionIDs = append(connectionIDs, myFrineds...)
		}

		// get connections someone else
		guard := make(chan struct{}, 10)
		wg := sync.WaitGroup{}
		mux := &sync.Mutex{}
		wg.Add(len(data.ConnectionsOfID))

		for i := 0; i < len(data.ConnectionsOfID); i++ {
			guard <- struct{}{}
			go func(n int) {
				defer wg.Done()
				connections, err := s.networkRPC.GetIDsOfFriends(ctx, data.ConnectionsOfID[n])
				if err == nil {
					mux.Lock()
					connectionIDs = append(connectionIDs, connections...)
					mux.Unlock()
				} else {
					s.tracer.LogError(span, err)
				}
				<-guard
			}(i)
		}

		wg.Wait()

		// if there are no connections
		if len(connectionIDs) < 1 {
			return nil, 0, nil
		}
	}

	blockedIDs, err := s.networkRPC.GetBlockedIDs(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	ids, total, err := s.search.UserSearch(ctx, data, connectionIDs, blockedIDs)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	// get users' profiles
	profiles, err := s.userRPC.GetProfilesByID(ctx, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return profiles, total, nil
}

// JobSearch searches for the job by the input fields
func (s Service) JobSearch(ctx context.Context, data *requests.JobSearch) (interface{}, int64, error) {
	span := s.tracer.MakeSpan(ctx, "JobSearch")
	defer span.Finish()

	err := jobSearchValidator(data)
	if err != nil {
		return nil, 0, err // TODO: check
	}

	ids := make([]string, 0)

	// get ids of following companies
	if data.IsFollowing {
		ids, err = s.networkRPC.GetIDsOfFollowingCompanies(ctx)
		if err != nil {
			s.tracer.LogError(span, err)
			// return nil, 0, err
		}
	}

	if len(data.CompanyIDs) > 0 {
		ids = append(ids, data.CompanyIDs...)
	}

	blockedIDs, err := s.networkRPC.GetBlockedIDs(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	ids, total, err := s.search.JobSearch(ctx, data, ids, blockedIDs)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	jobDetails := make([]*jobsRPC.JobViewForUser, len(ids))
	companiesResult := make(map[string]*companyRPC.Profile)

	// get jobs by ids
	guard := make(chan struct{}, 10)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(ids))

	for i := 0; i < len(ids); i++ {
		guard <- struct{}{}
		go func(n int) {
			defer wg.Done()
			job, err := s.jobsRPC.GetJob(ctx, ids[n])
			if err == nil {
				jobDetails[n] = job
				mu.Lock()
				companiesResult[job.GetCompanyInfo().GetCompanyId()] = nil // BUG: concurency
				mu.Unlock()
			} else {
				log.Println("error while getting jobs:", err)
			}
			<-guard
		}(i)
	}

	wg.Wait()

	// get companies by ids
	companyIDs := make([]string, 0, len(ids))

	for key := range companiesResult {
		companyIDs = append(companyIDs, key)
	}

	companies, err := s.companyRPC.GetCompanyProfiles(ctx, companyIDs)
	if err != nil {
		log.Println("error while getting companies:", err)
	}

	for i := range companies {
		companiesResult[companies[i].GetId()] = companies[i]
	}

	// combine all together
	searchResults := make([]*searchRPC.JobResult, 0, len(ids))

	for i := range ids {
		res := searchRPC.JobResult{
			ID: ids[i],
		}

		if len(jobDetails) >= i {
			res.Job = jobDetails[i].GetJobDetails()
			res.Company = companiesResult[jobDetails[i].GetCompanyInfo().GetCompanyId()]
			res.IsApplied = jobDetails[i].GetIsApplied()
			res.IsSaved = jobDetails[i].GetIsSaved()
		} else {
			res.Job = &jobsRPC.JobDetails{}
			log.Println("error: jobs less then ids. Ids:", len(ids), "jobs:", len(jobDetails))
		}

		searchResults = append(searchResults, &res)
	}

	return searchResults, total, nil
}

// CandidateSearch searches for the candidate by the input fields
func (s Service) CandidateSearch(ctx context.Context, companyID string, data *requests.CandidateSearch) (interface{}, int64, error) {
	span := s.tracer.MakeSpan(ctx, "CandidateSearch")
	defer span.Finish()

	err := candidateSearchValidator(data)
	if err != nil {
		return nil, 0, err // TODO: number ...
	}

	blockedIDs, err := s.networkRPC.GetBlockedIDs(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	ids, parentsIDs, total, err := s.search.CandidateSearch(ctx, data, blockedIDs)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	searchResults := make([]*searchRPC.CandidateResult, 0, len(ids))

	careerInterests := make([]*jobsRPC.CareerInterests, 0, len(ids))
	profiles := make([]*userRPC.Profile, 0, len(ids))

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Get jobs
	go func() {
		defer wg.Done()
		career, err := s.jobsRPC.GetCareerInterestsByIDs(ctx, companyID, parentsIDs)
		if err != nil {
			fmt.Println("error while recieving career interest:", err)
		} else {
			careerInterests = append(careerInterests, career...)
		}
	}()

	//get profiles
	go func() {
		defer wg.Done()
		profs, err := s.userRPC.GetProfilesByID(ctx, parentsIDs)
		if err != nil {
			fmt.Println("error while getting profiles:", err)
		} else {
			profiles = append(profiles, profs...)
		}
	}()

	wg.Wait()

	for i := range ids {
		res := searchRPC.CandidateResult{}

		if len(profiles) > i {
			res.Candidates = profiles[i]
		}
		if len(careerInterests) > i {
			for j := range careerInterests {
				if careerInterests[j].GetUserID() == profiles[i].GetID() {
					res.CareerInterests = careerInterests[j]
				}
			}
		}

		searchResults = append(searchResults, &res)
	}

	return searchResults, total, nil
}

// ServiceSearch  for the service by the input fields
func (s Service) ServiceSearch(ctx context.Context, data *requests.ServiceSearch) ([]string, int64, error) {
	span := s.tracer.MakeSpan(ctx, "ServiceSearch")
	defer span.Finish()

	ids, total, err := s.search.ServiceSearch(ctx, data)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return ids, total, nil

}

// ServiceRequestSearch ...
func (s Service) ServiceRequestSearch(ctx context.Context, data *requests.ServiceRequest) ([]string, int64, error) {
	span := s.tracer.MakeSpan(ctx, "ServiceRequestSearch")
	defer span.Finish()

	ids, total, err := s.search.ServiceRequestSearch(ctx, data)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return ids, total, nil
}

// CompanySearch searches for the company by the input fields
func (s Service) CompanySearch(ctx context.Context, data *requests.CompanySearch) (interface{}, int64, error) {
	span := s.tracer.MakeSpan(ctx, "CompanySearch")
	defer span.Finish()

	err := companySearchValidator(data)
	if err != nil {
		return nil, 0, err
	}

	blockedIDs, err := s.networkRPC.GetBlockedIDs(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	ids, total, err := s.search.CompanySearch(ctx, data, blockedIDs)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	// get users' profiles
	profiles, err := s.companyRPC.GetCompanyProfiles(ctx, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return profiles, total, nil
}

// SaveUserSearchFilter saves the filter for user search
func (s Service) SaveUserSearchFilter(ctx context.Context, data *requests.UserSearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveUserSearchFilter")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	err := saveUserSearchFilterValidator(data)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	filter := requests.UserSearchFilter{
		Name:       data.Name,
		UserSearch: data.UserSearch,
		Type:       requests.TypeUserFilterType,
	}

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	id := filter.GenerateID()
	err = filter.SetUserID(userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = s.filterRepository.SaveUserSearchFilter(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// SaveJobSearchFilter saves the filter for job search
func (s Service) SaveJobSearchFilter(ctx context.Context, data *requests.JobSearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveJobSearchFilter")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	err := saveJobSearchFilterValidator(data)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	filter := requests.JobSearchFilter{
		Name:      data.Name,
		JobSearch: data.JobSearch,
		Type:      requests.TypeJobFilterType,
	}

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	id := filter.GenerateID()
	err = filter.SetUserID(userID)
	if err != nil {
		return "", err
	}

	err = s.filterRepository.SaveJobSearchFilter(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// // SaveCandidateSearchFilter ...
// func (s Service) SaveCandidateSearchFilter(ctx context.Context, data *requests.CandidateSearchFilter) (string, error) {
// 	if data == nil {
// 		return "", errors.New("empty search")
// 	}
// 	filter := requests.CandidateSearchFilter{
// 		Name:            data.Name,
// 		CandidateSearch: data.CandidateSearch,
// 		Type:            requests.TypeCandidateFilterType,
// 	}

// 	candidateID, err := s.authRPC.GetUserID(ctx)
// 	if err != nil {
// 		return "", err
// 	}

// 	id := filter.GenerateID()
// 	filter.SetCandidateID(candidateID)

// 	return id, nil
// }

// SaveCompanySearchFilter saves the filter for company search
func (s Service) SaveCompanySearchFilter(ctx context.Context, data *requests.CompanySearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveCompanySearchFilter")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	err := saveCompanySearchFilterValidator(data)
	if err != nil {
		return "", err
	}
	filter := requests.CompanySearchFilter{
		Name:          data.Name,
		CompanySearch: data.CompanySearch,
		Type:          requests.TypeCompanyFilterType,
	}

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	id := filter.GenerateID()
	err = filter.SetUserID(userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = s.filterRepository.SaveCompanySearchFilter(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// SaveServiceSearchFilter saves the filter for company search
func (s Service) SaveServiceSearchFilter(ctx context.Context, data *requests.ServiceSearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveServiceSearchFilter")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	// err := saveCompanySearchFilterValidator(data)
	// if err != nil {
	// 	return "", err
	// }
	filter := requests.ServiceSearchFilter{
		Name:          data.Name,
		ServiceSearch: data.ServiceSearch,
		Type:          requests.TypeServiceFilterType,
	}
	id := filter.GenerateID()

	if data.GetCompanyID() == "" {
		userID, err := s.authRPC.GetUserID(ctx)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}

		err = filter.SetUserID(userID)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}
	} else {
		filter.SetCompanyID(data.GetCompanyID())
	}

	err := s.filterRepository.SaveServiceSearchFilter(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// SaveServiceRequestSearchFilter saves the filter for service request search
func (s Service) SaveServiceRequestSearchFilter(ctx context.Context, data *requests.ServiceRequestSearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveServiceRequestSearchFilter")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	filter := requests.ServiceRequestSearchFilter{
		Name:           data.Name,
		ServiceRequest: data.ServiceRequest,
		Type:           requests.TypeServiceRequestFilterType,
	}
	id := filter.GenerateID()

	if data.GetCompanyID() == "" {
		userID, err := s.authRPC.GetUserID(ctx)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}

		err = filter.SetUserID(userID)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}
	} else {
		filter.SetCompanyID(data.GetCompanyID())
	}

	err := s.filterRepository.SaveServiceRequestSearchFilter(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// GetFiltersByType gets the filter user asks for by type
// func (s Service) GetFiltersByType(ctx context.Context, filterType requests.FilterType) ([]interface{}, error) {
// 	span := s.tracer.MakeSpan(ctx, "GetFiltersByType")
// 	defer span.Finish()

// 	if filterType == "" {
// 		return nil, errors.New("empty filter")
// 	}

// 	userID, err := s.authRPC.GetUserID(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if userID == "" {
// 		return nil, errors.New("invalid id")
// 	}

// 	filter, err := s.filterRepository.GetFiltersByType(ctx, userID, filterType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return filter, nil
// }

// RemoveFilter removes the filter from database which user aks for by the filterID
func (s Service) RemoveFilter(ctx context.Context, filterID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilter")
	defer span.Finish()

	if filterID == "" {
		return errors.New("invalid id")
	}

	err := s.filterRepository.RemoveFilter(ctx, filterID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetAllFilters ...
func (s Service) GetAllFilters(ctx context.Context) ([]interface{}, error) {
	span := s.tracer.MakeSpan(ctx, "GetAllFilters")
	defer span.Finish()

	// lang := "en"

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	if userID == "" {
		s.tracer.LogError(span, err)
		return nil, errors.New("invalid id")
	}

	filter, err := s.filterRepository.GetAllFilters(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return filter, nil
}

// SaveUserSearchFilterForCompany saves the filter for user search by company
func (s Service) SaveUserSearchFilterForCompany(ctx context.Context, companyID string, data *requests.UserSearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveUserSearchFilterForCompany")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	if companyID == "" {
		return "", errors.New("empty companyID")
	}

	err := saveUserSearchFilterValidatorForCompany(data)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	filter := requests.UserSearchFilter{
		Name:       data.Name,
		UserSearch: data.UserSearch,
		Type:       requests.TypeUserFilterType,
	}

	id := filter.GenerateID()
	filter.SetCompanyID(companyID)

	err = s.filterRepository.SaveUserSearchFilterForCompany(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// SaveJobSearchFilterForCompany saves the filter for job search by company
func (s Service) SaveJobSearchFilterForCompany(ctx context.Context, companyID string, data *requests.JobSearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveJobSearchFilterForCompany")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	if companyID == "" {
		return "", errors.New("empty companyID")
	}

	err := saveJobSearchFilterValidatorForCompany(data)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	filter := requests.JobSearchFilter{
		Name:      data.Name,
		JobSearch: data.JobSearch,
		Type:      requests.TypeJobFilterType,
	}

	id := filter.GenerateID()
	filter.SetCompanyID(companyID)

	err = s.filterRepository.SaveJobSearchFilterForCompany(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// SaveCandidateSearchFilterForCompany saves the filter for candidate search by company
func (s Service) SaveCandidateSearchFilterForCompany(ctx context.Context, companyID string, data *requests.CandidateSearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveCandidateSearchFilterForCompany")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	if companyID == "" {
		return "", errors.New("empty companyID")
	}

	err := saveCandidateSearchFilterValidatorForCompany(data)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	filter := requests.CandidateSearchFilter{
		Name:            data.Name,
		CandidateSearch: data.CandidateSearch,
		Type:            requests.TypeCandidateFilterType,
	}

	id := filter.GenerateID()
	filter.SetCompanyID(companyID)

	err = s.filterRepository.SaveCandidateSearchFilterForCompany(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// SaveCompanySearchFilterForCompany saves the filter for company search by company
func (s Service) SaveCompanySearchFilterForCompany(ctx context.Context, companyID string, data *requests.CompanySearchFilter) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveCompanySearchFilterForCompany")
	defer span.Finish()

	if data == nil {
		return "", errors.New("empty search")
	}

	if companyID == "" {
		return "", errors.New("empty companyID")
	}

	err := saveCompanySearchFilterValidatorForCompany(data)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	filter := requests.CompanySearchFilter{
		Name:          data.Name,
		CompanySearch: data.CompanySearch,
		Type:          requests.TypeCompanyFilterType,
	}

	id := filter.GenerateID()
	filter.SetCompanyID(companyID)

	err = s.filterRepository.SaveCompanySearchFilterForCompany(ctx, &filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// GetFiltersByTypeForCompany gets the filter from database by using the companyID to know which company
// is requesting the search and gets filters by their type
// func (s Service) GetFiltersByTypeForCompany(ctx context.Context, companyID string, filterType requests.FilterType) ([]interface{}, error) {
// 	span := s.tracer.MakeSpan(ctx, "GetFiltersByTypeForCompany")
// 	defer span.Finish()

// 	if filterType == "" {
// 		return nil, errors.New("empty filter")
// 	}

// 	if companyID == "" {
// 		return nil, errors.New("invalid id")
// 	}

// 	filter, err := s.filterRepository.GetFiltersByTypeForCompany(ctx, companyID, filterType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return filter, nil
// }

// GetAllFiltersForCompany ...
func (s Service) GetAllFiltersForCompany(ctx context.Context, companyID string) ([]interface{}, error) {
	span := s.tracer.MakeSpan(ctx, "GetFiltersByTypeForCompany")
	defer span.Finish()

	if companyID == "" {
		return nil, errors.New("invalid id")
	}

	filter, err := s.filterRepository.GetAllFiltersForCompany(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range filter {
		switch flt := filter[i].(type) {

		case requests.UserSearchFilter:
			flt.City = make([]requests.City, 0, len(flt.CityID))
			for _, cityID := range flt.CityID {
				n, err := strconv.Atoi(cityID)
				if err != nil {
					return nil, err
				}
				name, sub, country, err := s.infoRPC.GetCityInformationByID(ctx, int32(n), nil)
				if err != nil {
					return nil, err
				}
				flt.City = append(flt.City, requests.City{
					City:        name,
					Country:     country,
					ID:          cityID,
					Subdivision: sub,
				})
			}

		case requests.CandidateSearchFilter:
			flt.City = make([]requests.City, 0, len(flt.CityID))
			for _, cityID := range flt.CityID {
				n, err := strconv.Atoi(cityID)
				if err != nil {
					return nil, err
				}
				name, sub, country, err := s.infoRPC.GetCityInformationByID(ctx, int32(n), nil)
				if err != nil {
					return nil, err
				}
				flt.City = append(flt.City, requests.City{
					City:        name,
					Country:     country,
					ID:          cityID,
					Subdivision: sub,
				})
			}

		case requests.CompanySearchFilter:
			flt.City = make([]requests.City, 0, len(flt.CityID))
			for _, cityID := range flt.CityID {
				n, err := strconv.Atoi(cityID)
				if err != nil {
					return nil, err
				}
				name, sub, country, err := s.infoRPC.GetCityInformationByID(ctx, int32(n), nil)
				if err != nil {
					return nil, err
				}
				flt.City = append(flt.City, requests.City{
					City:        name,
					Country:     country,
					ID:          cityID,
					Subdivision: sub,
				})
			}

		case requests.JobSearchFilter:
			flt.City = make([]requests.City, 0, len(flt.CityID))
			for _, cityID := range flt.CityID {
				n, err := strconv.Atoi(cityID)
				if err != nil {
					return nil, err
				}
				name, sub, country, err := s.infoRPC.GetCityInformationByID(ctx, int32(n), nil)
				if err != nil {
					return nil, err
				}
				flt.City = append(flt.City, requests.City{
					City:        name,
					Country:     country,
					ID:          cityID,
					Subdivision: sub,
				})
			}

		}
	}

	return filter, nil
}

// RemoveFilterForCompany removes the filter from database by using the companyID to know which company
// is requesting the search and removes the filter with filterID
func (s Service) RemoveFilterForCompany(ctx context.Context, filterID, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilterForCompany")
	defer span.Finish()

	if filterID == "" {
		return errors.New("invalid id")
	}

	if companyID == "" {
		return errors.New("invalid id")
	}

	err := s.filterRepository.RemoveFilterForCompany(ctx, filterID, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}
