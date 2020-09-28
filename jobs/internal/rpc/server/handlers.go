package serverRPC

import (
	"context"
	"log"
	"strconv"
	"time"

	careercenter "gitlab.lan/Rightnao-site/microservices/jobs/internal/career-center"
	jobShared "gitlab.lan/Rightnao-site/microservices/jobs/internal/job-functions"

	suitable "gitlab.lan/Rightnao-site/microservices/jobs/internal/suitablefor"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/company"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job"
)

// GetProfile returns candidate profile of user
// TODO: request user profile from userRPC
func (s *Server) GetProfile(ctx context.Context, data *jobsRPC.Empty) (*jobsRPC.CandidateProfile, error) {
	profile, err := s.service.GetCandidateProfile(ctx)
	if err != nil {
		return nil, err
	}
	return candidateProfileToCandidateProfileRPC(profile), nil
}

// SetCareerInterests ...
func (s Server) SetCareerInterests(ctx context.Context, data *jobsRPC.CareerInterests) (*jobsRPC.Empty, error) {
	err := s.service.SetCareerInterests(
		ctx,
		careerInterestRPCToCandidateCareerInterest(data),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// SetOpenFlag ...
func (s Server) SetOpenFlag(ctx context.Context, data *jobsRPC.BoolValue) (*jobsRPC.Empty, error) {
	err := s.service.SetOpenFlag(
		ctx,
		data.GetValue(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// PostJob ...
func (s Server) PostJob(ctx context.Context, data *jobsRPC.PostJobRequest) (*jobsRPC.ID, error) {

	jd := jobsDetailRPCToJoDetails(data.GetDetails())
	meta := jobsMetaRPCToJobMeta(data.GetMetadata())

	post := job.Posting{}

	if jd != nil {
		post.JobDetails = *jd
	}

	if meta != nil {
		post.JobMetadata = *meta
	}

	id, err := s.service.PostJob(
		ctx,
		data.GetCompanyId(),
		&post,
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.ID{
		Id: id,
	}, nil
}

// ChangePost ...
func (s Server) ChangePost(ctx context.Context, data *jobsRPC.PostJobRequest) (*jobsRPC.ID, error) {

	jd := jobsDetailRPCToJoDetails(data.GetDetails())
	meta := jobsMetaRPCToJobMeta(data.GetMetadata())

	post := job.Posting{}

	if jd != nil {
		post.JobDetails = *jd
	}

	if meta != nil {
		post.JobMetadata = *meta
	}

	id, err := s.service.ChangePost(
		ctx,
		data.GetDraftId(),
		data.GetCompanyId(),
		&post,
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.ID{
		Id: id,
	}, nil
}

// SaveDraft ...
func (s Server) SaveDraft(ctx context.Context, data *jobsRPC.PostJobRequest) (*jobsRPC.ID, error) {

	jd := jobsDetailRPCToJoDetails(data.GetDetails())
	meta := jobsMetaRPCToJobMeta(data.GetMetadata())

	post := job.Posting{}

	if jd != nil {
		post.JobDetails = *jd
	}

	if meta != nil {
		post.JobMetadata = *meta
	}

	id, err := s.service.SaveDraft(
		ctx,
		data.GetCompanyId(),
		&post,
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.ID{
		Id: id,
	}, nil
}

// DeleteExpiredPost ...
func (s Server) DeleteExpiredPost(ctx context.Context, data *jobsRPC.DeleteExpiredPostRequest) (*jobsRPC.Empty, error) {
	err := s.service.DeleteExpiredPost(ctx, data.GetPostID(), data.GetId())
	if err != nil {
		return nil, err
	}
	return &jobsRPC.Empty{}, nil
}

// ChangeDraft ...
func (s Server) ChangeDraft(ctx context.Context, data *jobsRPC.PostJobRequest) (*jobsRPC.ID, error) {

	jd := jobsDetailRPCToJoDetails(data.GetDetails())
	meta := jobsMetaRPCToJobMeta(data.GetMetadata())

	post := job.Posting{}

	if jd != nil {
		post.JobDetails = *jd
	}

	if meta != nil {
		post.JobMetadata = *meta
	}

	id, err := s.service.ChangeDraft(
		ctx,
		data.GetDraftId(),
		data.GetCompanyId(),
		&post,
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.ID{
		Id: id,
	}, nil
}

// ActivateJob ...
func (s Server) ActivateJob(ctx context.Context, data *jobsRPC.CompanyIdWithJobId) (*jobsRPC.Empty, error) {
	err := s.service.ActivateJob(ctx, data.GetCompanyId(), data.GetJobId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// PauseJob ...
func (s Server) PauseJob(ctx context.Context, data *jobsRPC.CompanyIdWithJobId) (*jobsRPC.Empty, error) {
	err := s.service.PauseJob(ctx, data.GetCompanyId(), data.GetJobId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// ApplyJob ...
func (s Server) ApplyJob(ctx context.Context, data *jobsRPC.ApplyJobRequest) (*jobsRPC.Empty, error) {
	err := s.service.ApplyJob(
		ctx,
		data.GetJobId(),
		applyJobRequestRPCToJobApplication(data),
		data.GetDocuments(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// IgnoreInvitation ...
func (s Server) IgnoreInvitation(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.Empty, error) {
	err := s.service.IgnoreInvitation(
		ctx,
		data.GetId(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// AddJobView ...
func (s Server) AddJobView(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.Empty, error) {
	err := s.service.AddJobView(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// GetRecommendedJobs ...
func (s Server) GetRecommendedJobs(ctx context.Context, data *jobsRPC.Pagination) (*jobsRPC.JobViewForUserArr, error) {

	recs, err := s.service.GetRecommendedJobs(
		ctx,
		data.GetFirst(),
		data.GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	recommendations := make([]*jobsRPC.JobViewForUser, 0, len(recs))
	for _, r := range recs {
		recommendations = append(recommendations, jobViewForUserViewForUserRPC(r))
	}

	return &jobsRPC.JobViewForUserArr{List: recommendations}, nil
}

// GetJob ...
func (s Server) GetJob(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.JobViewForUser, error) {
	jobView, err := s.service.GetJob(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("TRRRR %+v\n", jobView.JobDetails.AdditionalInfo)

	return jobViewForUserViewForUserRPC(jobView), nil
}

// GetDraft ...
func (s Server) GetDraft(ctx context.Context, data *jobsRPC.PostIDs) (*jobsRPC.PostJobRequest, error) {
	draftView, err := s.service.GetDraft(ctx, data.GetPostID(), data.GetCompanyID())
	if err != nil {
		return nil, err
	}

	return jobViewForCompanyViewForCompanyRPC(draftView), nil
}

// GetPost ...
func (s Server) GetPost(ctx context.Context, data *jobsRPC.PostIDs) (*jobsRPC.PostJobRequest, error) {
	postView, err := s.service.GetPost(ctx, data.GetPostID(), data.GetCompanyID())
	if err != nil {
		return nil, err
	}

	return jobViewForCompanyViewForCompanyRPC(postView), nil
}

// GetJobApplicants ...
func (s Server) GetJobApplicants(ctx context.Context, data *jobsRPC.GetJobApplicantsRequest) (*jobsRPC.JobApplicantArr, error) {
	candidates, err := s.service.GetJobApplicants(
		ctx,
		data.GetCompanyID(),
		data.GetJobID(),
		jobApplicantsSortRPCToCandidateApplicantSort(data.GetSort()),
		data.GetFirst(),
		data.GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	app := make([]*jobsRPC.JobApplicant, 0, len(candidates))
	for _, candidate := range candidates {
		app = append(app, jobApplicantToJobApplicantRPC(candidate))
	}

	return &jobsRPC.JobApplicantArr{
		List: app,
	}, nil
}

// SaveJob ...
func (s Server) SaveJob(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.Empty, error) {
	err := s.service.SaveJob(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// UnsaveJob ...
func (s Server) UnsaveJob(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.Empty, error) {
	err := s.service.UnsaveJob(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// SkipJob ...
func (s Server) SkipJob(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.Empty, error) {
	err := s.service.SkipJob(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// UnskipJob ...
func (s Server) UnskipJob(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.Empty, error) {
	err := s.service.UnskipJob(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// GetSavedJobs ...
func (s Server) GetSavedJobs(ctx context.Context, data *jobsRPC.Pagination) (*jobsRPC.JobViewForUserArr, error) {

	jobs, err := s.service.GetSavedJobs(ctx, data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	jobViews := make([]*jobsRPC.JobViewForUser, 0, len(jobs))
	for _, j := range jobs {
		jobViews = append(jobViews, jobViewForUserViewForUserRPC(j))
	}

	return &jobsRPC.JobViewForUserArr{
		List: jobViews,
	}, nil
}

// GetSkippedJobs ...
func (s Server) GetSkippedJobs(ctx context.Context, data *jobsRPC.Pagination) (*jobsRPC.JobViewForUserArr, error) {
	jobs, err := s.service.GetSkippedJobs(ctx, data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	jobViews := make([]*jobsRPC.JobViewForUser, 0, len(jobs))
	for _, j := range jobs {
		jobViews = append(jobViews, jobViewForUserViewForUserRPC(j))
	}

	return &jobsRPC.JobViewForUserArr{
		List: jobViews,
	}, nil
}

// GetAppliedJobs ...
func (s Server) GetAppliedJobs(ctx context.Context, data *jobsRPC.Pagination) (*jobsRPC.JobViewForUserArr, error) {
	jobs, err := s.service.GetAppliedJobs(ctx, data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	jobViews := make([]*jobsRPC.JobViewForUser, 0, len(jobs))
	for _, j := range jobs {
		jobViews = append(jobViews, jobViewForUserViewForUserRPC(j))
	}

	return &jobsRPC.JobViewForUserArr{
		List: jobViews,
	}, nil
}

// SetJobApplicationSeen ...
func (s Server) SetJobApplicationSeen(ctx context.Context, data *jobsRPC.SetJobApplicationSeenRequest) (*jobsRPC.Empty, error) {
	err := s.service.SetJobApplicationSeen(
		ctx,
		data.GetCompanyId(),
		data.GetJobId(),
		data.GetApplicantId(),
		data.GetSeen(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// SetJobApplicationCategory ...
func (s Server) SetJobApplicationCategory(ctx context.Context, data *jobsRPC.SetJobApplicationCategoryRequest) (*jobsRPC.Empty, error) {
	err := s.service.SetJobApplicationCategory(
		ctx,
		data.GetCompanyId(),
		data.GetJobId(),
		data.GetApplicantId(),
		applicantCategoryEnumRPCToApplicantCategory(data.GetCategory()),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// GetCandidates ...
func (s Server) GetCandidates(ctx context.Context, data *jobsRPC.PaginationWithId) (*jobsRPC.CandidateViewForCompanyArr, error) {
	candidates, err := s.service.GetCandidates(
		ctx,
		data.GetId(),
		data.GetPagination().GetFirst(),
		data.GetPagination().GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	app := make([]*jobsRPC.CandidateViewForCompany, 0, len(candidates))
	for _, candidate := range candidates {
		app = append(app, candidateViewForCompanyToCandidateViewForCompanyRPC(candidate))
	}

	return &jobsRPC.CandidateViewForCompanyArr{
		List: app,
	}, nil
}

// SaveCandidate ...
func (s Server) SaveCandidate(ctx context.Context, data *jobsRPC.CompanyIdWithCandidateId) (*jobsRPC.Empty, error) {
	err := s.service.SaveCandidate(
		ctx,
		data.GetCompanyId(),
		data.GetCandidateId(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// UnsaveCandidate ...
func (s Server) UnsaveCandidate(ctx context.Context, data *jobsRPC.CompanyIdWithCandidateId) (*jobsRPC.Empty, error) {
	err := s.service.UnsaveCandidate(
		ctx,
		data.GetCompanyId(),
		data.GetCandidateId(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// SkipCandidate ....
func (s Server) SkipCandidate(ctx context.Context, data *jobsRPC.CompanyIdWithCandidateId) (*jobsRPC.Empty, error) {
	err := s.service.SkipCandidate(
		ctx,
		data.GetCompanyId(),
		data.GetCandidateId(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// UnskipCandidate ....
func (s Server) UnskipCandidate(ctx context.Context, data *jobsRPC.CompanyIdWithCandidateId) (*jobsRPC.Empty, error) {
	err := s.service.UnskipCandidate(
		ctx,
		data.GetCompanyId(),
		data.GetCandidateId(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// GetListOfJobsWithSeenStat ...
func (s Server) GetListOfJobsWithSeenStat(ctx context.Context, data *jobsRPC.PaginationWithId) (*jobsRPC.ViewJobWithSeenStatArr, error) {
	var first, after int32

	if data.GetPagination() != nil {
		first = data.GetPagination().GetFirst()
		after = data.GetPagination().GetAfter()
	}

	jobs, err := s.service.GetListOfJobsWithSeenStat(
		ctx,
		data.GetId(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	views := make([]*jobsRPC.ViewJobWithSeenStat, 0, len(jobs))
	for _, view := range jobs {
		views = append(views, viewJobWithSeenStatArr(view))
	}

	return &jobsRPC.ViewJobWithSeenStatArr{
		List: views,
	}, nil
}

// GetAmountOfApplicantsPerCategory (ctx context.Context, companyID string) (total, unseen, favorite, inReview, disqualified int32, err error)
func (s Server) GetAmountOfApplicantsPerCategory(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.AmountOfApplicantsPerCategory, error) {
	total, unseen, favorite, inReview, disqualified, err := s.service.GetAmountOfApplicantsPerCategory(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.AmountOfApplicantsPerCategory{
		Disqualified: disqualified,
		Favorite:     favorite,
		Total:        total,
		InReview:     inReview,
		Unseen:       unseen,
	}, nil
}

// GetAmountsOfManageCandidates ...
func (s Server) GetAmountsOfManageCandidates(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.AmountsOfManageCandidates, error) {
	saved, skipped, alerts, err := s.service.GetAmountsOfManageCandidates(
		ctx,
		data.GetId(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.AmountsOfManageCandidates{
		Skipped: skipped,
		Saved:   saved,
		Alerts:  alerts,
	}, nil
}

// GetSavedCandidates ...
func (s Server) GetSavedCandidates(ctx context.Context, data *jobsRPC.PaginationWithId) (*jobsRPC.CandidateViewForCompanyArr, error) {
	var first, after int32

	if data.GetPagination() != nil {
		first = data.GetPagination().GetFirst()
		after = data.GetPagination().GetAfter()
	}

	candidates, err := s.service.GetSavedCandidates(
		ctx,
		data.GetId(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	views := make([]*jobsRPC.CandidateViewForCompany, 0, len(candidates))
	for _, view := range candidates {
		views = append(views, candidateViewForCompanyToCandidateViewForCompanyRPC(view))
	}

	return &jobsRPC.CandidateViewForCompanyArr{
		List: views,
	}, nil
}

// GetSkippedCandidates ...
func (s Server) GetSkippedCandidates(ctx context.Context, data *jobsRPC.PaginationWithId) (*jobsRPC.CandidateViewForCompanyArr, error) {
	var first, after int32

	if data.GetPagination() != nil {
		first = data.GetPagination().GetFirst()
		after = data.GetPagination().GetAfter()
	}

	candidates, err := s.service.GetSkippedCandidates(
		ctx,
		data.GetId(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	views := make([]*jobsRPC.CandidateViewForCompany, 0, len(candidates))
	for _, view := range candidates {
		views = append(views, candidateViewForCompanyToCandidateViewForCompanyRPC(view))
	}

	return &jobsRPC.CandidateViewForCompanyArr{
		List: views,
	}, nil
}

// func

// GetPostedJobs ...
func (s Server) GetPostedJobs(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.JobViewForCompanyArr, error) {
	jobs, err := s.service.GetPostedJobs(
		ctx,
		data.GetId(),
	)
	if err != nil {
		return nil, err
	}

	views := make([]*jobsRPC.JobViewForCompany, 0, len(jobs))
	for _, view := range jobs {
		views = append(views, jobViewForCompanyToJobViewForCompanyRPC(view))
	}

	return &jobsRPC.JobViewForCompanyArr{
		List: views,
	}, nil
}

// GetJobForCompany ...
func (s Server) GetJobForCompany(ctx context.Context, data *jobsRPC.CompanyIdWithJobId) (*jobsRPC.JobViewForCompany, error) {
	jobView, err := s.service.GetJobForCompany(
		ctx,
		data.GetCompanyId(),
		data.GetJobId(),
	)
	if err != nil {
		return nil, err
	}

	return jobViewForCompanyToJobViewForCompanyRPC(jobView), nil
}

// InviteUserToApply ...
func (s Server) InviteUserToApply(ctx context.Context, data *jobsRPC.InviteUserToApplyRequest) (*jobsRPC.Empty, error) {
	err := s.service.InviteUserToApply(
		ctx,
		data.GetCompanyId(),
		data.GetJobId(),
		data.GetUserId(),
		data.GetInvitationText(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// GetInvitedJobs ...
func (s Server) GetInvitedJobs(ctx context.Context, data *jobsRPC.Pagination) (*jobsRPC.JobViewForUserArr, error) {
	jobs, err := s.service.GetInvitedJobs(ctx, data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	jobViews := make([]*jobsRPC.JobViewForUser, 0, len(jobs))
	for _, j := range jobs {
		jobViews = append(jobViews, jobViewForUserViewForUserRPC(j))
	}

	return &jobsRPC.JobViewForUserArr{
		List: jobViews,
	}, nil
}

// ReportJob ...
func (s Server) ReportJob(ctx context.Context, data *jobsRPC.ReportJobRequest) (*jobsRPC.Empty, error) {
	err := s.service.ReportJob(
		ctx,
		data.GetJobId(),
		jobsRPCReportTypeEnumTojobReportType(data.GetType()),
		data.GetText(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// ReportCandidate ...
func (s Server) ReportCandidate(ctx context.Context, data *jobsRPC.ReportCandidateRequest) (*jobsRPC.Empty, error) {
	err := s.service.ReportCandidate(
		ctx,
		data.GetCompanyId(),
		data.GetCandidateId(),
		data.GetText(),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// SaveJobSearchFilter ...
func (s Server) SaveJobSearchFilter(ctx context.Context, data *jobsRPC.NamedJobSearchFilter) (*jobsRPC.Empty, error) {
	err := s.service.SaveJobSearchFilter(
		ctx,
		namedJobSearchFilterRPCToJobSearchFilter(data),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// GetSavedJobSearchFilters ...
func (s Server) GetSavedJobSearchFilters(ctx context.Context, data *jobsRPC.Empty) (*jobsRPC.NamedJobSearchFilterArr, error) {
	filters, err := s.service.GetSavedJobSearchFilters(ctx)
	if err != nil {
		return nil, err
	}

	searchFilters := make([]*jobsRPC.NamedJobSearchFilter, 0, len(filters))
	for _, f := range filters {
		searchFilters = append(searchFilters, jobNamedSearchFilterToNamedJobSearchFilterRPC(f))
	}

	return &jobsRPC.NamedJobSearchFilterArr{
		List: searchFilters,
	}, nil
}

// ...

// SaveCandidateSearchFilter ...
func (s Server) SaveCandidateSearchFilter(ctx context.Context, data *jobsRPC.SaveCandidateSearchFilterRequest) (*jobsRPC.Empty, error) {
	err := s.service.SaveCandidateSearchFilter(
		ctx,
		data.GetCompanyId(),
		namedCandidateSearchFilterRPCTocandidateNamedSearchFilter(data.GetFilter()),
	)
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// GetSavedCandidateSearchFilters ...
func (s Server) GetSavedCandidateSearchFilters(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.NamedCandidateSearchFilterArr, error) {

	filters, err := s.service.GetSavedCandidateSearchFilters(
		ctx,
		data.GetId(),
	)
	if err != nil {
		return nil, err
	}

	candidateFilters := make([]*jobsRPC.NamedCandidateSearchFilter, 0, len(filters))
	for _, f := range filters {
		candidateFilters = append(candidateFilters, candidateNamedSearchFilterToNamedCandidateSearchFilterRPC(f))
	}

	return &jobsRPC.NamedCandidateSearchFilterArr{
		List: candidateFilters,
	}, nil
}

// GetPlanPrices ...
func (s Server) GetPlanPrices(ctx context.Context, data *jobsRPC.GetPlanPricesRequest) (*jobsRPC.GetPlanPricesResult, error) {
	result, err := s.service.GetPlanPrices(ctx, data.GetCompanyID(), data.GetCountries(), data.GetCurrency())
	if err != nil {
		return nil, err
	}

	prices := make([]*jobsRPC.PlanPrice, 0, len(result))

	for _, r := range result {
		renew := make([]float32, 0, len(r.Discounts.Renewal))

		for _, n := range r.Discounts.Renewal {
			renew = append(renew, float32(n))
		}

		prices = append(prices, &jobsRPC.PlanPrice{
			Country:  r.Country,
			Currency: r.Currency,
			Features: &jobsRPC.Features{
				Anonymously: float32(r.Features.Anonymously),
				Language:    float32(r.Features.Language),
				Renewal:     renew,
			},
			PricesPerPlan: &jobsRPC.PricesPerPlan{
				Basic:            float32(r.Prices.Basic),
				Exclusive:        float32(r.Prices.Exclusive),
				Premium:          float32(r.Prices.Premium),
				Professional:     float32(r.Prices.Professional),
				ProfessionalPlus: float32(r.Prices.ProfessionalPlus),
				Standard:         float32(r.Prices.Standard),
				Start:            float32(r.Prices.Start),
			},
		})
	}

	plan := jobsRPC.GetPlanPricesResult{
		Prices: prices,
	}

	return &plan, nil
}

// GetPricingFor ...
func (s Server) GetPricingFor(ctx context.Context, data *jobsRPC.GetPricingRequest) (*jobsRPC.GetPricingResult, error) {

	result, err := s.service.GetPricingFor(
		ctx,
		data.GetCompanyId(),
		jobsMetaRPCToJobMeta(data.GetMeta()),
	)
	if err != nil {
		return nil, err
	}

	return jobPricingResultToGetPricingResultRPC(result), nil
}

// GetAmountOfActiveJobsOfCompany ...
func (s Server) GetAmountOfActiveJobsOfCompany(ctx context.Context, data *jobsRPC.ID) (*jobsRPC.Amount, error) {
	amount, err := s.service.GetAmountOfActiveJobsOfCompany(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Amount{
		Amount: amount,
	}, nil
}

// GetCareerInterestsByIds ...
func (s Server) GetCareerInterestsByIds(ctx context.Context, data *jobsRPC.IDs) (*jobsRPC.CareerInterestsArr, error) {

	careerInterests, err := s.service.GetCareerInterestsByIds(ctx, data.GetID(), data.GetIds(), 999, 0) // TODO: refactor
	if err != nil {
		return nil, err
	}

	careers := make([]*jobsRPC.CareerInterests, 0, len(careerInterests))
	for _, ci := range careerInterests {
		careers = append(careers, candidateCareerInterestToCareerInterestsRPC(ci))
	}

	return &jobsRPC.CareerInterestsArr{
		List: careers,
	}, nil
}

// UploadFileForJob ...
func (s Server) UploadFileForJob(ctx context.Context, data *jobsRPC.File) (*jobsRPC.ID, error) {
	id, err := s.service.UploadFileForJob(ctx, data.GetCompanyID(), data.GetTargetID(), jobsRPCFileToJobFile(data))
	if err != nil {
		return nil, err
	}

	return &jobsRPC.ID{
		Id: id,
	}, nil
}

// UploadFileForApplication ...
func (s Server) UploadFileForApplication(ctx context.Context, data *jobsRPC.File) (*jobsRPC.ID, error) {
	id, err := s.service.UploadFileForApplication(ctx, data.GetUserID(), data.GetTargetID(), jobsRPCFileToJobFile(data))
	if err != nil {
		return nil, err
	}

	return &jobsRPC.ID{
		Id: id,
	}, nil
}

// AddCVInCareerCenter ...
func (s Server) AddCVInCareerCenter(ctx context.Context, data *jobsRPC.AddCVInCareerCenterRequest) (*jobsRPC.Empty, error) {
	err := s.service.AddCVInCareerCenter(ctx, data.GetCompanyID(), careercenter.CVOptions{
		ExpierencedProfessionals: data.GetExpierencedProfessionals(),
		NewJobSeekers:            data.GetNewJobSeekers(),
		YoungProfessionals:       data.GetYoungProfessionals(),
	})
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// GetSavedCVs ...
func (s Server) GetSavedCVs(ctx context.Context, data *jobsRPC.GetSavedCVsRequest) (*jobsRPC.CandidateViewForCompanyArr, error) {
	candidates, err := s.service.GetSavedCVs(
		ctx,
		data.GetId(),
		data.GetPagination().GetFirst(),
		data.GetPagination().GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	app := make([]*jobsRPC.CandidateViewForCompany, 0, len(candidates))
	for _, candidate := range candidates {
		app = append(app, candidateViewForCompanyToCandidateViewForCompanyRPC(candidate))
	}

	return &jobsRPC.CandidateViewForCompanyArr{
		List: app,
	}, nil
}

// RemoveCVs ...
func (s Server) RemoveCVs(ctx context.Context, data *jobsRPC.IDs) (*jobsRPC.Empty, error) {
	err := s.service.RemoveCVs(ctx, data.GetID(), data.GetIds())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// MakeFavoriteCVs ...
func (s Server) MakeFavoriteCVs(ctx context.Context, data *jobsRPC.MakeFavoriteCVsRequest) (*jobsRPC.Empty, error) {
	err := s.service.MakeFavoriteCVs(ctx, data.GetID(), data.GetIDs(), data.GetIsFavourite())
	if err != nil {
		return nil, err
	}

	return &jobsRPC.Empty{}, nil
}

// From RPC

func careerInterestRPCToCandidateCareerInterest(data *jobsRPC.CareerInterests) *candidate.CareerInterests {
	if data == nil {
		return nil
	}

	ci := candidate.CareerInterests{
		Locations:      make([]candidate.Location, 0, len(data.GetLocations())),
		Experience:     eperienceEnumRPCTocandidateExperienceEnum(data.GetExperience()),
		Industry:       data.GetIndustry(),
		Jobs:           data.GetJobs(),
		Relocate:       data.GetRelocate(),
		Remote:         data.GetRemote(),
		SalaryCurrency: data.GetSalaryCurrency(),
		SalaryMax:      data.GetSalaryMax(),
		SalaryMin:      data.GetSalaryMin(),
		Subindustries:  data.GetSubindustry(),
		Travel:         data.GetTravel(),
		SalaryInterval: salaryIntervalRPCToCandidateSalaryInterval(data.GetSalaryInterval()),
		CompanySize:    companySizeRPCToCandidateCompanySize(data.GetCompanySize()),
		JobTypes:       make([]candidate.JobType, 0, len(data.GetJobTypes())),
		Suitable:       suitableForRPCToArray(data.GetSuitableFor()),
		// NormalizedSalary:
	}
	_ = ci.SetUserID(data.GetUserID())

	for _, t := range data.GetJobTypes() {
		ci.JobTypes = append(ci.JobTypes, jobTypeRPCToCandidateJobType(t))
	}

	for _, l := range data.GetLocations() {
		cityID, err := strconv.Atoi(l.GetCityID())
		if err != nil {
			log.Println("cant convert cityID:", err)
		}
		ci.Locations = append(ci.Locations, candidate.Location{
			City:        l.GetCityName(),
			CityID:      int32(cityID),
			Country:     l.GetCountry(),
			Subdivision: l.GetSubdivision(),
		})
	}

	return &ci
}

func eperienceEnumRPCTocandidateExperienceEnum(data jobsRPC.ExperienceEnum) candidate.ExperienceEnum {
	switch data {
	case jobsRPC.ExperienceEnum_WithoutExperience:
		return candidate.ExperienceEnumWithoutExperience
	case jobsRPC.ExperienceEnum_LessThenOneYear:
		return candidate.ExperienceEnumLessThenOneYear
	case jobsRPC.ExperienceEnum_OneTwoYears:
		return candidate.ExperienceEnumOneTwoYears
	case jobsRPC.ExperienceEnum_TwoThreeYears:
		return candidate.ExperienceEnumTwoThreeYears
	case jobsRPC.ExperienceEnum_ThreeFiveYears:
		return candidate.ExperienceEnumThreeFiveYears
	case jobsRPC.ExperienceEnum_FiveSevenyears:
		return candidate.ExperienceEnumFiveSevenYears
	case jobsRPC.ExperienceEnum_SevenTenYears:
		return candidate.ExperienceEnumSevenTenYears
	case jobsRPC.ExperienceEnum_TenYearsAndMore:
		return candidate.ExperienceEnumTenYearsAndMore
	}
	return candidate.ExperienceEnumExpericenUnknown
}

func salaryIntervalRPCToCandidateSalaryInterval(data jobsRPC.SalaryInterval) candidate.SalaryInterval {
	switch data {
	case jobsRPC.SalaryInterval_Hour:
		return candidate.SalaryIntervalHour
	case jobsRPC.SalaryInterval_Month:
		return candidate.SalaryIntervalMonth
	case jobsRPC.SalaryInterval_Year:
		return candidate.SalaryIntervalYear
	}

	return candidate.SalaryIntervalUnknown
}

func jobTypeRPCToCandidateJobType(data jobsRPC.JobType) candidate.JobType {
	switch data {
	case jobsRPC.JobType_Consultancy:
		return candidate.JobTypeConsultancy
	case jobsRPC.JobType_Partner:
		return candidate.JobTypePartner
	case jobsRPC.JobType_FullTime:
		return candidate.JobTypeFullTime
	case jobsRPC.JobType_PartTime:
		return candidate.JobTypePartTime
	case jobsRPC.JobType_Temporary:
		return candidate.JobTypeTemporary
	case jobsRPC.JobType_Volunteer:
		return candidate.JobTypeVolunteer
	case jobsRPC.JobType_Internship:
		return candidate.JobTypeInternship
	case jobsRPC.JobType_Contractual:
		return candidate.JobTypeContractual
	}

	return candidate.JobTypeUnknown
}

func companySizeRPCToCandidateCompanySize(data jobsRPC.CompanySize) company.Size {
	switch data {
	case jobsRPC.CompanySize_SIZE_10001_PLUS_EMPLOYEES:
		return company.Size10001Plus
	case jobsRPC.CompanySize_SIZE_1_10_EMPLOYEES:
		return company.Size1To10
	case jobsRPC.CompanySize_SIZE_11_50_EMPLOYEES:
		return company.Size11To50
	case jobsRPC.CompanySize_SIZE_51_200_EMPLOYEES:
		return company.Size51To200
	case jobsRPC.CompanySize_SIZE_201_500_EMPLOYEES:
		return company.Size201To500
	case jobsRPC.CompanySize_SIZE_501_1000_EMPLOYEES:
		return company.Size501To1000
	case jobsRPC.CompanySize_SIZE_1001_5000_EMPLOYEES:
		return company.Size1001To5000
	case jobsRPC.CompanySize_SIZE_5001_10000_EMPLOYEES:
		return company.Size5001To10000
	}

	return company.SizeUnknown
}

func jobsDetailRPCToJoDetails(data *jobsRPC.JobDetails) *job.Details {
	if data == nil {
		return nil
	}

	benefits := make([]job.Benefit, 0, len(data.GetBenefits()))
	for _, b := range data.GetBenefits() {
		benefits = append(benefits, job.Benefit(b.String()))
	}

	j := job.Details{
		Benefits:              benefits,
		City:                  data.GetCity(),
		Region:                data.GetRegion(),
		Country:               data.GetCountry(),
		LocationType:          jobLocationTypeRPCToEnum(data.GetLocationType()),
		CoverLetter:           data.GetCoverLetter(),
		DeadlineDay:           data.GetDeadlineDay(),
		DeadlineYear:          data.GetDeadlineYear(),
		DeadlineMonth:         data.GetDeadlineMonth(),
		HeaderURL:             data.GetHeaderUrl(),
		HiringMonth:           data.GetHiringMonth(),
		HiringDay:             data.GetHiringDay(),
		HiringYear:            data.GetHiringYear(),
		JobFunctions:          jobFunctionRPCToArray(data.GetJobFunctions()),
		NumberOfPositions:     data.GetNumberOfPositions(),
		PublishDay:            data.GetPublishDay(),
		PublishYear:           data.GetPublishYear(),
		PublishMonth:          data.GetPublishMonth(),
		Required:              applicantQuailificationRPCToApplicantQuailification(data.GetRequired()),
		Preferred:             applicantQuailificationRPCToApplicantQuailification(data.GetPreterred()),
		SalaryCurrency:        data.GetSalaryCurrency(),
		SalaryMax:             data.GetSalaryMax(),
		SalaryMin:             data.GetSalaryMin(),
		AddtionalCompensation: additiionalCompensationRPCToArray(data.GetAdditionalCompensation()),
		AdditionalInfo:        additionalInfoRPCToAdditionalInfo(data.GetAdditionalInfo()),
		Title:                 data.GetTitle(),
		WorkRemotly:           data.GetIsWillingToWorkRemotly(),
		SalaryInterval:        salaryIntervalRPCToCandidateSalaryInterval(data.GetSalaryInterval()),
		EmploymentTypes:       make([]candidate.JobType, 0, len(data.GetEmploymentTypes())),
		Descriptions:          make([]*job.Description, 0, len(data.GetDescriptions())),
	}

	for _, et := range data.GetEmploymentTypes() {
		j.EmploymentTypes = append(j.EmploymentTypes, jobTypeRPCToCandidateJobType(et))
	}

	for _, desc := range data.GetDescriptions() {
		j.Descriptions = append(j.Descriptions, jobDescriptionRPCToJobDescription(desc))
	}

	return &j
}

func jobFunctionRPCToArray(data []jobsRPC.JobFunction) []jobShared.JobFunction {

	jbFns := make([]jobShared.JobFunction, 0, len(data))

	for _, jbFn := range data {
		jbFns = append(jbFns, jobFunctionRPCToEnum(jbFn))
	}

	return jbFns
}

func jobFunctionRPCToEnum(data jobsRPC.JobFunction) jobShared.JobFunction {
	switch data {
	case jobsRPC.JobFunction_Accounting:
		return jobShared.JobFunctionAccounting
	case jobsRPC.JobFunction_Administrative:
		return jobShared.JobFunctionAdministrative
	case jobsRPC.JobFunction_Arts_Design:
		return jobShared.JobFunctionArts_Design
	case jobsRPC.JobFunction_Business_Development:
		return jobShared.JobFunctionBusiness_Development
	case jobsRPC.JobFunction_Community_Social_Services:
		return jobShared.JobFunctionCommunity_Social_Services
	case jobsRPC.JobFunction_Consulting:
		return jobShared.JobFunctionConsulting
	case jobsRPC.JobFunction_Education:
		return jobShared.JobFunctionEducation
	case jobsRPC.JobFunction_Engineering:
		return jobShared.JobFunctionEngineering
	case jobsRPC.JobFunction_Entrepreneurship:
		return jobShared.JobFunctionEntrepreneurship
	case jobsRPC.JobFunction_Finance:
		return jobShared.JobFunctionFinance
	case jobsRPC.JobFunction_Healthcare_Services:
		return jobShared.JobFunctionHealthcare_Services
	case jobsRPC.JobFunction_Human_Resources:
		return jobShared.JobFunctionHuman_Resources
	case jobsRPC.JobFunction_Information_Technology:
		return jobShared.JobFunctionInformation_Technology
	case jobsRPC.JobFunction_Legal:
		return jobShared.JobFunctionLegal
	case jobsRPC.JobFunction_Marketing:
		return jobShared.JobFunctionMarketing
	case jobsRPC.JobFunction_Media_Communications:
		return jobShared.JobFunctionMedia_Communications
	case jobsRPC.JobFunction_Military_Protective_Services:
		return jobShared.JobFunctionMilitary_Protective_Services
	case jobsRPC.JobFunction_Operations:
		return jobShared.JobFunctionOperator
	case jobsRPC.JobFunction_Product_Management:
		return jobShared.JobFunctionProduct_Management
	case jobsRPC.JobFunction_Program_Product_Management:
		return jobShared.JobFunctionProgram_Product_Management
	case jobsRPC.JobFunction_Purchasing:
		return jobShared.JobFunctionPurchasing
	case jobsRPC.JobFunction_Quality_Assurance:
		return jobShared.JobFunctionQuality_Assurance
	case jobsRPC.JobFunction_Real_Estate:
		return jobShared.JobFunctionReal_Estate
	case jobsRPC.JobFunction_Rersearch:
		return jobShared.JobFunctionRersearch
	case jobsRPC.JobFunction_Sales:
		return jobShared.JobFunctionSales
	case jobsRPC.JobFunction_Support:
		return jobShared.JobFunctionSupport

	}

	return jobShared.JobFunctionNone
}

func jobsMetaRPCToJobMeta(data *jobsRPC.JobMeta) *job.Meta {
	if data == nil {
		return nil
	}

	m := job.Meta{
		AdvertisementCountries: data.GetAdvertisementCountries(),
		Anonymous:              data.GetAnonymous(),
		Highlight:              jobRPCHighlightToHighlight(data.GetHighlight()),
		Currency:               data.GetCurrency(),
		NumOfLanguages:         data.GetNumOfLanguages(),
		Renewal:                data.GetRenewal(),
		// JobPlan:                jobPlanRPCToJobPlan(data.GetJobPlan()),
		AmountOfDays: data.GetAmountOfDays(),
	}

	return &m
}

func applicantQuailificationRPCToApplicantQuailification(data *jobsRPC.ApplicationQuailification) job.ApplicantQualification {
	return job.ApplicantQualification{
		Education:      data.GetEducations(),
		Work:           data.GetWork(),
		License:        data.GetLicense(),
		Experience:     jobsRPCJobDetailsExperienceEnumToCandidateExperienceEnum(data.GetExperience()),
		Skills:         data.GetSkills(),
		Language:       languageRPCToArray(data.GetLanguages()),
		ToolTechnology: toolsRPCToArray(data.GetTools()),
	}
}

func toolsRPCToArray(data []*jobsRPC.ApplcantToolsAndTechnology) []job.ToolTechnology {
	if data == nil {
		return nil
	}

	tools := make([]job.ToolTechnology, 0, len(data))

	for _, tool := range data {
		tools = append(tools, toolRPCToTool(tool))
	}

	return tools
}

func toolRPCToTool(data *jobsRPC.ApplcantToolsAndTechnology) job.ToolTechnology {
	if data == nil {
		return job.ToolTechnology{}
	}

	tool := job.ToolTechnology{
		ToolTechnology: data.GetTool(),
		Rank:           rankRPCToRank(data.GetRank()),
	}

	tool.GenerateID()

	return tool
}

func languageRPCToArray(data []*jobsRPC.ApplicantLanguage) []job.Language {
	if data == nil {
		return nil
	}

	languages := make([]job.Language, 0, len(data))

	for _, language := range data {
		languages = append(languages, languageRPCToLanguage(language))
	}

	return languages

}

func languageRPCToLanguage(data *jobsRPC.ApplicantLanguage) job.Language {
	if data == nil {
		return job.Language{}
	}

	language := job.Language{
		Language: data.GetLanguage(),
		Rank:     rankRPCToRank(data.GetRank()),
	}

	language.GenerateID()

	return language

}

func rankRPCToRank(data jobsRPC.ApplicationLevel) *job.Level {
	rank := job.LevelBeginner

	switch data {
	case jobsRPC.ApplicationLevel_Level_Advanced:
		rank = job.LevelAdvanced
	case jobsRPC.ApplicationLevel_Level_Intermediate:
		rank = job.LevelIntermediate
	case jobsRPC.ApplicationLevel_Level_Master:
		rank = job.LevelMaster
	}

	return &rank
}

func jobLocationTypeRPCToEnum(data jobsRPC.LocationType) job.LocationType {

	if data == jobsRPC.LocationType_On_Site {
		return job.LocationTypeOnSite
	}

	return job.LocationTypeRemote
}

func additionalInfoRPCToAdditionalInfo(data *jobsRPC.AdditionalInfo) job.AdditionalInfo {
	if data == nil {
		return job.AdditionalInfo{}
	}

	return job.AdditionalInfo{
		SuitableFor:       suitableForRPCToArray(data.GetSuitableFor()),
		TravelRequirement: travlerRequirementRPCToEnum(data.GetTravelRequirement()),
	}
}

func travlerRequirementRPCToEnum(data jobsRPC.TravelRequirement) job.TravelRequirement {

	switch data {
	case jobsRPC.TravelRequirement_All_time:
		return job.TravelRequirementAll
	case jobsRPC.TravelRequirement_Few_times:
		return job.TravelRequirementFew
	case jobsRPC.TravelRequirement_Once_month:
		return job.TravelRequirementMonth
	case jobsRPC.TravelRequirement_Once_week:
		return job.TravelRequirementWeek
	case jobsRPC.TravelRequirement_Once_year:
		return job.TravelRequirementYear

	}

	return job.TravelRequirementNone
}

func suitableForRPCToArray(data []jobsRPC.SuitableFor) []suitable.SuitableFor {

	suitables := make([]suitable.SuitableFor, 0, len(data))

	for _, suitable := range data {
		suitables = append(suitables, suitableForRPCToSuitableFor(suitable))
	}

	return suitables
}

func suitableForRPCToSuitableFor(data jobsRPC.SuitableFor) suitable.SuitableFor {

	switch data {
	case jobsRPC.SuitableFor_Person_With_Disability:
		return suitable.SuitableForPersonWithDisability
	case jobsRPC.SuitableFor_Student:
		return suitable.SuitableForStudent
	case jobsRPC.SuitableFor_Single_Parent:
		return suitable.SuitableForPersonSingleParent
	case jobsRPC.SuitableFor_Veterans:
		return suitable.SuitableForPersonVeterans

	}
	return suitable.SuitableForNone
}

func additiionalCompensationRPCToArray(data []jobsRPC.AdditionalCompensation) []job.AddtionalCompensation {

	comps := make([]job.AddtionalCompensation, 0, len(data))

	for _, comp := range data {
		comps = append(comps, additiionalCompensationRPCToEnum(comp))
	}

	return comps
}

func additiionalCompensationRPCToEnum(data jobsRPC.AdditionalCompensation) job.AddtionalCompensation {
	compensation := job.AddtionalCompensationBonus

	switch data {
	case jobsRPC.AdditionalCompensation_Profit_Sharing:
		compensation = job.AddtionalCompensationProfit
	case jobsRPC.AdditionalCompensation_Sales_Commission:
		compensation = job.AddtionalCompensationSales
	case jobsRPC.AdditionalCompensation_Tips_Gratuities:
		compensation = job.AddtionalCompensationTips
	}

	return compensation
}

func additiionalCompensationArrayToRPC(data []job.AddtionalCompensation) []jobsRPC.AdditionalCompensation {

	comps := make([]jobsRPC.AdditionalCompensation, 0, len(data))

	for _, comp := range data {
		comps = append(comps, additiionalCompensationToRPC(comp))
	}

	return comps
}

func additiionalCompensationToRPC(data job.AddtionalCompensation) jobsRPC.AdditionalCompensation {

	switch data {
	case job.AddtionalCompensationProfit:
		return jobsRPC.AdditionalCompensation_Profit_Sharing
	case job.AddtionalCompensationSales:
		return jobsRPC.AdditionalCompensation_Sales_Commission
	case job.AddtionalCompensationTips:
		return jobsRPC.AdditionalCompensation_Tips_Gratuities
	}

	return jobsRPC.AdditionalCompensation_Bonus
}

func jobRPCHighlightToHighlight(data jobsRPC.JobHighlight) job.Highlight {
	highlight := job.HighlightNone

	switch data {
	case jobsRPC.JobHighlight_Blue:
		highlight = job.HighlightBlue
	case jobsRPC.JobHighlight_White:
		highlight = job.HighlightWhite
	}
	return highlight
}

func applyJobRequestRPCToJobApplication(data *jobsRPC.ApplyJobRequest) *job.Application {
	if data == nil {
		return nil
	}

	app := job.Application{
		CoverLetter: data.GetCoverLetter(),
		// Documents:   data.GetDocuments(),
		// Documents:   make([]job.File, 0, len(data.GetDocuments())),
		Email: data.GetEmail(),
		Phone: data.GetPhone(),
	}

	// for _, f := range data.GetDocuments() {
	// 	if f != nil {
	// 		app.Documents = append(app.Documents, *jobsRPCFileToJobFile(f))
	// 	}
	// }

	if meta := applicationMetaRPCToJobApplicationMeta(data.GetMetadata()); meta != nil {
		app.Metadata = *meta
	}

	return &app
}

func applicationRPCToJobApplication(data *jobsRPC.Application) *job.Application {
	if data == nil {
		return nil
	}

	app := job.Application{
		CoverLetter: data.GetCoverLetter(),
		// Documents:   data.GetDocuments(),
		Documents: make([]*job.File, 0, len(data.GetDocuments())),
		Email:     data.GetEmail(),
		Phone:     data.GetPhone(),
	}

	for _, f := range data.GetDocuments() {
		if f != nil {
			app.Documents = append(app.Documents, jobsRPCFileToJobFile(f))
		}
	}

	if meta := applicationMetaRPCToJobApplicationMeta(data.GetMetadata()); meta != nil {
		app.Metadata = *meta
	}

	return &app
}

func applicationMetaRPCToJobApplicationMeta(data *jobsRPC.ApplicationMeta) *job.ApplicationMeta {
	if data == nil {
		return nil
	}
	app := job.ApplicationMeta{
		Category: applicantCategoryEnumRPCToApplicantCategory(data.GetCategory()),
		Seen:     data.GetSeen(),
	}

	return &app
}

func jobDescriptionRPCToJobDescription(data *jobsRPC.JobDescription) *job.Description {
	if data == nil {
		return nil
	}

	jd := job.Description{
		Description: data.GetDescription(),
		Language:    data.GetLanguage(),
		WhyUs:       data.GetWhyUs(),
	}

	return &jd
}

func jobPlanRPCToJobPlan(data jobsRPC.JobPlan) job.Plan {
	switch data {
	case jobsRPC.JobPlan_Basic:
		return job.PlanBasic
	case jobsRPC.JobPlan_Start:
		return job.PlanStart
	case jobsRPC.JobPlan_Premium:
		return job.PlanPremium
	case jobsRPC.JobPlan_Exclusive:
		return job.PlanExclusive
	case jobsRPC.JobPlan_Standard:
		return job.PlanStandard
	case jobsRPC.JobPlan_Professional:
		return job.PlanProfessional
	case jobsRPC.JobPlan_ProfessionalPlus:
		return job.PlanProfessionalPlus
	}

	return job.PlanUnknown
}

func namedJobSearchFilterRPCToJobSearchFilter(data *jobsRPC.NamedJobSearchFilter) *job.NamedSearchFilter {
	if data == nil {
		return nil
	}

	filter := job.NamedSearchFilter{
		Name:   data.GetName(),
		Filter: jobSearchFilterRPCToJobSearchFilter(data.GetFilter()),
	}
	filter.SetID(data.GetId())

	return &filter
}

func jobSearchFilterRPCToJobSearchFilter(data *jobsRPC.JobSearchFilter) *job.SearchFilter {
	if data == nil {
		return nil
	}

	filter := job.SearchFilter{
		Keyword: data.Keyword,
		// DatePosted:         dates, // TODO:
		ExperienceLevel:    jobsRPCJobDetailsExperienceEnumToJobExperienceEnum(data.ExperienceLevel),
		Degrees:            data.Degree,
		Countries:          data.Country,
		Cities:             data.City,
		JobTypes:           data.JobType,
		Languages:          data.Language,
		Industries:         data.Industry,
		Subindustries:      data.Subindustry,
		CompanyNames:       data.CompanyName,
		CompanySizes:       companySizeRPCToCandidateCompanySize(data.CompanySize),
		Currency:           data.Currency,
		MinSalary:          data.MinSalary,
		MaxSalary:          data.MaxSalary,
		Period:             data.Period,
		Skills:             data.Skill,
		FollowingCompanies: data.IsFollowing,
		WithoutCoverLetter: data.WithoutCoverLetter,
		WithSalary:         data.WithSalary,
		DatePosted:         jobsRPCDatePostedEnumToJobDatePostedEnum(data.DatePosted),
	}

	return &filter
}

func namedCandidateSearchFilterRPCTocandidateNamedSearchFilter(data *jobsRPC.NamedCandidateSearchFilter) *candidate.NamedSearchFilter {
	if data == nil {
		return nil
	}

	filter := candidate.NamedSearchFilter{
		Name:   data.Name,
		Filter: candidateSearchFilterRPCTocandidateSearchFilter(data.GetFilter()),
	}

	filter.SetID(data.GetId())

	return &filter
}

func candidateSearchFilterRPCTocandidateSearchFilter(data *jobsRPC.CandidateSearchFilter) *candidate.SearchFilter {
	if data == nil {
		return nil
	}

	filter := candidate.SearchFilter{
		Keyword:                data.Keywords,
		Country:                data.Country,
		City:                   data.City,
		CurrentCompany:         data.CurrentCompany,
		PastCompany:            data.PastCompany,
		Industry:               data.Industry,
		SubIndustry:            data.SubIndustry,
		ExperienceLevel:        jobsRPCJobDetailsExperienceEnumToCandidateExperienceEnum(data.ExperienceLevel),
		JobType:                data.JobType,
		Skill:                  data.Skill,
		Language:               data.Language,
		School:                 data.School,
		Degree:                 data.Degree,
		FieldOfStudy:           data.FieldOfStudy,
		IsStudent:              data.IsStudent,
		Currency:               data.Currency,
		Period:                 data.Period,
		MinSalary:              data.MinSalary,
		MaxSalary:              data.MaxSalary,
		IsWillingToTravel:      data.IsWillingToTravel,
		IsWillingToWorkRemotly: data.IsWillingToWorkRemotly,
		IsPossibleToRelocate:   data.IsPossibleToRelocate,
	}

	return &filter
}

func jobApplicantsSortRPCToCandidateApplicantSort(data jobsRPC.GetJobApplicantsRequest_JobApplicantsSort) candidate.ApplicantSort {
	switch data {
	case jobsRPC.GetJobApplicantsRequest_Firstname:
		return candidate.AppicantFirstname
	case jobsRPC.GetJobApplicantsRequest_Lastname:
		return candidate.AppicantLastname
	case jobsRPC.GetJobApplicantsRequest_PostedDate:
		return candidate.AppicantPostedDate
	case jobsRPC.GetJobApplicantsRequest_ExpeirenceYears:
		return candidate.AppicantExperienceYears
	}

	return candidate.AppicantFirstname
}

func applicantCategoryEnumRPCToApplicantCategory(data jobsRPC.ApplicantCategoryEnum) job.ApplicantCategory {
	switch data {
	case jobsRPC.ApplicantCategoryEnum_ApplicantCategoryDisqualified:
		return job.ApplicantCategoryDisqualified
	case jobsRPC.ApplicantCategoryEnum_ApplicantCategoryFavorite:
		return job.ApplicantCategoryFavorite
	case jobsRPC.ApplicantCategoryEnum_ApplicantCategoryInReview:
		return job.ApplicantCategoryInReview
	}

	return job.ApplicantCategoryNone
}

func jobsRPCFileToJobFile(data *jobsRPC.File) *job.File {
	if data == nil {
		return nil
	}

	file := job.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		URL:      data.GetURL(),
	}

	_ = file.SetID(data.GetID())

	return &file
}

func jobsRPCReportTypeEnumTojobReportType(data jobsRPC.ReportJobRequest_ReportJobTypeEnum) job.ReportType {
	switch data {
	case jobsRPC.ReportJobRequest_scam:
		return job.ReportTypeScam
	case jobsRPC.ReportJobRequest_expired:
		return job.ReportTypeExpired
	case jobsRPC.ReportJobRequest_incorrect:
		return job.ReportTypeIncorrect
	case jobsRPC.ReportJobRequest_offensive:
		return job.ReportTypeOffensive
	}
	return job.ReportTypeOther
}

// To RPC

func candidateProfileToCandidateProfileRPC(data *candidate.Profile) *jobsRPC.CandidateProfile {
	if data == nil {
		return nil
	}

	profile := jobsRPC.CandidateProfile{
		IsOpen:          data.IsOpen,
		CareerInterests: candidateCareerInterestToCareerInterestsRPC(data.CareerInterests),
	}

	return &profile
}

func candidateCareerInterestToCareerInterestsRPC(data *candidate.CareerInterests) *jobsRPC.CareerInterests {
	if data == nil {
		return nil
	}

	ci := jobsRPC.CareerInterests{
		Locations:      make([]*jobsRPC.Location, 0, len(data.Locations)),
		Experience:     candidateExperienceEnumToEperienceEnumRPC(data.Experience),
		Industry:       data.Industry,
		Subindustry:    data.Subindustries,
		Jobs:           data.Jobs,
		Relocate:       data.Relocate,
		Remote:         data.Remote,
		SalaryCurrency: data.SalaryCurrency,
		SalaryMax:      data.SalaryMax,
		SalaryMin:      data.SalaryMin,
		Travel:         data.Travel,
		UserID:         data.GetUserID(),
		SalaryInterval: candidateSalaryIntervalToSalaryIntervalRPC(data.SalaryInterval),
		JobTypes:       make([]jobsRPC.JobType, 0, len(data.JobTypes)),
		CompanySize:    candidateCompanySizeToCompanySizeRPC(data.CompanySize),
		IsInvited:      data.IsInvited,
		IsSaved:        data.IsSaved,
		SuitableFor:    suitableForArrayToRPC(data.Suitable),
	}

	for _, t := range data.JobTypes {
		ci.JobTypes = append(ci.JobTypes, candidateJobTypeToJobTypeRPC(t))
	}

	for _, l := range data.Locations {
		ci.Locations = append(ci.Locations, &jobsRPC.Location{
			CityID:      strconv.Itoa(int(l.CityID)),
			CityName:    l.City,
			Country:     l.Country,
			Subdivision: l.Subdivision,
		})
	}

	return &ci
}

func candidateExperienceEnumToEperienceEnumRPC(data candidate.ExperienceEnum) jobsRPC.ExperienceEnum {
	switch data {
	case candidate.ExperienceEnumWithoutExperience:
		return jobsRPC.ExperienceEnum_WithoutExperience
	case candidate.ExperienceEnumLessThenOneYear:
		return jobsRPC.ExperienceEnum_LessThenOneYear
	case candidate.ExperienceEnumOneTwoYears:
		return jobsRPC.ExperienceEnum_OneTwoYears

		// return jobsRPC.ExperienceEnum_UnknownExperience
	// case candidate.ExperienceEnumLessThenOneYear:
	// 	return jobsRPC.ExperienceEnum_LessThenOneYear

	case candidate.ExperienceEnumTwoThreeYears:
		return jobsRPC.ExperienceEnum_TwoThreeYears
	case candidate.ExperienceEnumThreeFiveYears:
		return jobsRPC.ExperienceEnum_ThreeFiveYears
	case candidate.ExperienceEnumFiveSevenYears:
		return jobsRPC.ExperienceEnum_FiveSevenyears
	case candidate.ExperienceEnumSevenTenYears:
		return jobsRPC.ExperienceEnum_SevenTenYears
	case candidate.ExperienceEnumTenYearsAndMore:
		return jobsRPC.ExperienceEnum_TenYearsAndMore
	}
	return jobsRPC.ExperienceEnum_UnknownExperience
}

func candidateSalaryIntervalToSalaryIntervalRPC(data candidate.SalaryInterval) jobsRPC.SalaryInterval {
	switch data {
	case candidate.SalaryIntervalHour:
		return jobsRPC.SalaryInterval_Hour
	case candidate.SalaryIntervalMonth:
		return jobsRPC.SalaryInterval_Month
	case candidate.SalaryIntervalYear:
		return jobsRPC.SalaryInterval_Year
	}

	return jobsRPC.SalaryInterval_Unknown
}

func candidateJobTypeToJobTypeRPC(data candidate.JobType) jobsRPC.JobType {
	switch data {
	case candidate.JobTypeConsultancy:
		return jobsRPC.JobType_Consultancy
	case candidate.JobTypePartner:
		return jobsRPC.JobType_Partner
	case candidate.JobTypeFullTime:
		return jobsRPC.JobType_FullTime
	case candidate.JobTypePartTime:
		return jobsRPC.JobType_PartTime
	case candidate.JobTypeTemporary:
		return jobsRPC.JobType_Temporary
	case candidate.JobTypeVolunteer:
		return jobsRPC.JobType_Volunteer
	case candidate.JobTypeInternship:
		return jobsRPC.JobType_Internship
	case candidate.JobTypeContractual:
		return jobsRPC.JobType_Contractual
	}

	return jobsRPC.JobType_UnknownJobType
}

func candidateCompanySizeToCompanySizeRPC(data company.Size) jobsRPC.CompanySize {
	switch data {
	case company.Size10001Plus:
		return jobsRPC.CompanySize_SIZE_10001_PLUS_EMPLOYEES
	case company.Size1To10:
		return jobsRPC.CompanySize_SIZE_1_10_EMPLOYEES
	case company.Size11To50:
		return jobsRPC.CompanySize_SIZE_11_50_EMPLOYEES
	case company.Size51To200:
		return jobsRPC.CompanySize_SIZE_51_200_EMPLOYEES
	case company.Size201To500:
		return jobsRPC.CompanySize_SIZE_201_500_EMPLOYEES
	case company.Size501To1000:
		return jobsRPC.CompanySize_SIZE_501_1000_EMPLOYEES
	case company.Size1001To5000:
		return jobsRPC.CompanySize_SIZE_1001_5000_EMPLOYEES
	case company.Size5001To10000:
		return jobsRPC.CompanySize_SIZE_5001_10000_EMPLOYEES
	}

	return jobsRPC.CompanySize_SIZE_UNDEFINED
}

func jobViewForUserViewForUserRPC(data *job.ViewForUser) *jobsRPC.JobViewForUser {
	if data == nil {
		return nil
	}

	view := jobsRPC.JobViewForUser{
		Id:             data.GetID(),
		JobDetails:     jobDetailsToJobDetailsRPC(&data.JobDetails),
		CompanyInfo:    companyDetailsToCompanyDetailsRPC(&data.CompanyDetails),
		Application:    jobApplicationToApplicationRPC(&data.Application),
		Metadata:       jobMetaToJobMetaRPC(&data.Metadata),
		IsSaved:        data.IsSaved,
		IsApplied:      data.IsApplied,
		InvitationText: data.InvitationText,
	}

	return &view
}
func jobViewForCompanyViewForCompanyRPC(data *job.ViewForCompany) *jobsRPC.PostJobRequest {
	if data == nil {
		return nil
	}

	view := jobsRPC.PostJobRequest{
		CompanyId: data.GetUserID(),
		Details:   jobDetailsToJobRPCDetails(data.JobDetails),
		Metadata:  jobMetaToJobRPCMeta(data.Metadata),
		DraftId:   data.GetID(), //Question is this correct?

		// Id:                   data.GetID(),
		// JobDetails:           jobDetailsToJobDetailsRPC(&data.JobDetails),
		// Metadata:             jobMetaToJobMetaRPC(&data.Metadata),
		// NumberOfApplications: data.NumberOfApplications,
		// NumberOfViews:        data.NumberOfViews,
		// Status:               data.Status,
		// CreatedAt:            data.CreatedAt.Unix(),
	}

	return &view
}

func jobDetailsToJobDetailsRPC(data *job.Details) *jobsRPC.JobDetails {
	if data == nil {
		return nil
	}

	benefits := make([]jobsRPC.JobDetails_JobBenefit, 0, len(data.Benefits))
	for _, b := range data.Benefits {
		benefits = append(benefits, jobsRPC.JobDetails_JobBenefit(jobsRPC.JobDetails_JobBenefit_value[string(b)]))
	}

	log.Printf("Travle %+v\n", data.AdditionalInfo)

	details := jobsRPC.JobDetails{
		Benefits:               benefits,
		City:                   data.City,
		Country:                data.Country,
		LocationType:           locationTypeToRPC(data.LocationType),
		CoverLetter:            data.CoverLetter,
		DeadlineDay:            data.DeadlineDay,
		DeadlineYear:           data.DeadlineYear,
		DeadlineMonth:          data.DeadlineMonth,
		HeaderUrl:              data.HeaderURL,
		HiringDay:              data.HiringDay,
		HiringYear:             data.HiringYear,
		HiringMonth:            data.HiringMonth,
		JobFunctions:           jobFunctionArrayToRPC(data.JobFunctions),
		NumberOfPositions:      data.NumberOfPositions,
		PublishDay:             data.PublishDay,
		PublishYear:            data.PublishYear,
		PublishMonth:           data.PublishMonth,
		Region:                 data.Region,
		Required:               quailificationToRPC(data.Required),
		Preterred:              quailificationToRPC(data.Preferred),
		AdditionalCompensation: additiionalCompensationArrayToRPC(data.AddtionalCompensation),
		AdditionalInfo:         additionalInfoToRPC(data.AdditionalInfo),
		SalaryCurrency:         data.SalaryCurrency,
		SalaryMax:              data.SalaryMax,
		SalaryMin:              data.SalaryMin,
		Title:                  data.Title,
		Descriptions:           make([]*jobsRPC.JobDescription, 0, len(data.Descriptions)),
		EmploymentTypes:        make([]jobsRPC.JobType, 0, len(data.EmploymentTypes)),
		SalaryInterval:         candidateSalaryIntervalToSalaryIntervalRPC(data.SalaryInterval),
		Files:                  filesToRPC(data.Files),
	}

	for _, desc := range data.Descriptions {
		details.Descriptions = append(details.Descriptions, jobDescriptionToJobDescriptionRPC(desc))
	}

	for _, t := range data.EmploymentTypes {
		details.EmploymentTypes = append(details.EmploymentTypes, candidateJobTypeToJobTypeRPC(t))
	}

	details.Location = &jobsRPC.Location{
		CityID:      data.Location.CityID,
		CityName:    data.Location.CityName,
		Country:     data.Location.Country,
		Subdivision: data.Location.Subdivision,
	}

	return &details
}

func jobFunctionArrayToRPC(data []jobShared.JobFunction) []jobsRPC.JobFunction {

	jbFns := make([]jobsRPC.JobFunction, 0, len(data))

	for _, jbFn := range data {
		jbFns = append(jbFns, jobFunctionEnumToRPC(jbFn))
	}

	return jbFns

}

func jobFunctionEnumToRPC(data jobShared.JobFunction) jobsRPC.JobFunction {

	switch data {
	case "accounting":
		return jobsRPC.JobFunction_Accounting
	case "administrative":
		return jobsRPC.JobFunction_Administrative
	case "arts_design":
		return jobsRPC.JobFunction_Arts_Design
	case "business_development":
		return jobsRPC.JobFunction_Business_Development
	case "community_social_services":
		return jobsRPC.JobFunction_Community_Social_Services
	case "consulting":
		return jobsRPC.JobFunction_Consulting
	case "education":
		return jobsRPC.JobFunction_Education
	case "engineering":
		return jobsRPC.JobFunction_Engineering
	case "entrepreneurship":
		return jobsRPC.JobFunction_Entrepreneurship
	case "finance":
		return jobsRPC.JobFunction_Finance
	case "healthcare_services":
		return jobsRPC.JobFunction_Healthcare_Services
	case "human_resources":
		return jobsRPC.JobFunction_Human_Resources
	case "information_technology":
		return jobsRPC.JobFunction_Information_Technology
	case "legal":
		return jobsRPC.JobFunction_Legal
	case "marketing":
		return jobsRPC.JobFunction_Marketing
	case "media_communications":
		return jobsRPC.JobFunction_Media_Communications
	case "military_protective_services":
		return jobsRPC.JobFunction_Military_Protective_Services
	case "operations":
		return jobsRPC.JobFunction_Operations
	case "product_management":
		return jobsRPC.JobFunction_Product_Management
	case "program_product_management":
		return jobsRPC.JobFunction_Program_Product_Management
	case "purchasing":
		return jobsRPC.JobFunction_Purchasing
	case "quality_assurance":
		return jobsRPC.JobFunction_Quality_Assurance
	case "real_estate":
		return jobsRPC.JobFunction_Real_Estate
	case "rersearch":
		return jobsRPC.JobFunction_Rersearch
	case "sales":
		return jobsRPC.JobFunction_Sales
	case "support":
		return jobsRPC.JobFunction_Support

	}

	return jobsRPC.JobFunction_None_Job_Func
}

func filesToRPC(data []job.File) []*jobsRPC.File {
	files := make([]*jobsRPC.File, 0, len(data))

	for _, file := range data {
		files = append(files, fileToRPC(file))
	}

	return files
}

func fileToRPC(data job.File) *jobsRPC.File {
	return &jobsRPC.File{
		ID:       data.GetID(),
		Name:     data.Name,
		MimeType: data.MimeType,
		URL:      data.URL,
	}
}

func quailificationToRPC(data job.ApplicantQualification) *jobsRPC.ApplicationQuailification {
	return &jobsRPC.ApplicationQuailification{
		Experience: candidateExperienceEnumTojobsRPCJobDetailsExperienceEnum(data.Experience),
		Languages:  languagesToRPC(data.Language),
		Tools:      toolsToRPC(data.ToolTechnology),
		Skills:     data.Skills,
		Educations: data.Education,
		License:    data.License,
		Work:       data.Work,
	}
}

func toolsToRPC(data []job.ToolTechnology) []*jobsRPC.ApplcantToolsAndTechnology {
	tools := make([]*jobsRPC.ApplcantToolsAndTechnology, 0, len(data))

	for _, tool := range data {
		tools = append(tools, toolToRPC(tool))
	}

	return tools
}

func toolToRPC(data job.ToolTechnology) *jobsRPC.ApplcantToolsAndTechnology {
	return &jobsRPC.ApplcantToolsAndTechnology{
		ID:   data.GetID(),
		Tool: data.ToolTechnology,
		Rank: applicantLevelToRPC(data.Rank),
	}
}

func languagesToRPC(data []job.Language) []*jobsRPC.ApplicantLanguage {

	languages := make([]*jobsRPC.ApplicantLanguage, 0, len(data))

	for _, language := range data {
		languages = append(languages, languageToRPC(language))
	}

	return languages
}

func languageToRPC(data job.Language) *jobsRPC.ApplicantLanguage {
	return &jobsRPC.ApplicantLanguage{
		ID:       data.GetID(),
		Language: data.Language,
		Rank:     applicantLevelToRPC(data.Rank),
	}
}

func applicantLevelToRPC(data *job.Level) jobsRPC.ApplicationLevel {

	if data != nil {
		level := *data

		switch level {
		case job.LevelBeginner:
			return jobsRPC.ApplicationLevel_Level_Begginer
		case job.LevelAdvanced:
			return jobsRPC.ApplicationLevel_Level_Advanced
		case job.LevelIntermediate:
			return jobsRPC.ApplicationLevel_Level_Intermediate
		case job.LevelMaster:
			return jobsRPC.ApplicationLevel_Level_Master
		}

		return jobsRPC.ApplicationLevel_Level_Unknown
	}

	return jobsRPC.ApplicationLevel_Level_Unknown
}

func additionalInfoToRPC(data job.AdditionalInfo) *jobsRPC.AdditionalInfo {

	return &jobsRPC.AdditionalInfo{
		SuitableFor:       suitableForArrayToRPC(data.SuitableFor),
		TravelRequirement: travelRequirementToRPC(data.TravelRequirement),
	}
}

func travelRequirementToRPC(data job.TravelRequirement) jobsRPC.TravelRequirement {

	switch data {
	case job.TravelRequirementWeek:
		return jobsRPC.TravelRequirement_Once_week
	case job.TravelRequirementAll:
		return jobsRPC.TravelRequirement_All_time
	case job.TravelRequirementMonth:
		return jobsRPC.TravelRequirement_Once_month
	case job.TravelRequirementYear:
		return jobsRPC.TravelRequirement_Once_year
	case job.TravelRequirementFew:
		return jobsRPC.TravelRequirement_Few_times

	}

	return jobsRPC.TravelRequirement_Travel_req_none
}

func suitableForArrayToRPC(data []suitable.SuitableFor) []jobsRPC.SuitableFor {

	suitables := make([]jobsRPC.SuitableFor, 0, len(data))

	for _, suitable := range data {
		suitables = append(suitables, suitableForToRPC(suitable))
	}

	return suitables
}

func suitableForToRPC(data suitable.SuitableFor) jobsRPC.SuitableFor {

	switch data {
	case suitable.SuitableForStudent:
		return jobsRPC.SuitableFor_Student
	case suitable.SuitableForPersonWithDisability:
		return jobsRPC.SuitableFor_Person_With_Disability
	case suitable.SuitableForPersonSingleParent:
		return jobsRPC.SuitableFor_Single_Parent
	case suitable.SuitableForPersonVeterans:
		return jobsRPC.SuitableFor_Veterans
	}

	return jobsRPC.SuitableFor_None_Suitable
}

func jobDetailsToJobRPCDetails(data job.Details) *jobsRPC.JobDetails {

	benefits := make([]jobsRPC.JobDetails_JobBenefit, 0, len(data.Benefits))
	for _, b := range data.Benefits {
		benefits = append(benefits, jobsRPC.JobDetails_JobBenefit(jobsRPC.JobDetails_JobBenefit_value[string(b)]))
	}

	details := jobsRPC.JobDetails{
		Benefits:          benefits,
		City:              data.City,
		Country:           data.Country,
		CoverLetter:       data.CoverLetter,
		DeadlineDay:       data.DeadlineDay,
		DeadlineYear:      data.DeadlineYear,
		DeadlineMonth:     data.DeadlineMonth,
		HeaderUrl:         data.HeaderURL,
		HiringDay:         data.HiringDay,
		HiringYear:        data.HiringYear,
		HiringMonth:       data.HiringMonth,
		JobFunctions:      jobFunctionArrayToRPC(data.JobFunctions),
		NumberOfPositions: data.NumberOfPositions,
		PublishDay:        data.PublishDay,
		PublishYear:       data.PublishYear,
		PublishMonth:      data.PublishMonth,
		Region:            data.Region,
		SalaryCurrency:    data.SalaryCurrency,
		SalaryMax:         data.SalaryMax,
		SalaryMin:         data.SalaryMin,
		Title:             data.Title,
		Descriptions:      make([]*jobsRPC.JobDescription, 0, len(data.Descriptions)),
		EmploymentTypes:   make([]jobsRPC.JobType, 0, len(data.EmploymentTypes)),
		SalaryInterval:    candidateSalaryIntervalToSalaryIntervalRPC(data.SalaryInterval),
	}

	for _, desc := range data.Descriptions {
		details.Descriptions = append(details.Descriptions, jobDescriptionToJobDescriptionRPC(desc))
	}

	for _, t := range data.EmploymentTypes {
		details.EmploymentTypes = append(details.EmploymentTypes, candidateJobTypeToJobTypeRPC(t))

	}

	return &details
}

func locationTypeToRPC(data job.LocationType) jobsRPC.LocationType {
	if data == job.LocationTypeOnSite {
		return jobsRPC.LocationType_On_Site
	}

	return jobsRPC.LocationType_Remote
}

func jobDescriptionToJobDescriptionRPC(data *job.Description) *jobsRPC.JobDescription {
	if data == nil {
		return nil
	}

	jd := jobsRPC.JobDescription{
		Description: data.Description,
		Language:    data.Language,
		WhyUs:       data.WhyUs,
	}

	return &jd
}

func companyDetailsToCompanyDetailsRPC(data *company.Details) *jobsRPC.CompanyDetails {
	if data == nil {
		return nil
	}

	cd := jobsRPC.CompanyDetails{
		Avatar:      data.CompanyAvatar,
		CompanyId:   data.GetCompanyID(),
		Industry:    data.Industry,
		Subindustry: data.Subindustry,
		URL:         data.CompanyURL,
		// : data.CompanyName
		// ...
	}

	return &cd
}

func jobApplicantToJobApplicantRPC(data *job.Applicant) *jobsRPC.JobApplicant {
	if data == nil {
		return nil
	}
	ja := jobsRPC.JobApplicant{
		UserId:          data.Application.GetUserID(),
		Application:     jobApplicationToApplicationRPC(&data.Application),
		CareerInterests: candidateCareerInterestToCareerInterestsRPC(data.CareerInterests),
	}

	return &ja
}

func jobApplicationToApplicationRPC(data *job.Application) *jobsRPC.Application {
	if data == nil {
		return nil
	}

	app := jobsRPC.Application{
		CoverLetter: data.CoverLetter,
		// Documents:   data.Documents,
		Documents: make([]*jobsRPC.File, 0, len(data.Documents)),
		Email:     data.Email,
		Phone:     data.Phone,
		Metadata:  jobApplicationMetaToApplicationMetaRPC(&data.Metadata),
		CreatedAt: timeToStringDayMonthAndYear(data.CreatedAt),
	}

	for _, f := range data.Documents {
		app.Documents = append(app.Documents, jobFileTojobsRPCFile(f))
	}

	return &app
}

func jobApplicationMetaToApplicationMetaRPC(data *job.ApplicationMeta) *jobsRPC.ApplicationMeta {
	if data == nil {
		return nil
	}

	am := jobsRPC.ApplicationMeta{
		Category: jobApplicantCategoryToApplicantCategoryEnumRPC(data.Category),
		Seen:     data.Seen,
	}

	return &am
}

func candidateViewForCompanyToCandidateViewForCompanyRPC(data *candidate.ViewForCompany) *jobsRPC.CandidateViewForCompany {
	if data == nil {
		return nil
	}

	cv := jobsRPC.CandidateViewForCompany{
		UserId:          data.GetUserID(),
		CareerInterests: candidateCareerInterestToCareerInterestsRPC(&data.CareerInterests),
	}
	if saved := data.IsSaved; saved != nil {
		cv.IsSaved = *saved
	}
	return &cv
}

func jobViewForCompanyToJobViewForCompanyRPC(data *job.ViewForCompany) *jobsRPC.JobViewForCompany {
	if data == nil {
		return nil
	}

	view := jobsRPC.JobViewForCompany{
		CreatedAt:            data.CreatedAt.UnixNano(),
		Id:                   data.GetID(),
		NumberOfApplications: data.NumberOfApplications,
		NumberOfViews:        data.NumberOfViews,
		Status:               data.Status,
		UserId:               data.GetUserID(),
		Metadata:             jobMetaToJobMetaRPC(&data.Metadata),
		JobDetails:           jobDetailsToJobDetailsRPC(&data.JobDetails),
	}

	return &view
}

func jobMetaToJobMetaRPC(data *job.Meta) *jobsRPC.JobMeta {
	if data == nil {
		return nil
	}

	meta := jobsRPC.JobMeta{
		AdvertisementCountries: data.AdvertisementCountries,
		Anonymous:              data.Anonymous,
		Currency:               data.Currency,
		NumOfLanguages:         data.NumOfLanguages,
		Renewal:                data.Renewal,
		Highlight:              highlightToRPC(data.Highlight),
		// JobPlan:                jobPlanToJobPlanRPC(data.JobPlan),
		AmountOfDays: data.AmountOfDays,
	}

	return &meta
}

func jobMetaToJobRPCMeta(data job.Meta) *jobsRPC.JobMeta {

	meta := jobsRPC.JobMeta{
		AdvertisementCountries: data.AdvertisementCountries,
		Anonymous:              data.Anonymous,
		Currency:               data.Currency,
		NumOfLanguages:         data.NumOfLanguages,
		Renewal:                data.Renewal,
		Highlight:              highlightToRPC(data.Highlight),
		// JobPlan:                jobPlanToJobPlanRPC(data.JobPlan),
		AmountOfDays: data.AmountOfDays,
	}

	return &meta
}

func highlightToRPC(data job.Highlight) jobsRPC.JobHighlight {
	highlight := jobsRPC.JobHighlight_None

	switch data {
	case job.HighlightBlue:
		highlight = jobsRPC.JobHighlight_Blue
	case job.HighlightWhite:
		highlight = jobsRPC.JobHighlight_White
	}

	return highlight
}

func jobPlanToJobPlanRPC(data job.Plan) jobsRPC.JobPlan {
	switch data {
	case job.PlanBasic:
		return jobsRPC.JobPlan_Basic
	case job.PlanStart:
		return jobsRPC.JobPlan_Start
	case job.PlanPremium:
		return jobsRPC.JobPlan_Premium
	case job.PlanExclusive:
		return jobsRPC.JobPlan_Exclusive
	case job.PlanStandard:
		return jobsRPC.JobPlan_Standard
	case job.PlanProfessional:
		return jobsRPC.JobPlan_Professional
	case job.PlanProfessionalPlus:
		return jobsRPC.JobPlan_ProfessionalPlus
	}

	return jobsRPC.JobPlan_Start
}

func jobNamedSearchFilterToNamedJobSearchFilterRPC(data *job.NamedSearchFilter) *jobsRPC.NamedJobSearchFilter {
	if data == nil {
		return nil
	}

	filter := jobsRPC.NamedJobSearchFilter{
		Id:     data.GetID(),
		Name:   data.Name,
		Filter: jobSearchFilterToJobSearchFilterRPC(data.Filter),
	}

	return &filter
}

func jobSearchFilterToJobSearchFilterRPC(data *job.SearchFilter) *jobsRPC.JobSearchFilter {
	if data == nil {
		return nil
	}

	filter := jobsRPC.JobSearchFilter{
		Keyword: data.Keyword,
		// DatePosted:         dates, // TODO:
		ExperienceLevel:    jobExperienceEnumTojobsRPCJobDetailsExperienceEnum(data.ExperienceLevel),
		Degree:             data.Degrees,
		Country:            data.Countries,
		City:               data.Cities,
		JobType:            data.JobTypes,
		Language:           data.Languages,
		Industry:           data.Industries,
		Subindustry:        data.Subindustries,
		CompanyName:        data.CompanyNames,
		CompanySize:        candidateCompanySizeToCompanySizeRPC(data.CompanySizes),
		Currency:           data.Currency,
		MinSalary:          data.MinSalary,
		MaxSalary:          data.MaxSalary,
		Period:             data.Period,
		Skill:              data.Skills,
		IsFollowing:        data.FollowingCompanies,
		WithoutCoverLetter: data.WithoutCoverLetter,
		WithSalary:         data.WithSalary,
		DatePosted:         jobDatePostedEnumToSearchRPCDatePostedEnum(data.DatePosted),
	}

	return &filter
}

func candidateNamedSearchFilterToNamedCandidateSearchFilterRPC(data *candidate.NamedSearchFilter) *jobsRPC.NamedCandidateSearchFilter {
	if data == nil {
		return nil
	}

	filter := jobsRPC.NamedCandidateSearchFilter{
		Id:     data.GetID(),
		Name:   data.Name,
		Filter: candidateSearchFilterToCandidateSearchFilter(data.Filter),
	}

	return &filter
}

func candidateSearchFilterToCandidateSearchFilter(data *candidate.SearchFilter) *jobsRPC.CandidateSearchFilter {
	if data == nil {
		return nil
	}

	filter := jobsRPC.CandidateSearchFilter{
		Keywords:               data.Keyword,
		Country:                data.Country,
		City:                   data.City,
		CurrentCompany:         data.CurrentCompany,
		PastCompany:            data.PastCompany,
		Industry:               data.Industry,
		SubIndustry:            data.SubIndustry,
		ExperienceLevel:        candidateExperienceEnumTojobsRPCJobDetailsExperienceEnum(data.ExperienceLevel),
		JobType:                data.JobType,
		Skill:                  data.Skill,
		Language:               data.Language,
		School:                 data.School,
		Degree:                 data.Degree,
		FieldOfStudy:           data.FieldOfStudy,
		IsStudent:              data.IsStudent,
		Currency:               data.Currency,
		Period:                 data.Period,
		MinSalary:              data.MinSalary,
		MaxSalary:              data.MaxSalary,
		IsWillingToTravel:      data.IsWillingToTravel,
		IsWillingToWorkRemotly: data.IsWillingToWorkRemotly,
		IsPossibleToRelocate:   data.IsPossibleToRelocate,
	}

	return &filter
}

func jobPricingResultToGetPricingResultRPC(data *job.PricingResult) *jobsRPC.GetPricingResult {
	if data == nil {
		return nil
	}

	pricing := jobsRPC.GetPricingResult{
		Currency: data.Currency,
		Total:    float32(data.Total),
	}

	countries := make([]*jobsRPC.PricingResultByCountry, 0, len(data.ByCountries))
	for _, country := range data.ByCountries {
		countries = append(countries, jobPricingResultByCountryToPricingResultByCountry(&country))
	}

	pricing.Countries = countries

	return &pricing
}

func jobPricingResultByCountryToPricingResultByCountry(data *job.PricingResultByCountry) *jobsRPC.PricingResultByCountry {
	if data == nil {
		return nil
	}

	pricing := jobsRPC.PricingResultByCountry{
		Country:                 data.Country,
		PlanPrice:               float32(data.PlanPrice),
		RenewalPrice:            float32(data.RenewalPrice),
		LanguagePrice:           float32(data.AdditionalFeatures["language"]),
		PublishAnonymouslyPrice: float32(data.AdditionalFeatures["publish_anonymously"]),
		TotalPrice:              float32(data.TotalPrice),
	}

	return &pricing
}

func viewJobWithSeenStatArr(data *job.ViewJobWithSeenStat) *jobsRPC.ViewJobWithSeenStat {
	if data == nil {
		return nil
	}

	view := jobsRPC.ViewJobWithSeenStat{
		ID:           data.GetID(),
		Title:        data.Title,
		Status:       data.Status,
		TotalAmount:  data.TotalAmount,
		UnseenAmount: data.UnseenAmount,
	}

	return &view
}

func jobApplicantCategoryToApplicantCategoryEnumRPC(data job.ApplicantCategory) jobsRPC.ApplicantCategoryEnum {
	switch data {
	case job.ApplicantCategoryFavorite:
		return jobsRPC.ApplicantCategoryEnum_ApplicantCategoryFavorite
	case job.ApplicantCategoryInReview:
		return jobsRPC.ApplicantCategoryEnum_ApplicantCategoryInReview
	case job.ApplicantCategoryDisqualified:
		return jobsRPC.ApplicantCategoryEnum_ApplicantCategoryDisqualified
	}

	return jobsRPC.ApplicantCategoryEnum_ApplicantCategoryNone
}

func jobFileTojobsRPCFile(data *job.File) *jobsRPC.File {
	if data == nil {
		return nil
	}

	file := jobsRPC.File{
		MimeType: data.MimeType,
		Name:     data.Name,
		URL:      data.URL,
	}

	file.ID = data.GetID()

	return &file
}

// func stringDateToTime(s string) time.Time {
// 	if date, err := time.Parse("2-1-2006", s); err == nil {
// 		return date
// 	}
// 	return time.Time{}
// }
//
// func stringDayMonthAndYearToTime(s string) time.Time {
// 	if date, err := time.Parse("1-2006", s); err == nil {
// 		return date
// 	}
// 	return time.Time{}
// }
//
// func stringYearToDate(s string) time.Time {
// 	if date, err := time.Parse("2006", s); err == nil {
// 		return date
// 	}
// 	return time.Time{}
// }
//
// func timeToStringMonthAndYear(t time.Time) string {
// 	if t == (time.Time{}) {
// 		return ""
// 	}
//
// 	y, m, _ := t.UTC().Date()
// 	return strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
// }
//
func timeToStringDayMonthAndYear(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	y, m, d := t.UTC().Date()
	return strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
}

func jobsRPCJobDetailsExperienceEnumToCandidateExperienceEnum(s jobsRPC.ExperienceEnum) candidate.ExperienceEnum {
	switch s {
	case jobsRPC.ExperienceEnum_WithoutExperience:
		return candidate.ExperienceEnumWithoutExperience
	case jobsRPC.ExperienceEnum_LessThenOneYear:

		return candidate.ExperienceEnumLessThenOneYear
	case jobsRPC.ExperienceEnum_OneTwoYears:
		return candidate.ExperienceEnumOneTwoYears

		// return candidate.ExperienceEnumTwoThreeYears

	case jobsRPC.ExperienceEnum_TwoThreeYears:
		return candidate.ExperienceEnumTwoThreeYears
	case jobsRPC.ExperienceEnum_ThreeFiveYears:
		return candidate.ExperienceEnumThreeFiveYears
	case jobsRPC.ExperienceEnum_FiveSevenyears:
		return candidate.ExperienceEnumFiveSevenYears
	case jobsRPC.ExperienceEnum_SevenTenYears:
		return candidate.ExperienceEnumSevenTenYears
	case jobsRPC.ExperienceEnum_TenYearsAndMore:
		return candidate.ExperienceEnumTenYearsAndMore
	}

	return candidate.ExperienceEnumExpericenUnknown
}

func candidateExperienceEnumTojobsRPCJobDetailsExperienceEnum(s candidate.ExperienceEnum) jobsRPC.ExperienceEnum {
	switch s {
	case candidate.ExperienceEnumWithoutExperience:

		return jobsRPC.ExperienceEnum_WithoutExperience
	case candidate.ExperienceEnumLessThenOneYear:
		return jobsRPC.ExperienceEnum_LessThenOneYear
	case candidate.ExperienceEnumOneTwoYears:
		return jobsRPC.ExperienceEnum_OneTwoYears

		// return jobsRPC.ExperienceEnum_UnknownExperience
	// case candidate.ExperienceEnumLessThenOneYear:
	// 	return jobsRPC.ExperienceEnum_LessThenOneYear

	case candidate.ExperienceEnumTwoThreeYears:
		return jobsRPC.ExperienceEnum_TwoThreeYears
	case candidate.ExperienceEnumThreeFiveYears:
		return jobsRPC.ExperienceEnum_ThreeFiveYears
	case candidate.ExperienceEnumFiveSevenYears:
		return jobsRPC.ExperienceEnum_FiveSevenyears
	case candidate.ExperienceEnumSevenTenYears:
		return jobsRPC.ExperienceEnum_SevenTenYears
	case candidate.ExperienceEnumTenYearsAndMore:
		return jobsRPC.ExperienceEnum_TenYearsAndMore
	}

	return jobsRPC.ExperienceEnum_UnknownExperience
}

func jobExperienceEnumTojobsRPCJobDetailsExperienceEnum(s job.ExperienceEnum) jobsRPC.ExperienceEnum {
	switch s {
	case job.ExperienceEnumWithoutExperience:

		return jobsRPC.ExperienceEnum_WithoutExperience
	case job.ExperienceEnumLessThenOneYear:
		return jobsRPC.ExperienceEnum_LessThenOneYear
	case job.ExperienceEnumOneTwoYears:
		return jobsRPC.ExperienceEnum_OneTwoYears

		// return jobsRPC.ExperienceEnum_UnknownExperience
	// case job.ExperienceEnumLessThenOneYear:
	// 	return jobsRPC.ExperienceEnum_LessThenOneYear

	case job.ExperienceEnumTwoThreeYears:
		return jobsRPC.ExperienceEnum_TwoThreeYears
	case job.ExperienceEnumThreeFiveYears:
		return jobsRPC.ExperienceEnum_ThreeFiveYears
	case job.ExperienceEnumFiveSevenYears:
		return jobsRPC.ExperienceEnum_FiveSevenyears
	case job.ExperienceEnumSevenTenYears:
		return jobsRPC.ExperienceEnum_SevenTenYears
	case job.ExperienceEnumTenYearsAndMore:
		return jobsRPC.ExperienceEnum_TenYearsAndMore
	}

	return jobsRPC.ExperienceEnum_UnknownExperience
}

func jobsRPCJobDetailsExperienceEnumToJobExperienceEnum(s jobsRPC.ExperienceEnum) job.ExperienceEnum {
	switch s {
	case jobsRPC.ExperienceEnum_WithoutExperience:
		return job.ExperienceEnumWithoutExperience
	case jobsRPC.ExperienceEnum_LessThenOneYear:

		return job.ExperienceEnumLessThenOneYear
	case jobsRPC.ExperienceEnum_OneTwoYears:
		return job.ExperienceEnumOneTwoYears

		// return job.ExperienceEnumTwoThreeYears

	case jobsRPC.ExperienceEnum_TwoThreeYears:
		return job.ExperienceEnumTwoThreeYears
	case jobsRPC.ExperienceEnum_ThreeFiveYears:
		return job.ExperienceEnumThreeFiveYears
	case jobsRPC.ExperienceEnum_FiveSevenyears:
		return job.ExperienceEnumFiveSevenYears
	case jobsRPC.ExperienceEnum_SevenTenYears:
		return job.ExperienceEnumSevenTenYears
	case jobsRPC.ExperienceEnum_TenYearsAndMore:
		return job.ExperienceEnumTenYearsAndMore
	}

	return job.ExperienceEnumExpericenUnknown
}

func jobDatePostedEnumToSearchRPCDatePostedEnum(t job.DatePostedEnum) jobsRPC.DatePostedEnum {
	switch t {
	case job.DateEnumPast24Hours:
		return jobsRPC.DatePostedEnum_Past24Hours
	case job.DateEnumPastWeek:
		return jobsRPC.DatePostedEnum_PastWeek
	case job.DateEnumPastMonth:
		return jobsRPC.DatePostedEnum_PastMonth
	}

	return jobsRPC.DatePostedEnum_Anytime
}

func jobsRPCDatePostedEnumToJobDatePostedEnum(t jobsRPC.DatePostedEnum) job.DatePostedEnum {
	switch t {
	case jobsRPC.DatePostedEnum_Past24Hours:
		return job.DateEnumPast24Hours
	case jobsRPC.DatePostedEnum_PastWeek:
		return job.DateEnumPastWeek
	case jobsRPC.DatePostedEnum_PastMonth:
		return job.DateEnumPastMonth
	}

	return job.DateEnumAnytime
}
